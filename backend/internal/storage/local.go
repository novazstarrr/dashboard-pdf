// internal/storage/local.go
package storage

import (
    "context"
    "io"
    "os"
    "path/filepath"
)

type LocalStorage struct {
    basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
    return &LocalStorage{
        basePath: basePath,
    }
}

func (s *LocalStorage) Save(ctx context.Context, path string, file io.Reader) error {
    fullPath := filepath.Join(s.basePath, path)

    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        return err
    }

    dst, err := os.Create(fullPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    _, err = io.Copy(dst, file)
    return err
}

func (s *LocalStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
    return os.Open(filepath.Join(s.basePath, path))
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
    return os.Remove(filepath.Join(s.basePath, path))
}

