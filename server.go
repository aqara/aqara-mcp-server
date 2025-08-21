package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
)

func httpContext(ctx context.Context, r *http.Request) context.Context {
	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		log.Println("Authorization header present")
	} else {
		log.Println("No authorization header provided")
	}

	if after, ok := strings.CutPrefix(authorization, "Bearer "); ok {
		authorization = after
	}

	var apiBase, apiKey string

	if authorization != "" {
		apiKey = authorization
		apiBase = r.Header.Get("BaseUrl")
		if apiBase == "" {
			apiBase = CloudAPIBase
		}
	}

	if apiBase != "" && apiKey != "" {
		ctx = context.WithValue(ctx, CloudCtx{}, CloudCtx{
			ApiBase: apiBase,
			ApiKey:  apiKey,
		})
	}

	return ctx
}

func serverStart(transport string) {
	hooks := &server.Hooks{}

	mcpServer := server.NewMCPServer(
		"Aqara MCP Server",
		Version,
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(hooks),
		BindTools(),
	)

	switch transport {
	case "http":
		host, port := "", ""

		if err := viper.UnmarshalKey("port", &port); err != nil {
			log.Printf("failed to unmarshal port from config: %v. Using default or expecting environment override.\n", err)
		}
		if err := viper.UnmarshalKey("host", &host); err != nil {
			log.Printf("failed to unmarshal host from config: %v. Using default or expecting environment override.\n", err)
		}

		httpServer := server.NewStreamableHTTPServer(mcpServer,
			server.WithEndpointPath("/echo/mcp"),
			server.WithHTTPContextFunc(httpContext),
			server.WithStateLess(true),
		)

		if httpServer == nil {
			log.Fatalln("HTTP server initialization failed")
		}

		baseURL := host + ":" + port
		log.Printf("HTTP server listening on %s", baseURL)
		// Start the server
		if err := httpServer.Start(baseURL); err != nil {
			log.Fatalf("HTTP Server error: %v", err)
		}
	case "sse":
		host, port := "", ""

		if err := viper.UnmarshalKey("port", &port); err != nil {
			log.Printf("failed to unmarshal port from config: %v. Using default or expecting environment override.\n", err)
		}
		if err := viper.UnmarshalKey("host", &host); err != nil {
			log.Printf("failed to unmarshal host from config: %v. Using default or expecting environment override.\n", err)
		}

		sseServer := server.NewSSEServer(mcpServer)

		if sseServer == nil {
			log.Fatalln("SSE server initialization failed")
		}

		baseURL := host + ":" + port
		log.Printf("SSE server listening on %s", baseURL)
		// Start the server
		if err := sseServer.Start(baseURL); err != nil {
			log.Fatalf("SSE Server error: %v", err)
		}
	default:
		log.Println("Starting server with stdio transport.")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Stdio Server error: %v", err)
		}
	}
}
