package services

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/contract"
	user "github.com/gildo-cordeiro/mapleplan-api/internal/domain"
	"github.com/gildo-cordeiro/mapleplan-api/internal/internalErrors"
	"github.com/gildo-cordeiro/mapleplan-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindByEmailAndPass(email, pass string) (*user.User, error)
	RegisterUser(contract.CreateNewUserDto) (string, error)
}

type userService struct {
	repo repository.Repository
}

func NewUserService(r repository.Repository) UserService {
	return &userService{repo: r}
}

func (s *userService) FindByEmailAndPass(email, pass string) (*user.User, error) {
	hashed, err := hashPassword(pass)
	if err != nil {
		return &user.User{}, err
	}

	if checkPasswordHash(hashed, pass) {
		userFounded, err := s.repo.FindByEmailAndPass(email, hashed)
		if err != nil {
			return &user.User{}, internalErrors.ErrInvalidCredentials
		}
		return userFounded, nil
	}

	return &user.User{}, internalErrors.ErrInvalidCredentials
}

func (s *userService) RegisterUser(newUser contract.CreateNewUserDto) (string, error) {
	hashed, err := hashPassword(newUser.Password)
	if err != nil {
		return "", err
	}
	user, err := user.NewFromDTO(newUser, hashed)
	if err != nil {
		return "", err
	}

	id, err := s.repo.Save(user)
	if err != nil {
		return "", err
	}

	return id, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", internalErrors.ErrInternal
	}
	return string(hashed), nil
}

func checkPasswordHash(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
