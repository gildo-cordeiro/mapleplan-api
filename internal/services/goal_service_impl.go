package services

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/mapper"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/response"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/goal"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
)

type GoalServiceImpl struct {
	UserRepository repositories.UserRepository
	Couple         repositories.CoupleRepository
	GoalRepository repositories.GoalRepository
	txManager      ports.TransactionManager
}

func NewGoalService(userRepo repositories.UserRepository, goalRepo repositories.GoalRepository, couple repositories.CoupleRepository, txManager ports.TransactionManager) services.GoalService {
	return &GoalServiceImpl{
		UserRepository: userRepo,
		GoalRepository: goalRepo,
		Couple:         couple,
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
		for _, widgetGoal := range widgetGoals {
			goals = append(goals, mapper.ToWidgetGoalResponse(widgetGoal))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (g *GoalServiceImpl) GetStatusCounts(ctx context.Context, userID string) (response.GoalStatusCountResponse, error) {
	var countResponse response.GoalStatusCountResponse

	err := g.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		counts, err := g.GoalRepository.CountGoalsByStatus(txCtx, userID)
		if err != nil {
			return err
		}
		countResponse = response.GoalStatusCountResponse{
			NotStarted: counts[goal.NotStartedStatus],
			InProgress: counts[goal.InProgressStatus],
			Completed:  counts[goal.CompletedStatus],
			Total:      counts[goal.NotStartedStatus] + counts[goal.InProgressStatus] + counts[goal.CompletedStatus],
		}
		return nil
	})

	if err != nil {
		return response.GoalStatusCountResponse{}, err
	}
	return countResponse, nil
}

func (g *GoalServiceImpl) GetGoalByID(ctx context.Context, goalID string) (response.GoalResponse, error) {
	var goalResponse response.GoalResponse

	err := g.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		foundGoal, err := g.GoalRepository.FindByID(txCtx, goalID)
		if err != nil {
			return utils.ErrRecordNotFound
		}

		goalResponse = mapper.ToGoalUpdateResponse(foundGoal)
		return nil
	})

	if err != nil {
		return goalResponse, err
	}
	return goalResponse, nil
}

func (g *GoalServiceImpl) CreateGoal(ctx context.Context, req request.CreateGoalRequest) error {
	goalToSave := mapper.CreateGoalRequestToGoalDomain(&req)
	err := g.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if req.AssignedToUser != nil && *req.AssignedToUser != "" {
			user, err := g.UserRepository.FindByID(txCtx, *req.AssignedToUser)
			if err != nil {
				return err
			}
			goalToSave.UserId = &user.ID
		} else if req.AssignedToCouple != nil && *req.AssignedToCouple != "" {
			c, err := g.Couple.FindByID(txCtx, *req.AssignedToCouple)
			if err != nil {
				return err
			}
			goalToSave.CoupleID = &c.ID
		} else {
			return utils.ErrInvalidGoalAssignment
		}
		return g.GoalRepository.Save(txCtx, goalToSave)
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *GoalServiceImpl) GetGoals(ctx context.Context, userID string) ([]response.GoalResponse, error) {
	goals := make([]response.GoalResponse, 0)

	err := g.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		foundGoals, err := g.GoalRepository.FindGoals(txCtx, userID)
		if err != nil {
			return err
		}
		for _, widgetGoal := range foundGoals {
			goals = append(goals, mapper.ToGoalResponse(widgetGoal))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (g *GoalServiceImpl) UpdateStatus(ctx context.Context, goalID string, status string) error {
	return g.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		parsedStatus, ok := goal.StringToStatus(status)
		if !ok {
			return utils.ErrInvalidGoalStatus
		}
		return g.GoalRepository.UpdateStatus(txCtx, goalID, parsedStatus)
	})
}

func (g *GoalServiceImpl) UpdateGoal(ctx context.Context, goalID string, req request.UpdateGoalRequestBody) error {
	return g.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		goalToUpdate := mapper.UpdateGoalRequestToGoalDomain(&req)
		if goalToUpdate == nil {
			return utils.ErrInvalidParse
		}

		if req.AssignedToUser != "" {
			user, err := g.UserRepository.FindByID(txCtx, req.AssignedToUser)
			if err != nil {
				return err
			}
			goalToUpdate.UserId = &user.ID
			goalToUpdate.CoupleID = nil
		} else if req.AssignedToCouple != "" {
			c, err := g.Couple.FindByID(txCtx, req.AssignedToCouple)
			if err != nil {
				return err
			}
			goalToUpdate.CoupleID = &c.ID
			goalToUpdate.UserId = nil
		}

		return g.GoalRepository.Update(txCtx, goalID, goalToUpdate)
	})
}

func (g *GoalServiceImpl) DeleteGoal(userID string, goalID string) error {
	//TODO implement me
	panic("implement me")
}
