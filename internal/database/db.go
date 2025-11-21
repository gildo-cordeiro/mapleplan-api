package database

import (
	"fmt"
	"time"

	finance2 "github.com/gildo-cordeiro/mapleplan-api/internal/domain/finance"
	"github.com/gildo-cordeiro/mapleplan-api/internal/domain/tasks"
	"github.com/gildo-cordeiro/mapleplan-api/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// The ConnectDB function receives the DSN and returns the DB instance.
func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro fatal ao conectar ao PostgreSQL: %w", err)
	}

	err = db.AutoMigrate(&user.User{}, &tasks.Task{}, &finance2.Transaction{}, &finance2.Goal{})
	if err != nil {
		return nil, err
	}

	// Advanced configuration of the connection pool (sql/database)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter a instância PSQL: %w", err)
	}

	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
