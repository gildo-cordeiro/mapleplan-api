package repositories

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/finance"
)

type FinanceRepository interface {
	FindByID(ctx context.Context, id string) (*finance.Finance, error)
	FindByUserID(ctx context.Context, userID string) ([]*finance.Finance, error)
	FindByProfileID(ctx context.Context, profileID string) ([]*finance.Finance, error)
	FindByUserAndProfile(ctx context.Context, userID, profileID string) ([]*finance.Finance, error)
	Save(ctx context.Context, f *finance.Finance) error
	Update(ctx context.Context, id string, f *finance.Finance) error
	Delete(ctx context.Context, id string) error
}
