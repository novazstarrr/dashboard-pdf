package file

import (
	"context"
	"tech-test/backend/internal/domain"
)

type Reader interface {
	GetByID(ctx context.Context, id uint) (*domain.File, error)
	GetByShareID(ctx context.Context, shareID string) (*domain.File, error)
	List(ctx context.Context) ([]domain.File, error)
	GetUserFilesPaginated(ctx context.Context, userID uint, page, pageSize int) ([]domain.File, int64, error)
	SearchFiles(ctx context.Context, userID uint, searchTerm string) ([]domain.File, error)
} 