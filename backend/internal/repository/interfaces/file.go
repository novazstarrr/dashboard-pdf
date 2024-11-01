package interfaces

import (
	"context"
	"tech-test/backend/internal/domain"
)

type FileRepository interface {
    Create(ctx context.Context, file *domain.File) error
    GetByID(ctx context.Context, id uint) (*domain.File, error)
    GetByUserID(ctx context.Context, userID uint) ([]domain.File, error)
    GetUserFilesPaginated(ctx context.Context, userID uint, page, pageSize int) ([]domain.File, int64, error)
    List(ctx context.Context) ([]domain.File, error)
    Delete(ctx context.Context, id uint) error
    SearchFiles(ctx context.Context, userID uint, searchTerm string) ([]domain.File, error)
    UpdateShareableID(ctx context.Context, fileID uint, shareableID string) error
    GetFileByShareID(ctx context.Context, shareID string) (*domain.File, error)
}