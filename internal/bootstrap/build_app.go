package bootstrap

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/api"
	"github.com/gildo-cordeiro/mapleplan-api/internal/api/handlers"
	"github.com/gildo-cordeiro/mapleplan-api/internal/business"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/database"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/repositories"

	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
)

func Build() (*api.BuildRegistry, error) {
	cfg := LoadConfig()
	db, err := database.NewGormDB()
	if err != nil {
		return nil, err
	}

	storageAdapter, err := BuildS3Storage(cfg)
	if err != nil {
		utils.Log.Warnf("buildStorage returned error (non-fatal in this context): %v", err)
		return nil, err
	}

	userRepo := repositories.NewGormUserRepository(db)
	goalRepo := repositories.NewGormGoalRepository(db)
	coupleRepo := repositories.NewGormCoupleRepository(db)
	txtManager := repositories.NewGormTransactionManager(db)

	userService := business.NewUserService(userRepo, coupleRepo, txtManager)
	goalService := business.NewGoalService(userRepo, goalRepo, coupleRepo, txtManager)
	business.NewStorageService(storageAdapter)

	health := handlers.HealthCheck{}
	userHandler := handlers.UserHandler{UserService: userService}
	goalHandler := handlers.GoalHandler{GoalService: goalService}
	authHandler := handlers.AuthHandler{UserService: userService}

	return &api.BuildRegistry{
		HealthHandler: health,
		UserHandler:   userHandler,
		GoalHandler:   goalHandler,
		AuthHandler:   authHandler,
	}, nil
}
