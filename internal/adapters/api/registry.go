package api

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers/auth"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers/user"
)

type HandlerRegistry struct {
	HealthHandler handlers.HealthCheck
	UserHandler   user.Handler
	AuthHandler   auth.Handler
}
