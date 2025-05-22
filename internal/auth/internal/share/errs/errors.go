package errs

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal error")
	ErrAlreadyExists = errors.New("already exists")
)

// AppError represents an application error
type AppError struct {
	Op      string // Operation that failed
	Kind    error  // Category of error
	Err     error  // The actual error
	Message string // User-friendly message
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(op string, kind error, err error, message string) *AppError {
	return &AppError{
		Op:      op,
		Kind:    kind,
		Err:     err,
		Message: message,
	}
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return errors.Is(appErr.Kind, ErrNotFound)
	}
	return false
}

// IsInvalidInput checks if the error is an invalid input error
func IsInvalidInput(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return errors.Is(appErr.Kind, ErrInvalidInput)
	}
	return false
}

// IsUnauthorized checks if the error is an unauthorized error
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return errors.Is(appErr.Kind, ErrUnauthorized)
	}
	return false
}

// IsForbidden checks if the error is a forbidden error
func IsForbidden(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return errors.Is(appErr.Kind, ErrForbidden)
	}
	return false
}

// IsInternal checks if the error is an internal error
func IsInternal(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return errors.Is(appErr.Kind, ErrInternal)
	}
	return false
}

// IsAlreadyExists checks if the error is an already exists error
func IsAlreadyExists(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return errors.Is(appErr.Kind, ErrAlreadyExists)
	}
	return false
}
