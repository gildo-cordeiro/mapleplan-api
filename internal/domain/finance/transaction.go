package finance

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID uint

	Amount          decimal.Decimal `gorm:"type:numeric(10, 2);not null"`
	Currency        string          `gorm:"type:varchar(3);not null"`
	Type            string          `gorm:"type:numeric(10, 2);not null"`
	Description     string
	TransactionDate time.Time
}
