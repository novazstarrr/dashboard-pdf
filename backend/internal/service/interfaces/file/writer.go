package file

import (
	"context"
	"tech-test/backend/internal/domain"
)

type Writer interface {
	Upload(ctx context.Context, file *domain.File) error
	Delete(ctx context.Context, id uint) error
	UpdateShareableID(ctx context.Context, fileID string, shareableID string) error
} 