// internal/repository/sqlite/file.go
package sqlite

import (
    "context"
    "gorm.io/gorm"
    "tech-test/backend/internal/domain"
    "tech-test/backend/internal/repository/interfaces"
    "log"
)

type fileRepository struct {
    db *gorm.DB
}

func NewFileRepository(db *gorm.DB) interfaces.FileRepository {
    return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *domain.File) error {
    result := r.db.Create(file)
    if result.Error != nil {
        return domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to create file",
            result.Error,
        )
    }
    return nil
}

func (r *fileRepository) GetByID(ctx context.Context, id uint) (*domain.File, error) {
    var file domain.File
    if err := r.db.First(&file, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrFileNotFound
        }
        return nil, domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to get file",
            err,
        )
    }
    return &file, nil
}

func (r *fileRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.File, error) {
    var files []domain.File
    if err := r.db.Where("user_id = ?", userID).Find(&files).Error; err != nil {
        return nil, err
    }
    return files, nil
}

func (r *fileRepository) GetUserFilesPaginated(ctx context.Context, userID uint, page, pageSize int) ([]domain.File, int64, error) {
    var files []domain.File
    var total int64

    if err := r.db.Model(&domain.File{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * pageSize

    if err := r.db.Where("user_id = ?", userID).
        Offset(offset).
        Limit(pageSize).
        Find(&files).Error; err != nil {
        return nil, 0, err
    }

    return files, total, nil
}

func (r *fileRepository) List(ctx context.Context) ([]domain.File, error) {
    var files []domain.File
    if err := r.db.Find(&files).Error; err != nil {
        return nil, err
    }
    return files, nil
}

func (r *fileRepository) Delete(ctx context.Context, id uint) error {
    return r.db.Delete(&domain.File{}, id).Error
}

func (r *fileRepository) SearchFiles(ctx context.Context, userID uint, searchTerm string) ([]domain.File, error) {
    var files []domain.File
    
    log.Printf("Searching for term: %s", searchTerm)
    
    query := r.db.Where("user_id = ?", userID)
    
    if searchTerm != "" {
        searchPattern := "%" + searchTerm + "%"
        query = query.Where(
            "LOWER(name) LIKE LOWER(?) OR LOWER(mime_type) LIKE LOWER(?)",
            searchPattern,
            searchPattern,
        )
    }
    
    if err := query.Find(&files).Error; err != nil {
        log.Printf("Search error: %v", err)
        return nil, err
    }
    
    log.Printf("Found %d files", len(files))
    return files, nil
}

func (r *fileRepository) UpdateShareableID(ctx context.Context, fileID uint, shareableID string) error {
    result := r.db.Model(&domain.File{}).
        Where("id = ?", fileID).
        Update("shareable_id", shareableID)
    
    if result.Error != nil {
        return domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to update shareable ID",
            result.Error,
        )
    }
    
    if result.RowsAffected == 0 {
        return domain.ErrFileNotFound
    }
    
    return nil
}

func (r *fileRepository) GetFileByShareID(ctx context.Context, shareID string) (*domain.File, error) {
    var file domain.File
    result := r.db.Where("shareable_id = ?", shareID).First(&file)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, domain.ErrFileNotFound
        }
        return nil, domain.NewAPIError(
            500,
            domain.ErrCodeInternal,
            "Failed to get file by share ID",
            result.Error,
        )
    }
    return &file, nil
}
