package services

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/mapper"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/response"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
)

type GoalServiceImpl struct {
	UserRepository repositories.UserRepository
	GoalRepository repositories.GoalRepository
	txManager      ports.TransactionManager
}

func NewGoalService(userRepo repositories.UserRepository, goalRepo repositories.GoalRepository, txManager ports.TransactionManager) services.GoalService {
	return &GoalServiceImpl{
		UserRepository: userRepo,
		GoalRepository: goalRepo,
		txManager:      txManager,
	}
}

func (g *GoalServiceImpl) GetWidgetGoals(ctx context.Context, userID string, limit int) ([]response.WidgetGoalResponse, error) {
	var goals []response.WidgetGoalResponse

	err := g.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		widgetGoals, err := g.GoalRepository.FindWidgetGoals(txCtx, userID, limit)
		if err != nil {
			return err
		}
		for _, goal := range widgetGoals {
			goals = append(goals, mapper.ToWidgetGoalResponse(goal))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (g *GoalServiceImpl) CreateGoal(userID string, req request.CreateGoalRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoalServiceImpl) GetGoals(userID string) ([]response.GoalResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoalServiceImpl) UpdateGoal(userID string, goalID string, req request.UpdateGoalRequest) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoalServiceImpl) DeleteGoal(userID string, goalID string) error {
	//TODO implement me
	panic("implement me")
}
