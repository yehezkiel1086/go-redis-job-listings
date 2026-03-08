package domain

import "errors"

var (
	ErrUnauthorized        error = errors.New("unauthorized")
	ErrForbidden           error = errors.New("forbidden: insufficient permissions")
	ErrInvalidToken        error = errors.New("invalid or expired token")
	ErrInvalidCredentials  error = errors.New("invalid credentials")
	ErrDuplicateEmail      error = errors.New("email already registered")
	ErrNotFound            error = errors.New("not found")
	ErrInternalServerError error = errors.New("internal server error")
)
