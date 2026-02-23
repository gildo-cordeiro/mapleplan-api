package repositories

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/goal"
	ports "github.com/gildo-cordeiro/mapleplan-api/internal/ports/repositories"
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

func NewGormGoalRepository(db *gorm.DB) ports.GoalRepository {
	return &GoalRepositoryImpl{db: db}
}

func (g *GoalRepositoryImpl) FindWidgetGoals(ctx context.Context, userID string, profileID *string, limit int) ([]*goal.Goal, error) {
	var goals []*goal.Goal
	db := g.getDB(ctx)

	if profileID != nil {
		err := db.Preload("ImmigrationProfile").Preload("User").
			Where("profile_id = ?", *profileID).
			Order("created_at DESC").
			Limit(limit).
			Find(&goals).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := db.Preload("ImmigrationProfile").Preload("User").
			Where("user_id = ?", userID).
			Order("created_at DESC").
			Limit(limit).
			Find(&goals).Error
		if err != nil {
			return nil, err
		}
	}

	return goals, nil
}

func (g *GoalRepositoryImpl) CountGoalsByStatus(ctx context.Context, userID string) (map[goal.Status]int, error) {
	type Result struct {
		Status goal.Status
		Count  int
	}
	var results []Result
	db := g.getDB(ctx)

	// Count goals where user is either the owner or a member of the profile
	err := db.
		Table("goals").
		Select("status, COUNT(*) as count").
		Joins("LEFT JOIN profile_members ON goals.immigration_profile_id = profile_members.immigration_profile_id").
		Where("profile_members.user_id = ? OR goals.user_id = ?", userID, userID).
		Group("status").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	resultList := make(map[goal.Status]int, len(results))
	for _, r := range results {
		resultList[r.Status] = r.Count
	}
	return resultList, nil
}

func (g *GoalRepositoryImpl) FindByID(ctx context.Context, id string) (*goal.Goal, error) {
	var foundedGoal goal.Goal
	err := g.getDB(ctx).Preload("User").Preload("ImmigrationProfile").First(&foundedGoal, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &foundedGoal, nil
}

func (g *GoalRepositoryImpl) Save(ctx context.Context, goal *goal.Goal) error {
	db := g.getDB(ctx)
	if err := db.Create(goal).Error; err != nil {
		return err
	}
	return nil
}

func (g *GoalRepositoryImpl) FindGoals(ctx context.Context, userID string, limit *int) ([]*goal.Goal, error) {
	var goals []*goal.Goal
	db := g.getDB(ctx)
	query := db.Preload("User").Preload("ImmigrationProfile").
		Joins("LEFT JOIN profile_members ON goals.immigration_profile_id = profile_members.immigration_profile_id").
		Where("profile_members.user_id = ? OR goals.user_id = ?", userID, userID).
		Order("goals.created_at DESC")

	if limit != nil {
		query = query.Limit(*limit)
	}

	err := query.Find(&goals).Error
	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (g *GoalRepositoryImpl) UpdateStatus(ctx context.Context, id string, status goal.Status) error {
	return g.getDB(ctx).Model(&goal.Goal{}).Where("id = ?", id).Update("status", status).Error
}

func (g *GoalRepositoryImpl) Update(ctx context.Context, id string, goal *goal.Goal) error {
	foundedGoal, err := g.FindByID(ctx, id)
	if err != nil {
		return err
	}
	err = foundedGoal.UpdateFields(goal)
	if err != nil {
		return err
	}

	// Ensure loaded associations are nil so GORM doesn't try to create/update related records
	foundedGoal.User = nil
	foundedGoal.ImmigrationProfile = nil

	return g.getDB(ctx).Save(foundedGoal).Error
}

func (g *GoalRepositoryImpl) Delete(ctx context.Context, id string) error {
	return g.getDB(ctx).Delete(&goal.Goal{}, "id = ?", id).Error
}
