package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/api"
	"github.com/gildo-cordeiro/mapleplan-api/internal/bootstrap"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	utils.InitLogger()

	if err := godotenv.Load(); err != nil {
		utils.Log.Infof("No .env file found or failed to load, continuing with environment variables: %v", err)
	} else {
		utils.Log.Info(".env file loaded successfully")
	}

	utils.Log.Info("Starting MaplePlan API...")

	registry, err := bootstrap.Build()
	if err != nil {
		utils.Log.Fatalf("failed to build app: %v", err)
	}

	router := api.RegisterRoutes(registry)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	utils.Log.Info("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		utils.Log.Fatalf("Could not start server: %v", err)
	}
}
