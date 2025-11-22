package user

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/goal"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/task"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/transaction"
)

type User struct {
	domain.Base

	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:text;not null" json:"-"`
	FirstName    string `gorm:"type:varchar(100)" json:"firstName,omitempty"`
	LastName     string `gorm:"type:varchar(100)" json:"lastName,omitempty"`

	Transaction *transaction.Transaction `gorm:"foreignKey:UserID" json:"transaction,omitempty"`
	Task        *task.Task               `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
	Goal        *goal.Goal               `gorm:"foreignKey:UserID" json:"goal,omitempty"`
}

func NewUser(email, passwordHash, firstName, lastName string, transaction *transaction.Transaction, task *task.Task, goal *goal.Goal) (*User, error) {
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
