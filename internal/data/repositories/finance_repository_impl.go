package repositories

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/finance"
	"github.com/gildo-cordeiro/mapleplan-api/internal/ports/repositories"
	"gorm.io/gorm"
)

type FinanceRepositoryImpl struct {
	db *gorm.DB
}

func (f *FinanceRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return f.db
	}
	if tx, ok := ctx.Value("tx_key").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return f.db.WithContext(ctx)
}

func NewGormFinanceRepository(db *gorm.DB) repositories.FinanceRepository {
	return &FinanceRepositoryImpl{db: db}
}

func (f *FinanceRepositoryImpl) FindByID(ctx context.Context, id string) (*finance.Finance, error) {
	var fin finance.Finance
	err := f.getDB(ctx).Where("id = ?", id).First(&fin).Error
	if err != nil {
		return nil, err
	}
	return &fin, nil
}

func (f *FinanceRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*finance.Finance, error) {
	var finances []*finance.Finance
	err := f.getDB(ctx).
		Where("user_id = ? AND profile_id IS NULL", userID).
		Order("date DESC").
		Find(&finances).Error
	if err != nil {
		return nil, err
	}
	return finances, nil
}

func (f *FinanceRepositoryImpl) FindByProfileID(ctx context.Context, profileID string) ([]*finance.Finance, error) {
	var finances []*finance.Finance
	err := f.getDB(ctx).
		Where("profile_id = ?", profileID).
		Order("date DESC").
		Find(&finances).Error
	if err != nil {
		return nil, err
	}
	return finances, nil
}

func (f *FinanceRepositoryImpl) FindByUserAndProfile(ctx context.Context, userID, profileID string) ([]*finance.Finance, error) {
	var finances []*finance.Finance
	err := f.getDB(ctx).
		Where("(user_id = ? AND profile_id IS NULL) OR profile_id = ?", userID, profileID).
		Order("date DESC").
		Find(&finances).Error
	if err != nil {
		return nil, err
	}
	return finances, nil
}

func (f *FinanceRepositoryImpl) Save(ctx context.Context, fin *finance.Finance) error {
	return f.getDB(ctx).Create(fin).Error
}

func (f *FinanceRepositoryImpl) Update(ctx context.Context, id string, fin *finance.Finance) error {
	foundFin, err := f.FindByID(ctx, id)
	if err != nil {
		return err
	}

	foundFin.Category = fin.Category
	foundFin.Type = fin.Type
	foundFin.Amount = fin.Amount
	foundFin.Currency = fin.Currency
	foundFin.Reference = fin.Reference
	foundFin.Note = fin.Note
	foundFin.Date = fin.Date
	foundFin.ImmigrationProfileID = fin.ImmigrationProfileID

	return f.getDB(ctx).Save(foundFin).Error
}

func (f *FinanceRepositoryImpl) Delete(ctx context.Context, id string) error {
	return f.getDB(ctx).Where("id = ?", id).Delete(&finance.Finance{}).Error
}
