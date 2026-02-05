package user

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/user/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
)

type User struct {
	domain.Base

	Email        string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string  `gorm:"type:text;not null" json:"-"`
	FirstName    string  `gorm:"type:varchar(100)" json:"firstName,omitempty"`
	LastName     string  `gorm:"type:varchar(100)" json:"lastName,omitempty"`
	Phone        *string `gorm:"type:varchar(20)" json:"phone,omitempty"`
}

func NewUser(email, passwordHash, firstName, lastName string) (*User, error) {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
	}, nil
}

func NewFromCreateDTO(dto request.CreateUserRequest, passwordHash string) (*User, error) {
	if dto == (request.CreateUserRequest{}) {
		return nil, utils.ErrInvalidInput
	}
	return NewUser(dto.Email, passwordHash, "", "")
}

func NewFromUpdateOnboardingDTO(dto request.UpdateUserOnboardingRequest, user *User) (*User, error) {
	if dto.FirstName == "" && dto.LastName == "" {
		return nil, utils.ErrNoFieldsToUpdate
	}
	return NewUser(user.Email, user.PasswordHash, dto.FirstName, dto.LastName)
}
