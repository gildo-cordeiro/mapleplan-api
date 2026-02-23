package finance

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/profile"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
)

type Finance struct {
	models.Base

	UserID               *string     `gorm:"type:uuid;index" json:"userId,omitempty"`
	ImmigrationProfileID *string     `gorm:"type:uuid;index" json:"immigrationProfileId,omitempty"` // NULL = private, filled = shared
	Category             Category    `gorm:"type:varchar(50);not null" json:"category"`
	Type                 FinanceType `gorm:"type:varchar(20);not null;index" json:"type"`
	Amount               float64     `gorm:"type:decimal(15,2);not null" json:"amount"`
	Currency             string      `gorm:"type:varchar(3);not null;default:CAD" json:"currency"`
	Reference            *string     `gorm:"type:varchar(255)" json:"reference,omitempty"`
	Note                 *string     `gorm:"type:text" json:"note,omitempty"`
	Date                 time.Time   `gorm:"type:date;not null;index" json:"date"`

	User               *user.User                  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	ImmigrationProfile *profile.ImmigrationProfile `gorm:"foreignKey:ImmigrationProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"immigrationProfile,omitempty"`
}

func (f *Finance) IsPrivate() bool {
	return f.ImmigrationProfileID == nil
}

func (f *Finance) ShareWithProfile(profileID string) {
	f.ImmigrationProfileID = &profileID
}

func (f *Finance) MakePrivate() {
	f.ImmigrationProfileID = nil
}
