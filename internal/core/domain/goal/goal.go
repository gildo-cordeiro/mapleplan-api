package goal

import (
	"errors"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/couple"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	"github.com/shopspring/decimal"
)

type Goal struct {
	domain.Base

	CoupleID *string `gorm:"type:uuid;index" json:"coupleId,omitempty"`
	UserId   *string `gorm:"type:uuid;index" json:"userId,omitempty"`

	Name          string          `gorm:"type:varchar(100);not null" json:"name"`
	Status        Status          `gorm:"type:varchar(20);not null" json:"status"`
	Progress      int             `gorm:"type:integer;not null" json:"progress"`
	TargetAmount  decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"targetAmount"`
	CurrentAmount decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"currentAmount"`
	DueDate       time.Time       `json:"type:date" json:"dueDate"`
	Phase         Phase           `gorm:"type:varchar(20);not null" json:"phase"`
	Priority      Priority        `gorm:"type:varchar(20);not null" json:"priority"`
	Description   *string         `gorm:"type:text" json:"description,omitempty"`

	Couple *couple.Couple `gorm:"foreignKey:CoupleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"couple,omitempty"`
	User   *user.User     `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}

func NewGoal(name string, status Status, progress int, targetAmount, currentAmount decimal.Decimal, dueDate time.Time, phase Phase, priority Priority, description *string, coupleID *string, userId *string) (*Goal, error) {
	return &Goal{
		Name:          name,
		Status:        status,
		Progress:      progress,
		TargetAmount:  targetAmount,
		CurrentAmount: currentAmount,
		DueDate:       dueDate,
		Phase:         phase,
		Priority:      priority,
		Description:   description,
		CoupleID:      coupleID,
		UserId:        userId,
	}, nil
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
