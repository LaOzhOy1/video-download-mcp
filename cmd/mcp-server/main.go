package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "video-download-mcp/internal/server"
)

// main is the entry point that starts the MCP server with either stdio or SSE transport
func main() {
    // Parse CLI flags for transport and SSE port
    transport := flag.String("transport", "stdio", "MCP transport: stdio | sse")
    port := flag.Int("port", 3000, "Port for SSE server when transport is 'sse'")
    flag.Parse()

    // Create the application MCP server instance
    srv := server.NewMCPServer()

    // Handle OS signals for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Start according to the selected transport
    switch *transport {
    case "stdio":
        log.Println("Starting MCP server with stdio transport...")
        go func() {
            if err := server.RunMCPServerWithStdio(srv); err != nil {
                log.Printf("MCP stdio server error: %v", err)
            }
        }()
    case "sse":
        log.Printf("Starting MCP server with SSE transport on port %d...", *port)
        go func() {
            if err := server.RunMCPServerWithSSE(srv, *port); err != nil {
                log.Printf("MCP SSE server error: %v", err)
            }
        }()
    default:
        log.Fatalf("unknown transport: %s (expected 'stdio' or 'sse')", *transport)
    }

    // Block until shutdown signal is received
    <-sigChan
    fmt.Println("Shutting down MCP server...")
}

