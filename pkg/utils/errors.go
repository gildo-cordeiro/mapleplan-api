package utils

import "errors"

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInternal              = errors.New("internal server error")
	ErrAlreadyExists         = errors.New("already exists")
	ErrNoFieldsToUpdate      = errors.New("no fields to update")
	ErrInvalidInput          = errors.New("invalid input")
	ErrRecordNotFound        = errors.New("record not found")
	ErrInvalidGoalAssignment = errors.New("invalid goal assignment")
)
