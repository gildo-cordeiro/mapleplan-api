package repositories

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/profile"
)

type ProfileRepository interface {
	FindByID(ctx context.Context, id string) (*profile.ImmigrationProfile, error)
	FindByUserID(ctx context.Context, userID string) ([]*profile.ImmigrationProfile, error)
	Save(ctx context.Context, p *profile.ImmigrationProfile) error
	Update(ctx context.Context, id string, p *profile.ImmigrationProfile) error
	Delete(ctx context.Context, id string) error
}

type ProfileMemberRepository interface {
	FindByID(ctx context.Context, id string) (*profile.ProfileMember, error)
	FindByProfileID(ctx context.Context, profileID string) ([]*profile.ProfileMember, error)
	FindByUserID(ctx context.Context, userID string) ([]*profile.ProfileMember, error)
	Save(ctx context.Context, m *profile.ProfileMember) error
	Update(ctx context.Context, id string, m *profile.ProfileMember) error
	Delete(ctx context.Context, id string) error
}
