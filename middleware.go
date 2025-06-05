package main

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// RequestWrapper wraps MCP tool handlers with authentication, timeout, and error recovery middleware.
func RequestWrapper(handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Add timeout to context
		ctx, cancel := context.WithTimeout(ctx, DefaultAPPTimeout)
		defer cancel()

		// Defer panic recovery
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERROR] [RequestWrapper] Panic recovered in handler for tool '%s': %v\n%s",
					request.Params.Name, r, debug.Stack())
			}
		}()

		// Validate authentication
		if err := validateAuth(); err != nil {
			log.Printf("[ERROR] [RequestWrapper] Authentication failed for tool '%s': %v", request.Params.Name, err)
			return mcp.NewToolResultText(fmt.Sprintf("Authentication failed: %v", err)), nil
		}

		// Log request start
		log.Printf("[INFO] [RequestWrapper] Processing tool request: %s", request.Params.Name)

		// Execute the actual handler
		result, err := handler(ctx, request)

		// Log completion
		if err != nil {
			log.Printf("[ERROR] [RequestWrapper] Tool '%s' failed: %v", request.Params.Name, err)
		} else {
			log.Printf("[INFO] [RequestWrapper] Tool '%s' completed successfully", request.Params.Name)
		}

		return result, err
	}
}

// validateAuth checks if the required authentication credentials are available.
func validateAuth() error {
	if strings.TrimSpace(Token) == "" {
		return fmt.Errorf("token is required but not provided")
	}
	if strings.TrimSpace(Region) == "" {
		return fmt.Errorf("region is required but not provided")
	}
	return nil
}
