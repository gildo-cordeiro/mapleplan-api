package bootstrap

import (
	"os"

	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDsn   string
	JwtSecret     string
	StorageURL    string
	StorageKey    string
	StorageBucket string
	AppEnv        string
}

func LoadConfig() *Config {
	utils.InitLogger()

	err := godotenv.Load()
	if err != nil {
		utils.Log.Logger.Printf("Error loading .env file: %v", err)
	}
	utils.Log.Info("Starting MaplePlan API...")

	return &Config{
		DatabaseDsn:   os.Getenv("DATABASE_DSN"),
		JwtSecret:     os.Getenv("JWT_SECRET"),
		StorageURL:    os.Getenv("STORAGE_URL"),
		StorageKey:    os.Getenv("STORAGE_KEY"),
		StorageBucket: os.Getenv("STORAGE_BUCKET"),
		AppEnv:        os.Getenv("APP_ENV"),
	}
}
