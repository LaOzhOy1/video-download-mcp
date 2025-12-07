package server

import (
    "context"
    "fmt"
    "log"
    "time"

    mcpserver "github.com/mark3labs/mcp-go/server"

    "video-download-mcp/internal/resources"
    "video-download-mcp/internal/prompt"
    "video-download-mcp/internal/tools"
)

// MCPServer wraps the underlying mcp-go server instance and tool registration
type MCPServer struct {
	server *mcpserver.MCPServer
}

// NewMCPServer constructs a new MCP server and registers all tools
func NewMCPServer() *MCPServer {
	// Create MCP server with basic middlewares
	srv := mcpserver.NewMCPServer(
		"video-download",
		"1.0.0",
		mcpserver.WithLogging(),
		mcpserver.WithRecovery(),
	)

	// Register tools
	tools.RegisterTools(srv)

    // Register prompts
    prompt.RegisterPrompts(srv)

    // Register resources
    resources.RegisterResources(srv)

	return &MCPServer{server: srv}
}

// Underlying exposes the raw mcp-go server for transports
func (s *MCPServer) Underlying() *mcpserver.MCPServer {
	return s.server
}

// RunMCPServerWithStdio starts the MCP server using stdio transport
func RunMCPServerWithStdio(s *MCPServer) error {
	return mcpserver.ServeStdio(s.server)
}

// RunMCPServerWithSSE starts the MCP server using SSE transport on the given port
func RunMCPServerWithSSE(s *MCPServer, port int) error {
	sse := mcpserver.NewSSEServer(s.server)
	addr := fmt.Sprintf(":%d", port)
	return sse.Start(addr)
}

// ShutdownSSE attempts to gracefully stop the SSE server with a timeout
// NOTE: This requires a reference to the SSE server. For simplicity, callers can manage lifecycle.
func ShutdownSSE(sse *mcpserver.SSEServer) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sse.Shutdown(ctx); err != nil {
		log.Printf("Error during SSE shutdown: %v", err)
	}
}
