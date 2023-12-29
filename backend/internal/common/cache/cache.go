package cache

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/realfabecker/photos/internal/core/ports"
)

type FileCacheHandler struct{}

func NewFileCache() ports.CacheHandler {
	return &FileCacheHandler{}
}

func (c *FileCacheHandler) Get(key string) ([]byte, error) {
	d, err := os.ReadFile(filepath.Join(os.TempDir(), key+".json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return d, nil
}

func (c *FileCacheHandler) Set(key string, data []byte) error {
	return os.WriteFile(filepath.Join(os.TempDir(), key+".json"), data, 0644)
}
