package services

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/response"
)

type GoalService interface {
	GetWidgetGoals(ctx context.Context, userID string, limit int) ([]response.WidgetGoalResponse, error)
	CreateGoal(userID string, req request.CreateGoalRequest) (string, error)
	GetGoals(userID string) ([]response.GoalResponse, error)
	UpdateGoal(userID string, goalID string, req request.UpdateGoalRequest) error
	DeleteGoal(userID string, goalID string) error
}
