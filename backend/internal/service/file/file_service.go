package file

import (
	"context"
	"strconv"
	"tech-test/backend/internal/domain"
	"tech-test/backend/internal/repository/interfaces"
	fileInterface "tech-test/backend/internal/service/interfaces/file"
	"go.uber.org/zap"
)

type service struct {
	repo      interfaces.FileRepository
	logger    *zap.Logger
	uploadDir string
}

func NewService(repo interfaces.FileRepository, logger *zap.Logger, uploadDir string) fileInterface.Service {
	return &service{
		repo:      repo,
		logger:    logger,
		uploadDir: uploadDir,
	}
}

func (s *service) GetByID(ctx context.Context, id uint) (*domain.File, error) {
	s.logger.Debug("Getting file by ID", zap.Uint("id", id))
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByShareID(ctx context.Context, shareID string) (*domain.File, error) {
	s.logger.Debug("Getting file by share ID", zap.String("shareId", shareID))
	return s.repo.GetFileByShareID(ctx, shareID)
}

func (s *service) List(ctx context.Context) ([]domain.File, error) {
	s.logger.Debug("Listing all files")
	return s.repo.List(ctx)
}

func (s *service) GetUserFilesPaginated(ctx context.Context, userID uint, page, pageSize int) ([]domain.File, int64, error) {
	s.logger.Debug("Getting paginated user files",
		zap.Uint("userID", userID),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))
	return s.repo.GetUserFilesPaginated(ctx, userID, page, pageSize)
}

func (s *service) SearchFiles(ctx context.Context, userID uint, searchTerm string) ([]domain.File, error) {
	s.logger.Debug("Searching files",
		zap.Uint("userID", userID),
		zap.String("searchTerm", searchTerm))
	return s.repo.SearchFiles(ctx, userID, searchTerm)
}

func (s *service) Upload(ctx context.Context, file *domain.File) error {
	s.logger.Debug("Uploading file",
		zap.String("name", file.Name),
		zap.Int64("size", file.Size))
	return s.repo.Create(ctx, file)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	s.logger.Debug("Deleting file", zap.Uint("id", id))
	return s.repo.Delete(ctx, id)
}

func (s *service) UpdateShareableID(ctx context.Context, fileID string, shareableID string) error {
	s.logger.Debug("Updating shareable ID",
		zap.String("fileID", fileID),
		zap.String("shareableID", shareableID))
	
	id, err := strconv.ParseUint(fileID, 10, 32)
	if err != nil {
		s.logger.Error("Invalid file ID format", 
			zap.String("fileID", fileID),
			zap.Error(err))
		return domain.NewAPIError(
			400,
			domain.ErrCodeInvalidInput,
			"Invalid file ID format",
			err,
		)
	}

	return s.repo.UpdateShareableID(ctx, uint(id), shareableID)
}
