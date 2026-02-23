package task

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/goal"
)

type Task struct {
	models.Base

	GoalID string `gorm:"type:uuid;not null;index" json:"goalId"`

	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	Completed   bool       `gorm:"not null;default:false;index" json:"completed"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`

	Goal *goal.Goal `gorm:"foreignKey:GoalID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"goal,omitempty"`
}
