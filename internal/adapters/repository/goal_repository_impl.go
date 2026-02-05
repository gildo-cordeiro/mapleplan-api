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

func NewGormGoalRepository(db *gorm.DB) repositories.GoalRepository {
	return &GoalRepositoryImpl{db: db}
}

func (g GoalRepositoryImpl) FindByID(ctx context.Context, id string) (*goal.Goal, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRepositoryImpl) Save(ctx context.Context, goal *goal.Goal) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRepositoryImpl) Update(ctx context.Context, id string, goal *goal.Goal) error {
	//TODO implement me
	panic("implement me")
}

func (g GoalRepositoryImpl) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
