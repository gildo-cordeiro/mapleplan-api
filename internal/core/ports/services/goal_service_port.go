package services

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/response"
)

type GoalService interface {
	GetWidgetGoals(ctx context.Context, userID string, limit int) ([]response.WidgetGoalResponse, error)
	CreateGoal(ctx context.Context, req request.CreateGoalRequest) error
	GetGoals(ctx context.Context, userID string) ([]response.GoalResponse, error)
	GetStatusCounts(ctx context.Context, userID string) (response.GoalStatusCountResponse, error)
	UpdateGoal(userID string, goalID string, req request.UpdateGoalRequest) error
	UpdateStatus(ctx context.Context, goalID string, status string) error
	DeleteGoal(userID string, goalID string) error
}
