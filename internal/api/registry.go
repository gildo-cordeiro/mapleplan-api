package api

import "github.com/gildo-cordeiro/mapleplan-api/internal/handlers"

type HandlerRegistry struct {
	HealthHandler handlers.HealthCheck
	UserHandler   handlers.UserHandler
}
