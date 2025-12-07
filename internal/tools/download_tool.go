package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"video-download-mcp/internal/storage"
	"video-download-mcp/internal/usecase"
)

// RegisterTools registers all MCP tools for the service
func RegisterTools(s *mcpserver.MCPServer) {
	// Tool: download_video_file
	// Required params: url, save_dir, filename
	downloadTool := mcp.NewTool(
		"download_video_file",
		mcp.WithDescription("Download a video from URL and save to target directory with a given filename"),
		mcp.WithString("url", mcp.Description("The video file URL (HTTP/HTTPS)"), mcp.Required()),
		mcp.WithString("save_dir", mcp.Description("Directory where the video will be saved"), mcp.Required()),
		mcp.WithString("filename", mcp.Description("Target filename (without path)"), mcp.Required()),
	)
	s.AddTool(downloadTool, handleDownloadVideoFile)
}

// handleDownloadVideoFile executes the download and returns the absolute saved path
func handleDownloadVideoFile(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Type-safe argument access – returns typed errors if wrong or missing
	url, err := req.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	saveDir, err := req.RequireString("save_dir")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fileName, err := req.RequireString("filename")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Sanitize and compute target path
	targetPath := filepath.Join(saveDir, fileName)

	// Perform download via use case (with progress callback placeholder)
	savedPath, derr := usecase.DownloadVideoWithProgress(ctx, url, targetPath, func(written int64, total int64) {})
	if derr != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to download video: %v", derr)), nil
	}

	// Record successfully downloaded path
	_ = storage.RecordDownload(savedPath)

	// Return absolute path
	return mcp.NewToolResultText(savedPath), nil
}
