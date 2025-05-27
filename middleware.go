package main

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// IsUserLoggedIn checks if the user is currently logged in by verifying Token and Region.
func IsUserLoggedIn() bool {
	return strings.TrimSpace(Token) != "" && strings.TrimSpace(Region) != ""
}

// MCPHandlerFunc defines the standard signature for MCP tool handlers.
type MCPHandlerFunc func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)

// RequestWrapper is a middleware that provides user login verification, request timeout, and panic recovery for MCP tool handlers.
// Renamed from Wrapper.
func RequestWrapper(handler MCPHandlerFunc) server.ToolHandlerFunc {
	const handlerTimeout = DefaultAPPTimeout // Use configured application timeout
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Panic recovery for the main wrapper function execution path
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] Unhandled panic in RequestWrapper: %v\nStack trace:\n%s", r, string(debug.Stack()))
				// Note: This recovery might not be reached if panic occurs inside the goroutine and isn't caught there.
				// However, the goroutine itself has its own defer for panic recovery.
			}
		}()

		// Create a new context with timeout for the handler execution.
		handlerCtx, cancel := context.WithTimeout(ctx, handlerTimeout)
		defer cancel() // Ensure cancellation to free resources

		if !IsUserLoggedIn() {
			log.Println("[AUTH] Access denied: User not logged in or configuration is incorrect.")
			return mcp.NewToolResultText("You're not logged in or your configuration is incorrect. Please log in and complete the configuration before proceeding."), nil
		}

		resultChan := make(chan *mcp.CallToolResult, 1)
		errorChan := make(chan error, 1)

		go func() {
			// Panic recovery for the handler goroutine
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[PANIC] Panic in handler goroutine: %v\n", r)
					errorChan <- fmt.Errorf("internal server error due to panic: %v", r)
				}
			}()

			res, err := handler(handlerCtx, request) // Pass the timeout-aware context
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- res
		}()

		select {
		case <-handlerCtx.Done(): // Context timed out or was cancelled
			log.Printf("[TIMEOUT] Handler execution timed out or was cancelled: %v", handlerCtx.Err())
			// Provide a user-friendly timeout message
			return mcp.NewToolResultText("Request timed out. Please try again later."), nil
		case err := <-errorChan:
			log.Printf("[ERROR] Handler returned an error: %v", err)
			// Provide a generic error message to the user, log the specific error internally.
			// Consider mapping specific errors to more user-friendly messages if needed.
			return mcp.NewToolResultText("An internal server error occurred. Please contact the administrator or try again later."), nil
		case res := <-resultChan:
			return res, nil
		}
	}
}
