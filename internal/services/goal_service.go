package services

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/response"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
)

type GoalServiceImpl struct {
	GoalRepository repositories.GoalRepository
	txManager      ports.TransactionManager
}

func NewGoalService(goalRepo repositories.GoalRepository, txManager ports.TransactionManager) services.GoalService {
	return &GoalServiceImpl{
		GoalRepository: goalRepo,
		txManager:      txManager,
	}
}

func (g GoalServiceImpl) GetWidgetGoals(userID string, limit string) ([]response.WidgetGoalResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalServiceImpl) CreateGoal(userID string, req request.CreateGoalRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalServiceImpl) GetGoals(userID string) ([]response.GoalResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalServiceImpl) UpdateGoal(userID string, goalID string, req request.UpdateGoalRequest) error {
	//TODO implement me
	panic("implement me")
}

func (g GoalServiceImpl) DeleteGoal(userID string, goalID string) error {
	//TODO implement me
	panic("implement me")
}
