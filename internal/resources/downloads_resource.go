package resources

import (
    "context"
    "encoding/json"

    "github.com/mark3labs/mcp-go/mcp"
    mcpserver "github.com/mark3labs/mcp-go/server"

    "video-download-mcp/internal/storage"
)

// RegisterResources registers server resources for clients to read
func RegisterResources(s *mcpserver.MCPServer) {
    // Resource: list of downloaded file paths
    res := mcp.NewResource(
        "downloads://list",
        "Downloaded Files",
        mcp.WithResourceDescription("List of all downloaded files' absolute paths"),
        mcp.WithMIMEType("application/json"),
    )

    s.AddResource(res, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
        paths := storage.ListDownloads()
        data, err := json.Marshal(paths)
        if err != nil {
            return nil, err
        }
        return []mcp.ResourceContents{
            mcp.TextResourceContents{
                URI:      "downloads://list",
                MIMEType: "application/json",
                Text:     string(data),
            },
        }, nil
    })
}

