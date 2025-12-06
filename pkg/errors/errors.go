package errors

import "fmt"

// AppError represents a custom error type for the application.
type AppError struct {
    Code    int
    Message string
}

// NewAppError creates a new AppError with the given code and message.
func NewAppError(code int, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
    }
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// IsNotFound checks if the error is a not found error.
func IsNotFound(err error) bool {
    if appErr, ok := err.(*AppError); ok {
        return appErr.Code == 404
    }
    return false
}

// IsInternal checks if the error is an internal server error.
func IsInternal(err error) bool {
    if appErr, ok := err.(*AppError); ok {
        return appErr.Code == 500
    }
    return false
}