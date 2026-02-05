package repository

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/couple"
	couplePort "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type CoupleRepositoryImpl struct {
	db *gorm.DB
}

func (r *CoupleRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return r.db
	}
	if tx, ok := ctx.Value("tx_key").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func NewGormCoupleRepository(db *gorm.DB) couplePort.CoupleRepository {
	return &CoupleRepositoryImpl{db: db}
}

func (r *CoupleRepositoryImpl) Save(ctx context.Context, c *couple.Couple) error {
	if err := r.getDB(ctx).Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (r *CoupleRepositoryImpl) FindByID(id string) (*couple.Couple, error) {
	var c couple.Couple
	err := r.db.Where("id = ?", id).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *CoupleRepositoryImpl) FindByUserID(ctx context.Context, userID string) (*couple.Couple, error) {
	var c couple.Couple
	err := r.getDB(ctx).Where("user_a_id = ? OR user_b_id = ?", userID, userID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
