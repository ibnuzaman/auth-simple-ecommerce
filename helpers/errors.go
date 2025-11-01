package helpers

import "net/http"

// AppError represents standardized application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error constructors
func ErrBadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message, "")
}

func ErrUnauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, "")
}

func ErrForbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, "")
}

func ErrNotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, "")
}

func ErrConflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message, "")
}

func ErrInternalServer(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, "")
}

func ErrValidation(details string) *AppError {
	return NewAppError(http.StatusBadRequest, "Validation error", details)
}
