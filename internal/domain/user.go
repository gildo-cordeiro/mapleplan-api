package domain

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/contract"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:text;not null" json:"-"`
	FirstName    string `gorm:"type:varchar(100)" json:"firstName,omitempty"`
	LastName     string `gorm:"type:varchar(100)" json:"lastName,omitempty"`

	Transaction *Transaction `gorm:"foreignKey:UserID" json:"transaction,omitempty"`
	Task        *Task        `gorm:"foreignKey:UserID" json:"task,omitempty"`
	Goal        *Goal        `gorm:"foreignKey:UserID" json:"goal,omitempty"`
}

func NewUser(email, passwordHash, firstName, lastName string, transaction *Transaction, task *Task, goal *Goal) (*User, error) {
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
