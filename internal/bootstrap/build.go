package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/database"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/handlers"
	"github.com/gildo-cordeiro/mapleplan-api/internal/adapters/repository"
	"github.com/gildo-cordeiro/mapleplan-api/internal/api"
	ports "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/internal/services"
	"github.com/gildo-cordeiro/mapleplan-api/internal/storage"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
)

func Build() (*api.HandlerBuilder, error) {
	cfg := LoadConfig()
	db, err := database.NewGormDB()
	if err != nil {
		return nil, err
	}

	storageAdapter, err := buildS3Storage(cfg)
	if err != nil {
		utils.Log.Warnf("buildStorage returned error (non-fatal in this context): %v", err)
		return nil, err
	}

	userRepo := repository.NewGormUserRepository(db)
	goalRepo := repository.NewGormGoalRepository(db)
	coupleRepo := repository.NewGormCoupleRepository(db)
	txtManager := repository.NewGormTransactionManager(db)

	userService := services.NewUserService(userRepo, coupleRepo, txtManager)
	goalService := services.NewGoalService(userRepo, goalRepo, coupleRepo, txtManager)
	services.NewStorageService(storageAdapter)

	health := handlers.HealthCheck{}
	userHandler := handlers.UserHandler{UserService: userService}
	goalHandler := handlers.GoalHandler{GoalService: goalService}
	authHandler := handlers.AuthHandler{UserService: userService}

	return &api.HandlerBuilder{
		HealthHandler: health,
		UserHandler:   userHandler,
		GoalHandler:   goalHandler,
		AuthHandler:   authHandler,
	}, nil
}

func buildS3Storage(cfg *Config) (ports.StorageService, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil config")
	}

	// Create AWS config with static credentials
	loadOptions := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, "")),
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), loadOptions...)
	if err != nil {
		utils.Log.Errorf("loading AWS config: %v", err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		if cfg.S3Endpoint != "" {
			options.BaseEndpoint = aws.String(cfg.S3Endpoint)
			options.UsePathStyle = true
		}
	})

	// Verify bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &cfg.StorageBucket})
	if err != nil {
		utils.Log.Warnf("Failed to access S3 bucket '%s': %v", cfg.StorageBucket, err)
		if cfg.S3Endpoint != "" {
			createInput := &s3.CreateBucketInput{Bucket: &cfg.StorageBucket}
			if cfg.AWSRegion != "" && cfg.AWSRegion != "us-east-1" {
				createInput.CreateBucketConfiguration = &types.CreateBucketConfiguration{
					LocationConstraint: types.BucketLocationConstraint(cfg.AWSRegion),
				}
			}

			_, createErr := s3Client.CreateBucket(ctx, createInput)
			if createErr != nil {
				utils.Log.Warnf("Failed to create S3 bucket '%s': %v", cfg.StorageBucket, createErr)
			} else {
				utils.Log.Infof("Created S3 bucket '%s'", cfg.StorageBucket)
			}
		}
	} else {
		utils.Log.Infof("Successfully connected to S3 bucket '%s'", cfg.StorageBucket)
	}

	return storage.NewS3Storage(s3Client, cfg.AWSRegion, cfg.S3Endpoint, utils.Log.Logger), nil
}
