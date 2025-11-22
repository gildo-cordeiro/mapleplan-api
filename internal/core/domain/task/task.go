package task

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
)

type Task struct {
	domain.Base
	UserID uint

	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text;not null"`
	DueDate     time.Time
	IsCompleted bool `gorm:"type:boolean;not null"`
}
