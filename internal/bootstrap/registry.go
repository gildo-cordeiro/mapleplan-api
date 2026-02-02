package bootstrap

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/api"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/database"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers/auth"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers/user"
	usersRepo "github.com/gildo-cordeiro/mapleplan-api/internal/adapters/repository/user"
	userServicePkg "github.com/gildo-cordeiro/mapleplan-api/internal/services/user"
)

func Build() (*api.HandlerRegistry, error) {
	db, err := database.NewGormDB()
	if err != nil {
		return nil, err
	}

	userRepo := usersRepo.NewGormUserRepository(db)

	userService := userServicePkg.NewUserService(userRepo)

	health := handlers.HealthCheck{}
	userHandler := user.Handler{UserService: userService}
	authHandler := auth.Handler{UserService: userService}

	return &api.HandlerRegistry{HealthHandler: health, UserHandler: userHandler, AuthHandler: authHandler}, nil
}
