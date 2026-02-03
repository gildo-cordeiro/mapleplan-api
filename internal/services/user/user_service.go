package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/couple"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	userRepository   repositories.UserRepository
	coupleRepository repositories.CoupleRepository
	txManager        ports.TransactionManager
}

func NewUserService(r repositories.UserRepository, c repositories.CoupleRepository, txManager ports.TransactionManager) services.UserService {
	return &ServiceImpl{userRepository: r, coupleRepository: c, txManager: txManager}
}

func (s *ServiceImpl) FindByEmailAndPass(email, pass string) (*user.User, error) {
	var found *user.User

	err := s.txManager.WithTransaction(context.Background(), func(txCtx context.Context) error {
		userFounded, err := s.userRepository.FindByEmail(txCtx, email)
		if err != nil {
			utils.Log.Errorf("error finding user by email: %v", err)
			return utils.ErrInternal
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userFounded.PasswordHash), []byte(pass)); err != nil {
			return utils.ErrInvalidCredentials
		}

		found = userFounded
		return nil
	})

	if err != nil {
		return nil, err
	}

	return found, nil
}

func (s *ServiceImpl) RegisterUser(newUser contract.CreateNewUserDto) (string, error) {
	if email := newUser.Email; email != "" {
		utils.Log.Infof("RegisterUser started for email=%s", email)
	} else {
		utils.Log.Infof("RegisterUser started")
	}

	hashed, err := hashPassword(newUser.Password)
	if err != nil {
		utils.Log.Errorf("error hashing password for email=%s: %v", newUser.Email, err)
		return "", err
	}

	userObj, err := user.NewFromCreateDTO(newUser, hashed)
	if err != nil {
		utils.Log.Errorf("error creating user object from DTO for email=%s: %v", newUser.Email, err)
		return "", err
	}

	id, err := s.userRepository.Save(userObj)
	if err != nil {
		utils.Log.Errorf("error saving user to repo for email=%s: %v", newUser.Email, err)
		return "", err
	}

	utils.Log.Infof("user registered successfully email=%s id=%s", newUser.Email, id)

	return id, nil
}

func (s *ServiceImpl) UpdateOnboarding(ctx context.Context, userId string, dto contract.UpdateUserOnboardingDto) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		userFounded, err := s.userRepository.FindByID(txCtx, userId)
		if err != nil {
			utils.Log.Errorf("error finding user by id: %v", err)
			return utils.ErrInternal
		}

		updatedUser, err := user.NewFromUpdateOnboardingDTO(dto, userFounded)
		if err != nil {
			utils.Log.Errorf("error creating user from update onboarding dto: %v", err)
			return utils.ErrInternal
		}

		if dto.PartnerEmail != "" {
			partner, err := s.userRepository.FindByEmail(txCtx, dto.PartnerEmail)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					utils.Log.Infof("partner not found, finishing onboarding without partner email=%s", dto.PartnerEmail)
				} else {
					utils.Log.Errorf("error finding partner by email: %v", err)
					return utils.ErrInternal
				}
			} else {
				// Partner found, create couple and link both users
				coupleName := dto.FirstName + " & " + partner.FirstName
				if coupleName == "" {
					coupleName = fmt.Sprintf("%s & %s", updatedUser.FirstName, partner.FirstName)
				}
				c := &couple.Couple{Name: coupleName}
				if err := s.coupleRepository.Save(txCtx, c); err != nil {
					utils.Log.Errorf("error creating couple: %v", err)
					return utils.ErrInternal
				}

				partner.CoupleID = &c.ID
				if err := s.userRepository.Update(txCtx, partner.ID, partner); err != nil {
					utils.Log.Errorf("error updating partner with couple id: %v", err)
					return utils.ErrInternal
				}

				updatedUser.CoupleID = &c.ID
			}
		}

		if err := s.userRepository.Update(txCtx, userId, updatedUser); err != nil {
			utils.Log.Errorf("error updating user onboarding: %v", err)
			return utils.ErrInternal
		}

		return nil
	})
}

func (s *ServiceImpl) SearchPartnerByName(userID string, name string) (contract.PartnersListDto, error) {
	users, err := s.userRepository.SearchByName(userID, name)
	if err != nil {
		utils.Log.Errorf("error searching users by name: %v", err)
		return contract.PartnersListDto{}, utils.ErrInternal
	}

	partners := make([]contract.Partner, 0, len(users))
	for _, u := range users {
		p := contract.Partner{
			Name:  u.FirstName + " " + u.LastName,
			Email: u.Email,
		}
		partners = append(partners, p)
	}

	return contract.PartnersListDto{Partners: partners}, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
