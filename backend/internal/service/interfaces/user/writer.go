package user

import (
	"context"
	"tech-test/backend/internal/domain"
)

type Writer interface {
	Register(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, id uint, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error
}
