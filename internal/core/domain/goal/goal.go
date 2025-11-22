package goal

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/shopspring/decimal"
)

type Goal struct {
	domain.Base
	UserID uint

	Name          string          `gorm:"type:varchar(100)"`
	TargetAmount  decimal.Decimal `gorm:"type:numeric(10, 2);not null"`
	CurrentAmount decimal.Decimal `gorm:"type:numeric(10, 2);not null"`
	DueDate       time.Time
}
