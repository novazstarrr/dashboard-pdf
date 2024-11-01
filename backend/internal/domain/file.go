// internal/domain/file.go
package domain

import (
	"fmt"
	"strings"
	"time"
)

type File struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	UserID      uint      `json:"userId" gorm:"not null;index"`
	Name        string    `json:"name" gorm:"not null;index"`
	Path        string    `json:"-" gorm:"not null"`  
	MimeType    string    `json:"mimeType" gorm:"not null"`
	Size        int64     `json:"size" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ShareableID string    `json:"shareableId" gorm:"index"`
	ContentType string    `json:"contentType" gorm:"not null"`
}


type FileDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	MimeType    string    `json:"mimeType"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"createdAt"`
	DownloadURL string    `json:"downloadUrl"`      
	ShareURL    string    `json:"shareUrl,omitempty"` 
}


const (
	MaxFileSize    = 100 * 1024 * 1024 // 100MB
	MaxFileNameLen = 255
)

var AllowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"application/pdf": true,
	"text/plain":      true,
}

func NewFile(userID uint, name string, size int64, contentType string) (*File, error) {
	if err := validateFileName(name); err != nil {
		return nil, err
	}

	if err := validateFileSize(size); err != nil {
		return nil, err
	}

	if err := validateContentType(contentType); err != nil {
		return nil, err
	}

	cleanName := cleanFileName(name)

	return &File{
		UserID:      userID,
		Name:        cleanName,
		MimeType:    contentType,
		Size:        size,
		ContentType: contentType,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (f *File) ToDTO(baseURL string) FileDTO {
	return FileDTO{
		ID:          f.ID,
		Name:        f.Name,
		MimeType:    f.MimeType,
		Size:        f.Size,
		CreatedAt:   f.CreatedAt,
		DownloadURL: f.generateDownloadURL(baseURL),
		ShareURL:    f.generateShareURL(baseURL),
	}
}


func validateFileName(name string) error {
	if name = strings.TrimSpace(name); name == "" {
		return NewAPIError(
			400,
			ErrCodeInvalidInput,
			"File name cannot be empty",
			nil,
		)
	}

	if len(name) > MaxFileNameLen {
		return NewAPIError(
			400,
			ErrCodeInvalidInput,
			fmt.Sprintf("File name cannot exceed %d characters", MaxFileNameLen),
			nil,
		)
	}

	return nil
}

func validateFileSize(size int64) error {
	if size <= 0 {
		return NewAPIError(
			400,
			ErrCodeInvalidInput,
			"File size must be greater than 0",
			nil,
		)
	}

	if size > MaxFileSize {
		return NewAPIError(
			400,
			ErrCodeInvalidInput,
			fmt.Sprintf("File size cannot exceed %d bytes", MaxFileSize),
			nil,
		)
	}

	return nil
}

func validateContentType(contentType string) error {
	if contentType == "" {
		return NewAPIError(
			400,
			ErrCodeInvalidInput,
			"Content type cannot be empty",
			nil,
		)
	}

	if !AllowedMimeTypes[contentType] {
		return NewAPIError(
			400,
			ErrCodeInvalidInput,
			"Unsupported file type",
			nil,
		)
	}

	return nil
}


func (f *File) generateDownloadURL(baseURL string) string {
	return fmt.Sprintf("%s/api/files/%d/download", baseURL, f.ID)
}

func (f *File) generateShareURL(baseURL string) string {
	if f.ShareableID == "" {
		return ""
	}
	return fmt.Sprintf("%s/api/files/shared/%s", baseURL, f.ShareableID)
}

func cleanFileName(name string) string {
	name = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			return '_'
		}
		return r
	}, name)
	
	return name
}

