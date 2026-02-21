package repositories

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Save(u *user.User) (string, error)
	FindByID(ctx context.Context, id string) (*user.User, error)
	Update(ctx context.Context, id string, u *user.User) error
	SearchByName(userID string, name string) ([]*user.User, error)
}
