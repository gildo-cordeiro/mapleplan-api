package goal

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/couple"
	"github.com/shopspring/decimal"
)

type Goal struct {
	domain.Base

	CoupleID *string `gorm:"type:uuid;index" json:"coupleId,omitempty"`

	Name          string          `gorm:"type:varchar(100);not null" json:"name"`
	Status        string          `gorm:"type:varchar(20);not null" json:"status"`
	Progress      int             `gorm:"type:integer;not null" json:"progress"`
	TargetAmount  decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"targetAmount"`
	CurrentAmount decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"currentAmount"`
	DueDate       time.Time       `json:"type:date" json:"dueDate"`
	Phase         string          `gorm:"type:varchar(20);not null" json:"phase"`
	Priority      string          `gorm:"type:varchar(20);not null" json:"priority"`
	Description   *string         `gorm:"type:text" json:"description,omitempty"`

	Couple *couple.Couple `gorm:"foreignKey:CoupleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"couple,omitempty"`
}
