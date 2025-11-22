package database

import (
	"fmt"
	"os"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/goal"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/task"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/transaction"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_DSN not set; check the .env file or environment variables")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("fatal error connecting to PostgreSQL: %w", err)
	}

	// Auto-migrate domain models
	err = db.AutoMigrate(&user.User{}, &task.Task{}, &transaction.Transaction{}, &goal.Goal{})
	if err != nil {
		return nil, fmt.Errorf("error migrating the database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error obtaining the database instance: %w", err)
	}

	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
