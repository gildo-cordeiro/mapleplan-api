package user

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	userRepository "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	userServicePort "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type ServiceImpl struct {
	repo userRepository.UserRepository
}

func NewUserService(r userRepository.UserRepository) userServicePort.UserService {
	return &ServiceImpl{repo: r}
}

func (s *ServiceImpl) FindByEmailAndPass(email, pass string) (*user.User, error) {
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

func (s *ServiceImpl) RegisterUser(newUser contract.CreateNewUserDto) (string, error) {
	if Log := utils.Log; Log != nil {
		if email := newUser.Email; email != "" {
			Log.Infof("RegisterUser started for email=%s", email)
		} else {
			Log.Infof("RegisterUser started")
		}
	}

	hashed, err := hashPassword(newUser.Password)
	if err != nil {
		if Log := utils.Log; Log != nil {
			Log.Errorf("error hashing password for email=%s: %v", newUser.Email, err)
		}
		return "", err
	}

	userObj, err := user.NewFromCreateDTO(newUser, hashed)
	if err != nil {
		if Log := utils.Log; Log != nil {
			Log.Errorf("error creating user object from DTO for email=%s: %v", newUser.Email, err)
		}
		return "", err
	}

	id, err := s.repo.Save(userObj)
	if err != nil {
		if Log := utils.Log; Log != nil {
			Log.Errorf("error saving user to repo for email=%s: %v", newUser.Email, err)
		}
		return "", err
	}

	if Log := utils.Log; Log != nil {
		Log.Infof("user registered successfully email=%s id=%s", newUser.Email, id)
	}

	return id, nil
}

func (s *ServiceImpl) UpdateOnboarding(userId string, dto contract.UpdateUserOnboardingDto) error {
	userFounded, err := s.repo.FindByID(userId)
	if err != nil {
		if Log := utils.Log; Log != nil {
			Log.Errorf("error finding user by id: %v", err)
		}
		return utils.ErrInternal
	}

	updatedUser, err := user.NewFromUpdateOnboardingDTO(dto, userFounded)
	if err != nil {
		if Log := utils.Log; Log != nil {
			Log.Errorf("error creating user from update onboarding dto: %v", err)
		}
		return utils.ErrInternal
	}

	if err := s.repo.Update(userId, updatedUser); err != nil {
		if Log := utils.Log; Log != nil {
			Log.Errorf("error updating user onboarding: %v", err)
		}
		return utils.ErrInternal
	}

	return nil
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
