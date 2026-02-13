package bootstrap

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/api"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/database"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/repository"
	"github.com/gildo-cordeiro/mapleplan-api/internal/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	storagego "github.com/supabase-community/storage-go"
)

func Build() (*api.HandlerBuilder, error) {
	cfg := LoadConfig()
	db, err := database.NewGormDB()
	if err != nil {
		return nil, err
	}

	var storageClient interface{}
	if cfg.AppEnv == "local" {
		storageClient, err = buildMinIOStorage(cfg)
	} else {
		storageClient, err = buildSupabaseStorage(cfg)
	}
	if err != nil {
		utils.Log.Warnf("buildStorage returned error (non-fatal in this context): %v", err)
	}

	userRepo := repository.NewGormUserRepository(db)
	goalRepo := repository.NewGormGoalRepository(db)
	coupleRepo := repository.NewGormCoupleRepository(db)
	txtManager := repository.NewGormTransactionManager(db)

	userService := services.NewUserService(userRepo, coupleRepo, txtManager)
	goalService := services.NewGoalService(userRepo, goalRepo, coupleRepo, txtManager)

	health := handlers.HealthCheck{}
	userHandler := handlers.UserHandler{UserService: userService}
	goalHandler := handlers.GoalHandler{GoalService: goalService}
	authHandler := handlers.AuthHandler{UserService: userService}

	return &api.HandlerBuilder{
		HealthHandler: health,
		UserHandler:   userHandler,
		GoalHandler:   goalHandler,
		AuthHandler:   authHandler,
		StorageClient: storageClient,
	}, nil
}

func buildMinIOStorage(cfg *Config) (*minio.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil config")
	}

	minioClient, err := minio.New(cfg.StorageURL, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.StorageKey, cfg.StorageSecret, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("initializing MinIO client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, cfg.StorageBucket)
	if err != nil {
		return nil, fmt.Errorf("checking bucket existence: %w", err)
	}

	if !exists {
		utils.Log.Warnf("MinIO bucket '%s' not found, attempting to create it...", cfg.StorageBucket)
		const maxAttempts = 5
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			err = minioClient.MakeBucket(ctx, cfg.StorageBucket, minio.MakeBucketOptions{Region: "us-east-1"})
			if err == nil {
				utils.Log.Infof("Successfully created MinIO bucket '%s' (attempt %d)", cfg.StorageBucket, attempt)
				break
			}

			utils.Log.Warnf("Failed to create bucket '%s' (attempt %d/%d): %v", cfg.StorageBucket, attempt, maxAttempts, err)
			sleep := time.Duration(1<<uint(attempt-1)) * time.Second
			jitter := time.Duration(r.Intn(500)) * time.Millisecond
			time.Sleep(sleep + jitter)
		}
		if err != nil {
			return nil, fmt.Errorf("unable to create MinIO bucket '%s' after %d attempts: %w", cfg.StorageBucket, maxAttempts, err)
		}
	} else {
		utils.Log.Infof("MinIO bucket '%s' already exists", cfg.StorageBucket)
	}

	return minioClient, nil
}

func buildSupabaseStorage(cfg *Config) (*storagego.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil config")
	}

	token, err := GenerateLocalToken(cfg.JwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generating storage token: %w", err)
	}

	storageClient := storagego.NewClient(cfg.StorageURL, token, nil)
	_, err = storageClient.GetBucket(cfg.StorageBucket)

	if err != nil {
		utils.Log.Warnf("Storage bucket '%s' not found, attempting to create it...", cfg.StorageBucket)
		const maxAttempts = 5
		var createErr error
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			_, createErr = storageClient.CreateBucket(cfg.StorageBucket, storagego.BucketOptions{Public: true})
			if createErr == nil {
				utils.Log.Infof("Successfully created storage bucket '%s' (attempt %d)", cfg.StorageBucket, attempt)
				break
			}

			utils.Log.Warnf("Failed to create bucket '%s' (attempt %d/%d): %v", cfg.StorageBucket, attempt, maxAttempts, createErr)
			sleep := time.Duration(1<<uint(attempt-1)) * time.Second
			jitter := time.Duration(r.Intn(500)) * time.Millisecond
			time.Sleep(sleep + jitter)
		}
		if createErr != nil {
			return nil, fmt.Errorf("unable to create storage bucket '%s' after %d attempts: %w", cfg.StorageBucket, maxAttempts, createErr)
		}
	} else {
		utils.Log.Infof("Storage bucket '%s' already exists", cfg.StorageBucket)
	}

	return storageClient, nil
}
