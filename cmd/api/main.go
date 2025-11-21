package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/api"
	"github.com/gildo-cordeiro/mapleplan-api/internal/database"
	"github.com/gildo-cordeiro/mapleplan-api/internal/handlers"
	"github.com/gildo-cordeiro/mapleplan-api/internal/repository"
	"github.com/gildo-cordeiro/mapleplan-api/internal/services"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// --- 1. Setup logging and load environment variables ---
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	if err := godotenv.Load(); err != nil {
		logrus.Infof("No .env file found or failed to load, continuing with environment variables: %v", err)
	} else {
		logrus.Info(".env file loaded successfully")
	}

	logrus.Info("Starting MaplePlan API...")
	db, err := database.NewGormDB()
	if err != nil {
		logrus.Fatalf("Failed to connect to the database: %v", err)
	}

	// --- 2. Initialize application components (repositories, services, handlers) ---

	// 2.1. Create Repositories (inject database connection)
	userRepo := repository.NewGormUserRepository(db)

	// 2.2. Create Services (inject repositories)
	userService := services.NewUserService(userRepo)

	// 2.3. Create Handlers (inject services)
	userHandler := handlers.UserHandler{UserService: userService}

	// 2.4. Create Handler Registry to pass to the router
	registry := &api.HandlerRegistry{UserHandler: userHandler}

	// --- 3. Initialize and configure the HTTP server ---
	router := api.RegisterRoutes(registry)

	// Optional: Add global middlewares (CORS, Logging, etc.)
	// router.Use(middlewares.LoggingMiddleware)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logrus.Info("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		logrus.Fatalf("Could not start server: %v", err)
	}
}
