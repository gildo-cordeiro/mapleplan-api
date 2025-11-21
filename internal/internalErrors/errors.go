package internalErrors

import "errors"

var (
	ErrInvalidEmail   = errors.New("invalid email")
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternal       = errors.New("internal server error")
)
