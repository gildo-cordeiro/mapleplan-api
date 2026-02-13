package api

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
	storagego "github.com/supabase-community/storage-go"
)

type HandlerBuilder struct {
	HealthHandler handlers.HealthCheck
	UserHandler   handlers.UserHandler
	GoalHandler   handlers.GoalHandler
	AuthHandler   handlers.AuthHandler

	StorageClient *storagego.Client
	StorageBucket *storagego.Bucket
}
