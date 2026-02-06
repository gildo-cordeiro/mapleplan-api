package repository

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/goal"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type GoalRepositoryImpl struct {
	db *gorm.DB
}

func (g *GoalRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return g.db
	}
	if tx, ok := ctx.Value("tx_key").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return g.db.WithContext(ctx)
}

func NewGormGoalRepository(db *gorm.DB) repositories.GoalRepository {
	return &GoalRepositoryImpl{db: db}
}

func (g *GoalRepositoryImpl) FindWidgetGoals(ctx context.Context, userID string, limit int) ([]*goal.Goal, error) {
	var goals []*goal.Goal
	err := g.getDB(ctx).
		Joins("JOIN couples ON goals.couple_id = couples.id").
		Where("couples.user_a_id = ? OR couples.user_b_id = ?", userID, userID).
		Order("goals.created_at DESC").
		Limit(limit).
		Find(&goals).Error
	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (g *GoalRepositoryImpl) FindByID(ctx context.Context, id string) (*goal.Goal, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoalRepositoryImpl) Save(ctx context.Context, goal *goal.Goal) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoalRepositoryImpl) Update(ctx context.Context, id string, goal *goal.Goal) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoalRepositoryImpl) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
