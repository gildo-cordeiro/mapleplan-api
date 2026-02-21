package repositories

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/ports"
	"gorm.io/gorm"
)

type gormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) ports.TransactionManager {
	return &gormTransactionManager{db: db}
}

func (g *gormTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Pass the transaction DB instance in the context
		txCtx := context.WithValue(ctx, "tx_key", tx)
		return fn(txCtx)
	})
}
