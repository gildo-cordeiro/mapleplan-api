package bootstrap

import (
	"os"

	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDsn        string
	JwtSecret          string
	StorageBucket      string
	AppEnv             string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	S3Endpoint         string
}

func LoadConfig() *Config {
	utils.InitLogger()

	err := godotenv.Load()
	if err != nil {
		utils.Log.Logger.Printf("Error loading .env file: %v", err)
	}
	utils.Log.Info("Starting MaplePlan API...")

	return &Config{
		DatabaseDsn:        os.Getenv("DATABASE_DSN"),
		JwtSecret:          os.Getenv("JWT_SECRET"),
		StorageBucket:      os.Getenv("STORAGE_BUCKET"),
		AppEnv:             os.Getenv("APP_ENV"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Endpoint:         os.Getenv("S3_ENDPOINT"),
	}
}
