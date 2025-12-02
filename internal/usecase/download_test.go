package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	gomonkey "github.com/agiledragon/gomonkey/v2"
)

func TestDownloadVideo_Success(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "video.mp4")

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&http.Client{}), "Do",
		func(_ *http.Client, req *http.Request) (*http.Response, error) {
			body := io.NopCloser(bytes.NewReader([]byte("test-bytes")))
			return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
		})
	defer patches.Reset()

	abs, err := DownloadVideo(context.Background(), "https://example.com/video.mp4", target)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if abs == "" {
		t.Fatalf("expected absolute path, got empty string")
	}
	if _, err := os.Stat(abs); err != nil {
		t.Fatalf("expected file to exist, stat error: %v", err)
	}
}

func TestDownloadVideo_StatusNotOK(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "bad.mp4")

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&http.Client{}), "Do",
		func(_ *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewReader(nil))}, nil
		})
	defer patches.Reset()

	_, err := DownloadVideo(context.Background(), "https://example.com/video.mp4", target)
	if err == nil {
		t.Fatalf("expected error when status != 200")
	}
}

func TestDownloadVideo_InvalidURL(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "x.mp4")

	_, err := DownloadVideo(context.Background(), "ftp://example.com/video.mp4", target)
	if err == nil {
		t.Fatalf("expected error for non-http(s) url")
	}
}

func TestDownloadVideo_MkdirFail(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "y.mp4")

	patches := gomonkey.ApplyFunc(os.MkdirAll, func(path string, perm os.FileMode) error { return errors.New("mkdir fail") })
	defer patches.Reset()

	_, err := DownloadVideo(context.Background(), "https://example.com/video.mp4", target)
	if err == nil {
		t.Fatalf("expected error when mkdir fails")
	}
}
