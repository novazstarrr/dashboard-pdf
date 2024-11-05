// internal/domain/requests.go
package domain

import "time"

type RegisterRequest struct {
    Email     string    `json:"email" validate:"required,email"`
    Password  string    `json:"password" validate:"required,min=6"`
    FirstName string    `json:"firstName" validate:"required"`
    Surname   string    `json:"surname" validate:"required"`
    DOB       time.Time `json:"dob" validate:"required"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}



