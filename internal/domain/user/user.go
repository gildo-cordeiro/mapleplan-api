package user

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/domain/finance"
	"github.com/gildo-cordeiro/mapleplan-api/internal/domain/tasks"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:text;not null" json:"-"`
	FirstName    string `gorm:"type:varchar(100)" json:"firstName,omitempty"`
	LastName     string `gorm:"type:varchar(100)" json:"lastName,omitempty"`

	Transaction *finance.Transaction `gorm:"foreignKey:UserID" json:"transaction,omitempty"`
	Task        *tasks.Task          `gorm:"foreignKey:UserID" json:"task,omitempty"`
	Goal        *finance.Goal        `gorm:"foreignKey:UserID" json:"goal,omitempty"`
}

func NewUser(email, passwordHash, firstName, lastName string, transaction *finance.Transaction, task *tasks.Task, goal *finance.Goal) (*User, error) {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Transaction:  transaction,
		Task:         task,
		Goal:         goal,
	}, nil
}

func NewFromDTO(dto contract.CreateNewUserDto, passwordHash string) (*User, error) {
	return NewUser(dto.Email, passwordHash, dto.Name, dto.LastName, nil, nil, nil)
}
