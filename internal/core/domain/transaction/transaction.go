package transaction

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/couple"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	domain.Base

	UserID   *string `gorm:"type:uuid;index" json:"userId,omitempty"`
	CoupleID *string `gorm:"type:uuid;index" json:"coupleId,omitempty"`

	Amount      decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"amount"`
	Currency    *string         `gorm:"type:varchar(10)" json:"currency,omitempty"`
	Type        string          `gorm:"type:varchar(50);not null" json:"type"` // ex: "income" | "expense"
	Description *string         `gorm:"type:text" json:"description,omitempty"`
	Date        time.Time       `gorm:"not null" json:"date"`

	User   *user.User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Couple *couple.Couple `gorm:"foreignKey:CoupleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"couple,omitempty"`
	// CategoryID  *string         `gorm:"type:uuid;index" json:"categoryId,omitempty"`
}
