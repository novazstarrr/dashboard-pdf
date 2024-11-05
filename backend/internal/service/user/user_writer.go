package user

import (
	"context"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"tech-test/backend/internal/domain"
	"tech-test/backend/internal/repository/interfaces"
	userInterface "tech-test/backend/internal/service/interfaces/user"
)

type writer struct {
	repo   interfaces.UserRepository
	logger *zap.Logger
}

func NewWriter(repo interfaces.UserRepository, logger *zap.Logger) userInterface.Writer {
	return &writer{
		repo:   repo,
		logger: logger,
	}
}

func (w *writer) Register(ctx context.Context, user *domain.User) error {
	w.logger.Debug("Registering new user", zap.String("email", user.Email))
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.logger.Error("Failed to hash password", zap.Error(err))
		return err
	}
	
	user.Password = string(hashedPassword)
	return w.repo.Create(ctx, user)
}

func (w *writer) UpdateUser(ctx context.Context, id uint, user *domain.User) error {
	w.logger.Debug("Updating user",
		zap.Uint("id", id),
		zap.String("email", user.Email),
	)
	
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			w.logger.Error("Failed to hash password during update", zap.Error(err))
			return err
		}
		user.Password = string(hashedPassword)
	}
	
	return w.repo.Update(ctx, id, user)
}

func (w *writer) DeleteUser(ctx context.Context, id uint) error {
	w.logger.Debug("Deleting user", zap.Uint("id", id))
	return w.repo.Delete(ctx, id)
}
