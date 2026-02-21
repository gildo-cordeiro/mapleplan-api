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
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/storage"
	"github.com/gildo-cordeiro/mapleplan-api/internal/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
)

func BuildS3Storage(cfg *Config) (services.StorageService, error) {
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
