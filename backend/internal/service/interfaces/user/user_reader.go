package user

import (
	"context"
	"tech-test/backend/internal/domain"
)

type Reader interface {
	GetUser(ctx context.Context, id uint) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
} 