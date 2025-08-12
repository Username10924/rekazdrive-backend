package storage

import (
	"os"
	"path/filepath"
	"strings"
)

type LocalBackend struct {
	BasePath string
}

func NewLocalBackend(basePath string) *LocalBackend {
	_ = os.MkdirAll(basePath, 0755)
	return &LocalBackend{BasePath: basePath}
}

func (l *LocalBackend) pathFor(id string) string {
	id = strings.ReplaceAll(id, "/", "_") // replace slashes to avoid directory traversal
	if id == "" {
		id = "someblob"
	}
	return filepath.Join(l.BasePath, id) // example output: /path/to/base/id
}

func (l *LocalBackend) Save(id string, data []byte) error {
	p := l.pathFor(id)
	dir := filepath.Dir(p) // path without the file name to ensure the directory exists
	// Create the directory if it does not exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

func (l *LocalBackend) Load(id string) ([]byte, error) {
	p := l.pathFor(id)
	return os.ReadFile(p)
}

func (l *LocalBackend) Delete(id string) error {
	p := l.pathFor(id)
	return os.Remove(p)
}