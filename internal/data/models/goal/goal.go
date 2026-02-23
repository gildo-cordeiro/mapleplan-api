package goal

import (
	"errors"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/profile"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
	"gorm.io/gorm"
)

type Goal struct {
	models.Base

	UserID               *string `gorm:"type:uuid;index" json:"userId,omitempty"`
	ImmigrationProfileID *string `gorm:"type:uuid;index" json:"immigrationProfileId,omitempty"` // NULL = private, filled = shared

	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Status      Status    `gorm:"type:varchar(20);not null" json:"status"`
	Progress    int       `gorm:"type:integer;not null;default:0" json:"progress"`
	DueDate     time.Time `gorm:"type:date" json:"dueDate"`
	Phase       Phase     `gorm:"type:varchar(20);not null" json:"phase"`
	Priority    Priority  `gorm:"type:varchar(20);not null" json:"priority"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`

	User               *user.User                  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	ImmigrationProfile *profile.ImmigrationProfile `gorm:"foreignKey:ImmigrationProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"immigrationProfile,omitempty"`
}

func (g *Goal) BeforeSave(tx *gorm.DB) error {
	if g.UserID != nil && g.ImmigrationProfileID != nil {
		return errors.New("goal cannot belong to both user and profile simultaneously")
	}

	if g.UserID == nil && g.ImmigrationProfileID == nil {
		return errors.New("goal must belong to either a user or a profile")
	}

	return nil
}

func (g *Goal) SetPhase(p Phase) error {
	if !IsValidPhase(p) {
		return errors.New("phase invalid")
	}
	g.Phase = p
	return nil
}

func (g *Goal) SetPriority(p Priority) error {
	if !IsValidPriority(p) {
		return errors.New("priority invalid")
	}
	g.Priority = p
	return nil
}

func (g *Goal) SetStatus(s Status) error {
	if !IsValidStatus(s) {
		return errors.New("status invalid")
	}
	g.Status = s
	return nil
}

func (g *Goal) UpdateFields(goal *Goal) error {
	if goal.UserID != nil && goal.ImmigrationProfileID != nil {
		return errors.New("goal cannot belong to both user and profile")
	}

	if err := g.SetStatus(goal.Status); err != nil {
		return err
	}
	if err := g.SetPhase(goal.Phase); err != nil {
		return err
	}
	if err := g.SetPriority(goal.Priority); err != nil {
		return err
	}

	// copy simple fields
	g.Name = goal.Name
	g.Description = goal.Description
	g.DueDate = goal.DueDate
	g.Progress = goal.Progress

	g.UserID = goal.UserID
	g.ImmigrationProfileID = goal.ImmigrationProfileID

	g.User = nil
	g.ImmigrationProfile = nil
	return nil
}
