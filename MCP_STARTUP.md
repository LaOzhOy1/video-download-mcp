# MCP Server Startup Guide
The video download assistant MCP supports downloading and saving links locally.


## MCP Server Startup (Studio & SSE)

This project provides an MCP server that supports both `stdio` (Studio) and `sse` transports. Use the following commands on Windows to start the server.

## Studio (stdio)

Start the server with stdio transport for use in Studio-compatible clients:

```powershell
go run .\cmd\mcp-server\main.go --transport=stdio
```

## SSE

Start the server with SSE transport, exposing an HTTP endpoint (default port 3000):

```powershell
go run .\cmd\mcp-server\main.go --transport=sse --port 3000
```

## Bash (Unix/macOS/Linux)

Start the server using bash-compatible commands on Unix-like systems.

### Studio (stdio)

```bash
go run ./cmd/mcp-server/main.go --transport=stdio
```

### SSE

```bash
go run ./cmd/mcp-server/main.go --transport=sse --port 3000
```

## Tool: download_video_file

Parameters (all required):

- `url`: Video file URL (HTTP/HTTPS)
- `save_dir`: Destination directory (will be created if missing)
- `filename`: Target file name (e.g., `myvideo.mp4`)

Returns: Absolute path of the saved video file.
