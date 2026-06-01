package domain

import "errors"

// ===========================================
// Domain Errors
// ===========================================
// These errors represent business rule violations
// and are used across all layers

var (
	// User errors
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailRequired    = errors.New("email is required")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmailExists      = errors.New("email already exists")
	ErrNameRequired     = errors.New("name is required")
	ErrNameTooShort     = errors.New("name must be at least 2 characters")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrInvalidPassword  = errors.New("invalid password")

	// Authentication errors
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("access forbidden")
	ErrTokenExpired    = errors.New("token has expired")
	ErrTokenInvalid    = errors.New("invalid token")
	ErrRefreshRequired = errors.New("token refresh required")

	// General errors
	ErrNotFound   = errors.New("resource not found")
	ErrConflict   = errors.New("resource already exists")
	ErrValidation = errors.New("validation error")
	ErrInternal   = errors.New("internal server error")
	ErrBadRequest = errors.New("bad request")

	// Data errors
	ErrInvalidID       = errors.New("invalid ID format")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidQuantity = errors.New("invalid quantity")
)

// AppError is a custom error type with additional context
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error codes
const (
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeBadRequest   = "BAD_REQUEST"
)
