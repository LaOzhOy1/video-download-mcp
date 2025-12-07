package prompt

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// registerPrompts defines MCP prompts that Cherry Studio can call to guide file downloads
func RegisterPrompts(s *mcpserver.MCPServer) {
    // Prompt: download_file_prompt – guides clients to call the tool with required args
	p := mcp.NewPrompt(
		"download_file_prompt",
		mcp.WithPromptDescription("Guide to download a file by URL to target directory with filename"),
		mcp.WithArgument("url", mcp.ArgumentDescription("The video/file URL (HTTP/HTTPS)"), mcp.RequiredArgument()),
		mcp.WithArgument("save_dir", mcp.ArgumentDescription("Destination directory to save the file"), mcp.RequiredArgument()),
		mcp.WithArgument("filename", mcp.ArgumentDescription("Target filename to use for saving"), mcp.RequiredArgument()),
	)

	s.AddPrompt(p, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments
        url := fmt.Sprintf("%v", args["url"])      // simple extraction; validated by client
        dir := fmt.Sprintf("%v", args["save_dir"])  // required by prompt
        name := fmt.Sprintf("%v", args["filename"]) // required by prompt

		title := "Download File Instruction"
        // Compose assistant message instructing client to invoke the tool with these args
		msg := mcp.NewPromptMessage(
			mcp.RoleAssistant,
			mcp.NewTextContent(
				fmt.Sprintf(
					"Use tool 'download_video_file' with arguments: url=%s, save_dir=%s, filename=%s. This will download the file to the destination and return the absolute path.",
					url, dir, name,
				),
			),
		)

		return mcp.NewGetPromptResult(title, []mcp.PromptMessage{msg}), nil
	})
}
