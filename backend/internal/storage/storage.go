// internal/storage/storage.go
package storage

import (
    "context"
    "io"
)

type Storage interface {
    Save(ctx context.Context, path string, file io.Reader) error
    Get(ctx context.Context, path string) (io.ReadCloser, error)
    Delete(ctx context.Context, path string) error
}

