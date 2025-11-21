package domain

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Goal struct {
	gorm.Model
	UserID uint

	Name          string          `gorm:"type:varchar(100)"`
	TargetAmount  decimal.Decimal `gorm:"type:numeric(10, 2);not null"`
	CurrentAmount decimal.Decimal `gorm:"type:numeric(10, 2);not null"`
	DueDate       time.Time
}
