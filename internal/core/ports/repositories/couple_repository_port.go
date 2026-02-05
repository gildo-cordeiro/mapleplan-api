package repositories

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/couple"
)

type CoupleRepository interface {
	Save(ctx context.Context, c *couple.Couple) error
	FindByID(id string) (*couple.Couple, error)
	FindByUserID(ctx context.Context, userID string) (*couple.Couple, error)
}
