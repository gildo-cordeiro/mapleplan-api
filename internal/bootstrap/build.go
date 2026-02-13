package bootstrap

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/api"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/database"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/repository"
	"github.com/gildo-cordeiro/mapleplan-api/internal/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	storagego "github.com/supabase-community/storage-go"
	"gorm.io/gorm"
)

func Build() (*api.HandlerBuilder, error) {
	cfg := LoadConfig()
	db, err := database.NewGormDB()
	if err != nil {
		return nil, err
	}

	storageClient, storageBucket, storageErr := buildStorage(cfg, db)
	if storageErr != nil {
		utils.Log.Warnf("buildStorage returned error (non-fatal in this context): %v", storageErr)
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
		StorageBucket: storageBucket,
	}, nil
}

func buildStorage(cfg *Config, db *gorm.DB) (*storagego.Client, *storagego.Bucket, error) {
	if cfg == nil {
		return nil, nil, fmt.Errorf("nil config")
	}
	// generate token for storage
	token, err := GenerateLocalToken(cfg.JwtSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("generating storage token: %w", err)
	}

	storageClient := storagego.NewClient(cfg.StorageURL, token, nil)
	bucket, err := storageClient.GetBucket(cfg.StorageBucket)

	if err != nil {
		utils.Log.Warnf("Storage bucket '%s' not found (GetBucket error: %v), attempting to create it...", cfg.StorageBucket, err)
		const maxAttempts = 5
		var createErr error
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rlsTried := false
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			bucket, createErr = storageClient.CreateBucket(cfg.StorageBucket, storagego.BucketOptions{Public: true})
			if createErr == nil {
				utils.Log.Infof("Successfully created storage bucket '%s' (attempt %d)", cfg.StorageBucket, attempt)
				createErr = nil
				break
			}

			// try to handle RLS in local env once
			if !rlsTried {
				if disableRLSIfRowLevelError(cfg, db, createErr) {
					rlsTried = true
					time.Sleep(500 * time.Millisecond)
					continue
				}
			}

			utils.Log.Warnf("Failed to create bucket '%s' (attempt %d/%d): %v", cfg.StorageBucket, attempt, maxAttempts, createErr)
			// exponential backoff with jitter
			sleep := time.Duration(1<<uint(attempt-1)) * time.Second
			jitter := time.Duration(r.Intn(500)) * time.Millisecond
			time.Sleep(sleep + jitter)
		}
		if createErr != nil {
			return nil, nil, fmt.Errorf("unable to create storage bucket '%s' after %d attempts: %w", cfg.StorageBucket, maxAttempts, createErr)
		}
	} else {
		utils.Log.Infof("Storage bucket '%s' already exists", cfg.StorageBucket)
	}

	return storageClient, &bucket, nil
}

func disableRLSIfRowLevelError(cfg *Config, db *gorm.DB, createErr error) bool {
	if cfg == nil || createErr == nil {
		return false
	}
	if cfg.AppEnv != "local" {
		utils.Log.Debugf("DisableRLSIfRowLevelError skipped because AppEnv=%s", cfg.AppEnv)
		return false
	}
	msg := createErr.Error()
	if !(strings.Contains(msg, "row-level security") || strings.Contains(msg, "violates row-level")) {
		return false
	}
	utils.Log.Warnf("Detected RLS error while creating bucket: %v; attempting to disable RLS on storage.buckets (local only)", createErr)
	if db == nil {
		utils.Log.Errorf("DB connection is nil; cannot disable RLS")
		return false
	}
	if execErr := db.Exec("ALTER TABLE IF EXISTS storage.buckets DISABLE ROW LEVEL SECURITY;").Error; execErr != nil {
		utils.Log.Errorf("Failed to disable RLS via DB: %v", execErr)
		return false
	}
	utils.Log.Infof("Disabled RLS on storage.buckets via DB; will allow CreateBucket to retry")
	// small pause to let DB propagate
	time.Sleep(500 * time.Millisecond)
	return true
}
