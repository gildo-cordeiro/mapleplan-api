package goal

import (
	"errors"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/couple"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
)

type Goal struct {
	domain.Base

	CoupleID *string `gorm:"type:uuid;index" json:"coupleId,omitempty"`
	UserId   *string `gorm:"type:uuid;index" json:"userId,omitempty"`

	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Status      Status    `gorm:"type:varchar(20);not null" json:"status"`
	Progress    int       `gorm:"type:integer;not null" json:"progress"`
	DueDate     time.Time `json:"type:date" json:"dueDate"`
	Phase       Phase     `gorm:"type:varchar(20);not null" json:"phase"`
	Priority    Priority  `gorm:"type:varchar(20);not null" json:"priority"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`

	Couple *couple.Couple `gorm:"foreignKey:CoupleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"couple,omitempty"`
	User   *user.User     `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
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

	g.UserId = goal.UserId
	g.CoupleID = goal.CoupleID

	g.User = nil
	g.Couple = nil

	return nil
}
