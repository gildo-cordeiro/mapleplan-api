package user

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	userRepository "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	userServicePort "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo userRepository.UserRepository
}

func NewUserService(r userRepository.UserRepository) userServicePort.UserService {
	return &UserServiceImpl{repo: r}
}

func (s *UserServiceImpl) FindByEmailAndPass(email, pass string) (*user.User, error) {
	userFounded, err := s.repo.FindByEmail(email)
	if err != nil {
		if Log := utils.Log; Log != nil {
			Log.Errorf("error finding user by email: %v", err)
		}
		return &user.User{}, utils.ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFounded.PasswordHash), []byte(pass)); err != nil {
		return &user.User{}, utils.ErrInvalidCredentials
	}

	return userFounded, nil
}

func (s *UserServiceImpl) RegisterUser(newUser contract.CreateNewUserDto) (string, error) {
	hashed, err := hashPassword(newUser.Password)
	if err != nil {
		return "", err
	}
	userObj, err := user.NewFromDTO(newUser, hashed)
	if err != nil {
		return "", err
	}

	id, err := s.repo.Save(userObj)
	if err != nil {
		return "", err
	}

	return id, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkPasswordHash(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
