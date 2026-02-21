package repositories

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/goal"
)

type GoalRepository interface {
	FindByID(ctx context.Context, id string) (*goal.Goal, error)
	FindGoals(ctx context.Context, userID string) ([]*goal.Goal, error)
	FindWidgetGoals(ctx context.Context, userID string, limit int) ([]*goal.Goal, error)
	CountGoalsByStatus(ctx context.Context, userID string) (map[goal.Status]int, error)
	Save(ctx context.Context, goal *goal.Goal) error
	Update(ctx context.Context, id string, goal *goal.Goal) error
	UpdateStatus(ctx context.Context, id string, status goal.Status) error
	Delete(ctx context.Context, id string) error
}
