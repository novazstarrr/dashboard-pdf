// internal/repository/sqlite/user.go
package sqlite

import (
    "context"
    "gorm.io/gorm"
    "tech-test/backend/internal/domain"
    "tech-test/backend/internal/repository/interfaces"
)

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
    return &userRepository{
        db: db,
    }
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
        return domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to create user",
            err,
        )
    }
    return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    var user domain.User
    if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to get user",
            err,
        )
    }
    return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to get user by email",
            err,
        )
    }
    return &user, nil
}

func (r *userRepository) Update(ctx context.Context, id uint, user *domain.User) error {
    result := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(user)
    if result.Error != nil {
        return domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to update user",
            result.Error,
        )
    }
    if result.RowsAffected == 0 {
        return domain.ErrUserNotFound
    }
    return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
    result := r.db.WithContext(ctx).Delete(&domain.User{}, id)
    if result.Error != nil {
        return domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to delete user",
            result.Error,
        )
    }
    if result.RowsAffected == 0 {
        return domain.ErrUserNotFound
    }
    return nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
    var users []domain.User
    if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
        return nil, domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to get all users",
            err,
        )
    }
    return users, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to find user by email",
            err,
        )
    }
    return &user, nil
}
