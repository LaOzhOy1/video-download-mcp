package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	mu        sync.RWMutex
	downloads []string
	dbPath    = filepath.Join(".", "downloads_db.json")
)

func init() {
	_ = load()
}

func RecordDownload(path string) error {
	mu.Lock()
	defer mu.Unlock()
	// avoid duplicates
	for _, p := range downloads {
		if p == path {
			return nil
		}
	}
	downloads = append(downloads, path)
	return save()
}

func ListDownloads() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, len(downloads))
	copy(out, downloads)
	return out
}

func load() error {
	b, err := os.ReadFile(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var arr []string
	if err := json.Unmarshal(b, &arr); err != nil {
		return err
	}
	downloads = arr
	return nil
}

func save() error {
	b, err := json.MarshalIndent(downloads, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dbPath, b, 0o644)
}
