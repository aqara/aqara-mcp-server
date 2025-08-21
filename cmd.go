package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// rootCmd is the base command when no subcommands are invoked.
	rootCmd = &cobra.Command{
		Use:     "aqara-mcp",
		Short:   "Aqara MCP Server",
		Long:    `Aqara MCP Server - Manage your Aqara devices and services through the MCP protocol.`,
		Version: Version,
		Run:     runDefaultServer,
	}

	// runCmd represents the base for 'run' subcommands.
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run related server commands",
		Long:  "Execute various runtime operations for the Aqara MCP server, such as starting it with different transport protocols.",
	}

	// stdioCmd represents the command to start the server with Stdio transport.
	stdioCmd = &cobra.Command{
		Use:   "stdio",
		Short: "Start server with stdio transport",
		Long:  `Starts the Aqara MCP server communicating via standard input/output (stdio) streams, typically using JSON-RPC messages.`,
		Run:   runStdioServer,
	}

	// httpCmd represents the command to start the server with StreamableHTTP transport.
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Start server with streamablehttp transport",
		Long:  `Starts the Aqara MCP server communicating via StreamableHTTP, typically using JSON-RPC messages over HTTP.`,
		Run:   runStreamableHTTPServer,
	}

	// sseCmd represents the command to start the server with SSE transport.
	sseCmd = &cobra.Command{
		Use:   "sse",
		Short: "Start server with SSE transport",
		Long:  `Starts the Aqara MCP server communicating via Server-Sent Events (SSE), typically using JSON-RPC messages over HTTP streams.`,
		Run:   runSSEServer,
	}
)

// runStdioServer is the function executed by the stdioCmd.
func runStdioServer(cmd *cobra.Command, args []string) {
	serverStart("stdio")
}

// runStreamableHTTPServer is the function executed by the httpCmd.
func runStreamableHTTPServer(cmd *cobra.Command, args []string) {
	serverStart("http")
}

// runSSEServer is the function executed by the sseCmd.
func runSSEServer(cmd *cobra.Command, args []string) {
	serverStart("sse")
}

// runDefaultServer is the default function executed when no subcommand is specified.
func runDefaultServer(cmd *cobra.Command, args []string) {
	log.Println("Starting Aqara MCP Server with default http transport...")
	serverStart("http")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func init() {
	// Persistent flags are global for the application.
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")

	// Flags specific to SSE server, but defined as persistent for potential broader use or simplification.
	// Alternatively, they could be local flags for sseCmd.
	rootCmd.PersistentFlags().String("host", "0.0.0.0", "The host address for the Streamable-HTTP server.")
	rootCmd.PersistentFlags().String("port", "8000", "The port for the SSE server.")

	// Bind these flags to viper keys for configuration management.
	// Errors are ignored here as per original, but consider handling them.
	_ = viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	_ = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	// viper.AutomaticEnv() // Optionally read from environment variables

	// Add subcommands.
	runCmd.AddCommand(stdioCmd)
	runCmd.AddCommand(httpCmd)
	runCmd.AddCommand(sseCmd)
	rootCmd.AddCommand(runCmd)
}
