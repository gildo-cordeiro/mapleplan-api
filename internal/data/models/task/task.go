package task

import (
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/goal"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
)

type Task struct {
	domain.Base

	UserID *string `gorm:"type:uuid;index" json:"userId,omitempty"`
	GoalID *string `gorm:"type:uuid;index" json:"goalId,omitempty"`

	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	Completed   bool       `gorm:"not null;default:false" json:"completed"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`

	User *user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Goal *goal.Goal `gorm:"foreignKey:GoalID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"goal,omitempty"`
}
