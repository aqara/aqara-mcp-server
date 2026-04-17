package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Defaults for the stdio lifecycle watchdogs.
//
// These values are deliberately conservative so a well-behaved host (Claude
// Code, OpenClaw, Cursor, ...) that keeps the server around for an interactive
// conversation is not killed prematurely, while a misbehaving or crashed host
// is still reaped within a few minutes instead of leaking forever.
const (
	defaultStdioIdleTimeout    = 5 * time.Minute
	defaultStdioParentPollTick = 2 * time.Second
	minStdioIdleTick           = 1 * time.Second
)

// Environment variables that tune the stdio lifecycle behaviour.
const (
	envStdioIdleTimeout = "AQARA_STDIO_IDLE_TIMEOUT"
	envStdioParentWatch = "AQARA_STDIO_PARENT_WATCH"
)

// stdioLifecycle coordinates the shutdown signals for a stdio MCP server.
//
// An MCP host launches a fresh server process per session and is expected to
// close the child's stdin when the session ends. A well-implemented stdio
// server should therefore exit when any of the following happens:
//
//  1. stdin reaches EOF (the host closed the pipe cleanly),
//  2. SIGINT / SIGTERM is delivered,
//  3. the parent process dies and we are reparented (the host crashed
//     without closing the pipe, which is the observed failure mode that
//     caused leaked aqara-mcp-server processes in Claude Code / OpenClaw),
//  4. no MCP traffic has been observed for the configured idle timeout.
//
// (1) is already handled by mark3labs/mcp-go: its processInputStream loop
// returns nil on io.EOF. (2)–(4) are the job of this type.
type stdioLifecycle struct {
	idleTimeout time.Duration
	parentWatch bool
	parentTick  time.Duration

	lastActivity atomic.Int64 // unix nano
}

func newStdioLifecycle() *stdioLifecycle {
	return &stdioLifecycle{
		idleTimeout: parseDurationEnv(envStdioIdleTimeout, defaultStdioIdleTimeout),
		parentWatch: parseBoolEnv(envStdioParentWatch, true),
		parentTick:  defaultStdioParentPollTick,
	}
}

// touch records that we have just observed activity from the client. Called
// from each MCP hook below.
func (l *stdioLifecycle) touch() {
	l.lastActivity.Store(time.Now().UnixNano())
}

// hooks returns an *server.Hooks wired to reset the idle timer on any RPC
// activity. The hooks fire inside MCPServer.HandleMessage, so they cover both
// requests and notifications regardless of the specific method.
func (l *stdioLifecycle) hooks() *server.Hooks {
	h := &server.Hooks{}
	h.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		l.touch()
	})
	h.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		l.touch()
	})
	h.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		l.touch()
	})
	h.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		l.touch()
	})
	h.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
		l.touch()
		return nil
	})
	return h
}

// serve runs the stdio server with all watchdogs attached. It returns nil on
// any expected shutdown (EOF, signal, parent death, idle timeout) and only
// surfaces unexpected I/O errors.
func (l *stdioLifecycle) serve(mcpServer *server.MCPServer) error {
	// Base context: cancelled on SIGINT / SIGTERM like upstream ServeStdio.
	sigCtx, stopSignal := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stopSignal()

	ctx, cancel := context.WithCancelCause(sigCtx)
	defer cancel(nil)

	l.touch()
	initialPpid := os.Getppid()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		l.watchIdle(ctx, cancel)
	}()
	if l.parentWatch && initialPpid > 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.watchParent(ctx, cancel, initialPpid)
		}()
	}

	stdio := server.NewStdioServer(mcpServer)
	stdio.SetErrorLogger(log.New(os.Stderr, "[stdio] ", log.LstdFlags))

	log.Printf(
		"stdio lifecycle: idle_timeout=%s parent_watch=%v parent_ppid=%d",
		l.idleTimeout, l.parentWatch, initialPpid,
	)

	err := stdio.Listen(ctx, os.Stdin, os.Stdout)

	cancel(nil)
	wg.Wait()

	switch {
	case err == nil:
		log.Printf("stdio lifecycle: shutdown on stdin EOF")
		return nil
	case errors.Is(err, context.Canceled):
		if cause := context.Cause(ctx); cause != nil && !errors.Is(cause, context.Canceled) {
			log.Printf("stdio lifecycle: shutdown cause: %v", cause)
		} else {
			log.Printf("stdio lifecycle: shutdown on signal")
		}
		return nil
	default:
		return err
	}
}

// watchIdle cancels ctx when no MCP activity has been observed for
// idleTimeout. Disabled when idleTimeout <= 0.
func (l *stdioLifecycle) watchIdle(ctx context.Context, cancel context.CancelCauseFunc) {
	if l.idleTimeout <= 0 {
		return
	}

	tick := l.idleTimeout / 4
	if tick < minStdioIdleTick {
		tick = minStdioIdleTick
	}
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			last := time.Unix(0, l.lastActivity.Load())
			if now.Sub(last) >= l.idleTimeout {
				cancel(fmt.Errorf("idle timeout exceeded (%s)", l.idleTimeout))
				return
			}
		}
	}
}

// watchParent cancels ctx when the original parent process disappears. This
// catches the leak scenario where an MCP host (e.g. OpenClaw gateway) dies
// without closing our stdin pipe: getppid() then flips away from its
// initial value (typically to 1 on Linux / launchd on macOS).
func (l *stdioLifecycle) watchParent(ctx context.Context, cancel context.CancelCauseFunc, initial int) {
	ticker := time.NewTicker(l.parentTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := os.Getppid()
			if current != initial || current <= 1 {
				cancel(fmt.Errorf(
					"parent process gone (ppid %d -> %d)",
					initial, current,
				))
				return
			}
		}
	}
}

// parseDurationEnv parses a duration-style env var. Accepts Go duration
// literals (e.g. "5m", "90s") and bare integers (interpreted as seconds).
// The values "0", "off", "disable", "disabled", "none" disable the feature
// by returning 0.
func parseDurationEnv(name string, def time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	switch strings.ToLower(raw) {
	case "0", "off", "disable", "disabled", "none":
		return 0
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	if n, err := strconv.Atoi(raw); err == nil && n > 0 {
		return time.Duration(n) * time.Second
	}
	log.Printf("stdio lifecycle: invalid %s=%q, falling back to %s", name, raw, def)
	return def
}

// parseBoolEnv parses a boolean env var with sensible defaults. Unknown
// values fall back to def.
func parseBoolEnv(name string, def bool) bool {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	switch strings.ToLower(raw) {
	case "1", "t", "true", "y", "yes", "on":
		return true
	case "0", "f", "false", "n", "no", "off":
		return false
	}
	return def
}
