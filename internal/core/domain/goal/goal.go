package goal

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/couple"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	"github.com/shopspring/decimal"
)

type Goal struct {
	domain.Base

	UserID   *string `gorm:"type:uuid;index" json:"userId,omitempty"`
	CoupleID *string `gorm:"type:uuid;index" json:"coupleId,omitempty"`

	Name          string          `gorm:"type:varchar(100);not null" json:"name"`
	TargetAmount  decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"targetAmount"`
	CurrentAmount decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"currentAmount"`
	DueDate       time.Time       `json:"dueDate"`

	User   *user.User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Couple *couple.Couple `gorm:"foreignKey:CoupleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"couple,omitempty"`
}
