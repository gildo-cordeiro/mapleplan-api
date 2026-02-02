package repositories

import "github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"

type UserRepository interface {
	// FindByEmail returns a user by email. Password verification is handled
	// by the service (use-case) layer.
	FindByEmail(email string) (*user.User, error)
	Save(u *user.User) (string, error)
	FindByID(id string) (*user.User, error)
	Update(id string, u *user.User) error
}
