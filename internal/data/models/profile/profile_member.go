package profile

import (
	"errors"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
)

type ProfileMember struct {
	models.Base

	ImmigrationProfileID string     `gorm:"type:uuid;not null;uniqueIndex:idx_profile_user" json:"ImmigrationProfileID"`
	UserID               string     `gorm:"type:uuid;not null;uniqueIndex:idx_profile_user" json:"userId"`
	Role                 MemberRole `gorm:"type:varchar(50);not null" json:"role"`
	JoinedAt             *time.Time `gorm:"type:date" json:"joinedAt,omitempty"`

	ImmigrationProfile *ImmigrationProfile `gorm:"foreignKey:ImmigrationProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"immigrationProfile,omitempty"`
	User               *user.User          `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func NewProfileMember(profileID, userID string, role MemberRole) (*ProfileMember, error) {
	if profileID == "" || userID == "" {
		return nil, errors.New("profileID and userID are required")
	}

	if !isValidRole(role) {
		return nil, errors.New("invalid role: must be Primary or Spouse")
	}

	now := time.Now()
	return &ProfileMember{
		ImmigrationProfileID: profileID,
		UserID:               userID,
		Role:                 role,
		JoinedAt:             &now,
	}, nil
}
