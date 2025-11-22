package transaction

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	domain.Base
	UserID uint

	Amount          decimal.Decimal `gorm:"type:numeric(10, 2);not null"`
	Currency        string          `gorm:"type:varchar(3);not null"`
	Type            string          `gorm:"type:numeric(10, 2);not null"`
	Description     string
	TransactionDate time.Time
}
