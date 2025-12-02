# MCP 服务启动指南（Studio 与 SSE）

该项目提供一个 MCP 服务，支持 `stdio`（Studio）与 `sse` 两种传输方式。服务包含工具 `download_video_file`，用于将视频链接下载到指定目录并以指定文件名保存，返回保存后的绝对路径。

## Windows（PowerShell）

### Studio（stdio）启动

```powershell
go run .\cmd\mcp-server\main.go --transport=stdio
```

### SSE 启动（默认端口 3000）

```powershell
go run .\cmd\mcp-server\main.go --transport=sse --port 3000
```

## Bash（Unix/macOS/Linux）

### Studio（stdio）启动

```bash
go run ./cmd/mcp-server/main.go --transport=stdio
```

### SSE 启动（默认端口 3000）

```bash
go run ./cmd/mcp-server/main.go --transport=sse --port 3000
```

## 工具：download_video_file

- 参数（均为必填）：
  - `url`：视频文件链接（HTTP/HTTPS）
  - `save_dir`：保存目录（不存在会自动创建）
  - `filename`：目标文件名（例如 `myvideo.mp4`）

- 返回：保存的视频文件的绝对路径

## 使用说明

- 选择 `--transport=stdio` 可用于 Studio 类客户端对接；选择 `--transport=sse` 将在指定端口开启 SSE HTTP 服务。
- 当网络环境存在代理限制时，请确保 `git` 与 `Go` 的代理配置正确，以便拉取 `github.com/mark3labs/mcp-go` 依赖。
