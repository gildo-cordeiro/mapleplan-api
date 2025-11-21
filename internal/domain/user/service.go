package user

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/internalErrors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repository Repository
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", internalErrors.ErrInternal
	}
	return string(hashed), nil
}

func CheckPasswordHash(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func (service *Service) RegisterUser(newUser contract.CreateNewUserDto) (string, error) {
	hashed, err := HashPassword(newUser.Password)
	if err != nil {
		return "", err
	}
	user, err := NewFromDTO(newUser, hashed)
	if err != nil {
		return "", err
	}

	id, err := service.Repository.Save(user)
	if err != nil {
		return "", err
	}

	return id, nil
}
