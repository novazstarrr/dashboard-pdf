package domain

import "time"


type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null" example:"user@example.com"`
	Password  string    `json:"password" gorm:"not null" example:"password123"` 
	FirstName string    `json:"firstName" gorm:"not null" example:"John"`
	Surname   string    `json:"surname" gorm:"not null" example:"Doe"`
	DOB       time.Time `json:"dob" gorm:"not null" example:"1990-01-01T00:00:00Z"`
	CreatedAt time.Time `json:"createdAt,omitempty" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" example:"2024-01-01T00:00:00Z"`
}

type UserContract struct {
	UserID     string    `json:"userId"`
	FileName   string    `json:"fileName"`
	UploadedAt time.Time `json:"uploadedAt"`
	Path       string    `json:"-"`  
}
