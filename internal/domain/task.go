package domain

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	UserID uint

	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text;not null"`
	DueDate     time.Time
	IsCompleted bool `gorm:"type:boolean;not null"`
}
