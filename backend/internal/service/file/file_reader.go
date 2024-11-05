package file

import (
	"context"
	"go.uber.org/zap"
	"tech-test/backend/internal/domain"
	"tech-test/backend/internal/repository/interfaces"
	fileInterface "tech-test/backend/internal/service/interfaces/file"
)

type reader struct {
	repo   interfaces.FileRepository
	logger *zap.Logger
}

func NewReader(repo interfaces.FileRepository, logger *zap.Logger) fileInterface.Reader {
	return &reader{
		repo:   repo,
		logger: logger,
	}
}

func (r *reader) GetByID(ctx context.Context, id uint) (*domain.File, error) {
	r.logger.Debug("Getting file by ID", zap.Uint("id", id))
	return r.repo.GetByID(ctx, id)
}

func (r *reader) GetByShareID(ctx context.Context, shareID string) (*domain.File, error) {
	r.logger.Debug("Getting file by share ID", zap.String("shareId", shareID))
	return r.repo.GetFileByShareID(ctx, shareID)
}

func (r *reader) List(ctx context.Context) ([]domain.File, error) {
	r.logger.Debug("Listing all files")
	return r.repo.List(ctx)
}

func (r *reader) GetUserFilesPaginated(ctx context.Context, userID uint, page, pageSize int) ([]domain.File, int64, error) {
	r.logger.Debug("Getting paginated user files",
		zap.Uint("userID", userID),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
	)
	return r.repo.GetUserFilesPaginated(ctx, userID, page, pageSize)
}

func (r *reader) SearchFiles(ctx context.Context, userID uint, searchTerm string) ([]domain.File, error) {
	r.logger.Debug("Searching files",
		zap.Uint("userID", userID),
		zap.String("searchTerm", searchTerm),
	)
	return r.repo.SearchFiles(ctx, userID, searchTerm)
}
