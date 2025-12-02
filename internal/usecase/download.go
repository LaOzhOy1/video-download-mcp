package usecase

import (
    "context"
    "errors"
    "fmt"
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

