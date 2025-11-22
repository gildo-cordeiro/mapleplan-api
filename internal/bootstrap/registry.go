package bootstrap

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/api"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/database"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
	usersRepo "github.com/gildo-cordeiro/mapleplan-api/internal/adapters/repository/users"
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
	userHandler := handlers.UserHandler{UserService: userService}

	return &api.HandlerRegistry{HealthHandler: health, UserHandler: userHandler}, nil
}
