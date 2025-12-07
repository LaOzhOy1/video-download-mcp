package usecase

import (
    "context"
    "errors"
    "fmt"
    "bufio"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// DownloadVideo downloads the content from the given URL and saves it to targetPath.
// It ensures the directory exists and returns the absolute path to the saved file.
func DownloadVideo(ctx context.Context, url string, targetPath string) (string, error) {
    // Basic validation
    if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
        return "", errors.New("url must start with http:// or https://")
    }

    // Ensure directory exists
    dir := filepath.Dir(targetPath)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return "", fmt.Errorf("failed to create directory: %w", err)
    }

    // Create HTTP client with timeout
    client := &http.Client{Timeout: 60 * time.Second}

    // Build request with context
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return "", fmt.Errorf("failed to build request: %w", err)
    }

    // Execute request
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Create (or truncate) the target file
    out, err := os.Create(targetPath)
    if err != nil {
        return "", fmt.Errorf("failed to create file: %w", err)
    }
    defer out.Close()

    // Stream copy response body to file
    if _, err := io.Copy(out, resp.Body); err != nil {
        return "", fmt.Errorf("failed to write file: %w", err)
    }

    // Compute and return absolute path
    abs, err := filepath.Abs(targetPath)
    if err != nil {
        return "", fmt.Errorf("failed to compute absolute path: %w", err)
    }
    return abs, nil
}

// DownloadVideoWithProgress downloads the URL to targetPath and invokes onProgress with (written, total)
func DownloadVideoWithProgress(ctx context.Context, url string, targetPath string, onProgress func(written int64, total int64)) (string, error) {
    if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
        return "", errors.New("url must start with http:// or https://")
    }
    dir := filepath.Dir(targetPath)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return "", fmt.Errorf("failed to create directory: %w", err)
    }

    client := &http.Client{Timeout: 60 * time.Second}
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return "", fmt.Errorf("failed to build request: %w", err)
    }
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    out, err := os.Create(targetPath)
    if err != nil {
        return "", fmt.Errorf("failed to create file: %w", err)
    }
    defer out.Close()

    var written int64
    total := resp.ContentLength
    reader := bufio.NewReaderSize(resp.Body, 64*1024)
    buf := make([]byte, 64*1024)
    for {
        n, rErr := reader.Read(buf)
        if n > 0 {
            if _, wErr := out.Write(buf[:n]); wErr != nil {
                return "", fmt.Errorf("failed to write file: %w", wErr)
            }
            written += int64(n)
            if onProgress != nil {
                onProgress(written, total)
            }
        }
        if rErr != nil {
            if rErr == io.EOF {
                break
            }
            return "", fmt.Errorf("read error: %w", rErr)
        }
    }

    abs, err := filepath.Abs(targetPath)
    if err != nil {
        return "", fmt.Errorf("failed to compute absolute path: %w", err)
    }
    return abs, nil
}
