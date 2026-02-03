package user

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/couple"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/goal"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/task"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/transaction"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
)

type User struct {
	domain.Base

	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:text;not null" json:"-"`
	FirstName    string `gorm:"type:varchar(100)" json:"firstName,omitempty"`
	LastName     string `gorm:"type:varchar(100)" json:"lastName,omitempty"`

	CoupleID *string        `gorm:"index" json:"coupleId,omitempty"`
	Couple   *couple.Couple `gorm:"foreignKey:CoupleID" json:"couple,omitempty"`

	Transactions []*transaction.Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
	Tasks        []*task.Task               `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
	Goals        []*goal.Goal               `gorm:"foreignKey:UserID" json:"goals,omitempty"`
}

func NewUser(email, passwordHash, firstName, lastName string, coupleID *string) (*User, error) {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		CoupleID:     coupleID,
	}, nil
}

func NewFromCreateDTO(dto contract.CreateNewUserDto, passwordHash string) (*User, error) {
	if dto == (contract.CreateNewUserDto{}) {
		return nil, utils.ErrInvalidInput
	}
	return NewUser(dto.Email, passwordHash, "", "", nil)
}

func NewFromUpdateOnboardingDTO(dto contract.UpdateUserOnboardingDto, user *User) (*User, error) {
	if dto.FirstName == "" && dto.LastName == "" {
		return nil, utils.ErrNoFieldsToUpdate
	}
	return NewUser(user.Email, user.PasswordHash, dto.FirstName, dto.LastName, nil)
}
