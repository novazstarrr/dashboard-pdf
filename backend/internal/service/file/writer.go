package file

import (
	"context"
	"go.uber.org/zap"
	"os"
	"strconv"
	"tech-test/backend/internal/domain"
	"tech-test/backend/internal/repository/interfaces"
	fileInterface "tech-test/backend/internal/service/interfaces/file"
)

type writer struct {
	repo      interfaces.FileRepository
	logger    *zap.Logger
	uploadDir string
}

func NewWriter(repo interfaces.FileRepository, logger *zap.Logger, uploadDir string) fileInterface.Writer {
	return &writer{
		repo:      repo,
		logger:    logger,
		uploadDir: uploadDir,
	}
}

func (w *writer) Upload(ctx context.Context, file *domain.File) error {
	w.logger.Debug("Uploading file",
		zap.String("name", file.Name),
		zap.Int64("size", file.Size),
	)
	
	if err := w.repo.Create(ctx, file); err != nil {
		w.logger.Error("Failed to create file record",
			zap.Error(err),
			zap.String("name", file.Name),
		)
		return err
	}
	
	return nil
}

func (w *writer) Delete(ctx context.Context, id uint) error {
	w.logger.Debug("Deleting file", zap.Uint("id", id))
	
	file, err := w.repo.GetByID(ctx, id)
	if err != nil {
		w.logger.Error("Failed to get file for deletion", 
			zap.Error(err),
			zap.Uint("id", id),
		)
		return err
	}

	if _, err := os.Stat(file.Path); !os.IsNotExist(err) {
		if err := os.Remove(file.Path); err != nil {
			w.logger.Error("Failed to delete physical file",
				zap.Error(err),
				zap.String("path", file.Path),
			)
			return domain.NewAPIError(
				500,
				domain.ErrCodeInternal,
				"Failed to delete physical file",
				err,
			)
		}
	}

	if err := w.repo.Delete(ctx, id); err != nil {
		w.logger.Error("Failed to delete file record",
			zap.Error(err),
			zap.Uint("id", id),
		)
		return err
	}

	return nil
}

func (w *writer) UpdateShareableID(ctx context.Context, fileID string, shareableID string) error {
	w.logger.Debug("Updating shareable ID",
		zap.String("fileID", fileID),
		zap.String("shareableID", shareableID),
	)
	
	id, err := strconv.ParseUint(fileID, 10, 32)
	if err != nil {
		w.logger.Error("Invalid file ID format",
			zap.Error(err),
			zap.String("fileID", fileID),
		)
		return domain.NewAPIError(
			400,
			domain.ErrCodeInvalidInput,
			"Invalid file ID format",
			err,
		)
	}

	if err := w.repo.UpdateShareableID(ctx, uint(id), shareableID); err != nil {
		w.logger.Error("Failed to update shareable ID",
			zap.Error(err),
			zap.String("fileID", fileID),
			zap.String("shareableID", shareableID),
		)
		return err
	}

	return nil
}
