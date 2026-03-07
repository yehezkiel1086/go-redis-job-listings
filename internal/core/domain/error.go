package domain

import "errors"

var (
	ErrInvalidCredentials  error = errors.New("invalid credentials")
	ErrDuplicateEmail      error = errors.New("email already registered")
	ErrNotFound            error = errors.New("not found")
	ErrInternalServerError error = errors.New("internal server error")
)
