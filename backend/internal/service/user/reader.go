package user

import (
	"context"
	"go.uber.org/zap"
	"tech-test/backend/internal/domain"
	"tech-test/backend/internal/repository/interfaces"
	userInterface "tech-test/backend/internal/service/interfaces/user"
)

type reader struct {
	repo   interfaces.UserRepository
	logger *zap.Logger
}

func NewReader(repo interfaces.UserRepository, logger *zap.Logger) userInterface.Reader {
	return &reader{
		repo:   repo,
		logger: logger,
	}
}

func (r *reader) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	r.logger.Debug("Getting user by ID", zap.Uint("id", id))
	return r.repo.GetByID(ctx, id)
}

func (r *reader) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.logger.Debug("Getting user by email", zap.String("email", email))
	return r.repo.GetByEmail(ctx, email)
}

func (r *reader) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	r.logger.Debug("Getting all users")
	return r.repo.GetAll(ctx)
}
