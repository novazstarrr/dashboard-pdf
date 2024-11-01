package user

import (
    "context"
    "tech-test/backend/internal/domain"
)

type UserService interface {
    UserReader
    UserWriter
    UserAuthenticator
}

type UserReader interface {
    GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
    
    GetUserByID(ctx context.Context, id uint) (*domain.User, error)
    
    GetAllUsers(ctx context.Context) ([]domain.User, error)
}

type UserWriter interface {
    Register(ctx context.Context, user *domain.User) error
    
    UpdateUser(ctx context.Context, id uint, user *domain.User) error
    
    DeleteUser(ctx context.Context, id uint) error
}

type UserAuthenticator interface {
    Login(ctx context.Context, email, password string) (*domain.User, error)
}