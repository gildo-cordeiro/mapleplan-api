package api

import (
	handlers "github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
)

type HandlerRegistry struct {
	HealthHandler handlers.HealthCheck
	UserHandler   handlers.UserHandler
}
