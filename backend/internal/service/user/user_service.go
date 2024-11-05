package user

import (
	"context"
	"tech-test/backend/internal/domain"
	"tech-test/backend/internal/repository/interfaces"
	userInterface "tech-test/backend/internal/service/interfaces/user"
	"tech-test/backend/internal/utils"
	"go.uber.org/zap"
	"errors"
)

type Service struct {
	repo   interfaces.UserRepository
	logger *zap.Logger
}

var _ userInterface.UserService = (*Service)(nil)

func NewService(repo interfaces.UserRepository, logger *zap.Logger) *Service {
	if repo == nil {
		panic("repo cannot be nil")
	}
	if logger == nil {
		panic("logger cannot be nil")
	}
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	s.logger.Debug("Getting user by ID", zap.Uint("id", id))
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	s.logger.Debug("Getting user by email", zap.String("email", email))
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	s.logger.Debug("Getting all users")
	return s.repo.GetAll(ctx)
}

func (s *Service) Register(ctx context.Context, user *domain.User) error {
	s.logger.Debug("Registering new user", zap.String("email", user.Email))
	
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return err
	}
	user.Password = hashedPassword

	return s.repo.Create(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, id uint, user *domain.User) error {
	s.logger.Debug("Updating user", zap.Uint("id", id))
	return s.repo.Update(ctx, id, user)
}

func (s *Service) DeleteUser(ctx context.Context, id uint) error {
	s.logger.Debug("Deleting user", zap.Uint("id", id))
	return s.repo.Delete(ctx, id)
}

func (s *Service) Login(ctx context.Context, email, password string) (*domain.User, error) {
	s.logger.Debug("Attempting login", zap.String("email", email))
	
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
