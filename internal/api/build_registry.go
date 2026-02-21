package api

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/api/handlers"
)

type BuildRegistry struct {
	HealthHandler handlers.HealthCheck
	UserHandler   handlers.UserHandler
	GoalHandler   handlers.GoalHandler
	AuthHandler   handlers.AuthHandler
}
