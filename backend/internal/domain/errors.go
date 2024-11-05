// internal/domain/errors.go
package domain

import (
	"fmt"
	"net/http"
)

const (
	ErrCodeInvalidInput     = 4000
	ErrCodeAuthentication   = 4001
	ErrCodeAuthorization    = 4003
	
	ErrCodeNotFound        = 4004
	ErrCodeConflict        = 4009
	ErrCodeInternal        = 5000
	ErrCodeUserNotFound    = 4040
	ErrCodeDuplicateEmail  = 4010
	ErrCodeInvalidFileType = 4011
	ErrCodeFileTooLarge    = 4012
	ErrCodeFileNotFound    = 4041
)

type APIError struct {
	StatusCode int         `json:"status"`     
	Code       int         `json:"code"`       
	Message    string      `json:"message"`    
	Details    interface{} `json:"details,omitempty"` 
	Internal   error       `json:"-"`          
}

func (e *APIError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}


func NewAPIError(statusCode, code int, message string, internal error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:      code,
		Message:   message,
		Internal:  internal,
	}
}


var (
	ErrInvalidInput = NewAPIError(
		http.StatusBadRequest,
		ErrCodeInvalidInput,
		"Invalid input provided",
		nil,
	)

	ErrUnauthorized = NewAPIError(
		http.StatusUnauthorized,
		ErrCodeAuthentication,
		"Authentication required",
		nil,
	)

	ErrForbidden = NewAPIError(
		http.StatusForbidden,
		ErrCodeAuthorization,
		"Access denied",
		nil,
	)

	ErrNotFound = NewAPIError(
		http.StatusNotFound,
		ErrCodeNotFound,
		"Resource not found",
		nil,
	)

	ErrConflict = NewAPIError(
		http.StatusConflict,
		ErrCodeConflict,
		"Resource conflict",
		nil,
	)

	ErrInternal = NewAPIError(
		http.StatusInternalServerError,
		ErrCodeInternal,
		"Internal server error",
		nil,
	)

	ErrUserNotFound = NewAPIError(
		http.StatusNotFound,
		ErrCodeUserNotFound,
		"User not found",
		nil,
	)

	ErrDuplicateEmail = NewAPIError(
		http.StatusConflict,
		ErrCodeDuplicateEmail,
		"Email already exists",
		nil,
	)

	ErrInvalidFileType = NewAPIError(
		http.StatusBadRequest,
		ErrCodeInvalidFileType,
		"Invalid file type",
		nil,
	)

	ErrFileTooLarge = NewAPIError(
		http.StatusBadRequest,
		ErrCodeFileTooLarge,
		"File exceeds maximum size limit",
		nil,
	)

	ErrFileNotFound = NewAPIError(
		http.StatusNotFound,
		ErrCodeFileNotFound,
		"File not found",
		nil,
	)
)


func NewInvalidInputError(details interface{}) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Code:      ErrCodeInvalidInput,
		Message:   "Invalid input provided",
		Details:   details,
	}
}

func NewNotFoundError(resource string) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Code:      ErrCodeNotFound,
		Message:   fmt.Sprintf("%s not found", resource),
	}
}

func WrapError(err error) *APIError {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:      ErrCodeInternal,
		Message:   "Internal server error",
		Internal:  err,
	}
}


func NewDuplicateEmailError(email string) *APIError {
	return &APIError{
		StatusCode: http.StatusConflict,
		Code:      ErrCodeDuplicateEmail,
		Message:   fmt.Sprintf("Email %s already exists", email),
	}
}

func NewFileTooLargeError(size int64, maxSize int64) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Code:      ErrCodeFileTooLarge,
		Message:   fmt.Sprintf("File size %d exceeds maximum size %d", size, maxSize),
		Details: map[string]int64{
			"size":     size,
			"maxSize": maxSize,
		},
	}
}

