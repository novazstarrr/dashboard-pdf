// internal/repository/memory/user.go
package memory

import (
    "context"
    "sync"
    "tech-test/backend/internal/domain"
    "tech-test/backend/internal/repository/interfaces"
)

var _ interfaces.UserRepository = (*userRepository)(nil)

type userRepository struct {
    users  map[uint]*domain.User
    mutex  sync.RWMutex
    nextID uint
}

func NewUserRepository() interfaces.UserRepository {
    return &userRepository{
        users:  make(map[uint]*domain.User),
        nextID: 1,
    }
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    for _, existingUser := range r.users {
        if existingUser.Email == user.Email {
            return domain.ErrDuplicateEmail
        }
    }

    user.ID = r.nextID
    r.users[r.nextID] = user
    r.nextID++
    return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    user, exists := r.users[id]
    if !exists {
        return nil, domain.ErrUserNotFound
    }
    return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    for _, user := range r.users {
        if user.Email == email {
            return user, nil
        }
    }
    return nil, domain.ErrUserNotFound
}

func (r *userRepository) Update(ctx context.Context, id uint, user *domain.User) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.users[id]; !exists {
        return domain.ErrUserNotFound
    }

    user.ID = id  // Ensure ID remains the same
    r.users[id] = user
    return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.users[id]; !exists {
        return domain.ErrUserNotFound
    }

    delete(r.users, id)
    return nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    users := make([]domain.User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, *user)
    }
    return users, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    for _, user := range r.users {
        if user.Email == email {
            return user, nil
        }
    }
    return nil, domain.ErrUserNotFound
}

