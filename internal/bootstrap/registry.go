package bootstrap

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/api"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/database"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/repository"
	"github.com/gildo-cordeiro/mapleplan-api/internal/services"
)

func Build() (*api.HandlerRegistry, error) {
	db, err := database.NewGormDB()
	if err != nil {
		return nil, err
	}

	coupRepo := repository.NewGormCoupleRepository(db)
	userRepo := repository.NewGormUserRepository(db)
	goalRepo := repository.NewGormGoalRepository(db)
	txtManager := repository.NewGormTransactionManager(db)

	userService := services.NewUserService(userRepo, coupRepo, txtManager)
	goalService := services.NewGoalService(userRepo, goalRepo, txtManager)

	health := handlers.HealthCheck{}
	userHandler := handlers.UserHandler{UserService: userService}
	goalHandler := handlers.GoalHandler{GoalService: goalService}
	authHandler := handlers.AuthHandler{UserService: userService}

	return &api.HandlerRegistry{HealthHandler: health, UserHandler: userHandler, GoalHandler: goalHandler, AuthHandler: authHandler}, nil
}
