package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/profile"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
	"github.com/gildo-cordeiro/mapleplan-api/internal/dto/user/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/dto/user/response"
	"github.com/gildo-cordeiro/mapleplan-api/internal/ports"
	"github.com/gildo-cordeiro/mapleplan-api/internal/ports/repositories"
	"github.com/gildo-cordeiro/mapleplan-api/internal/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	userRepository          repositories.UserRepository
	profileRepository       repositories.ProfileRepository
	profileMemberRepository repositories.ProfileMemberRepository
	txManager               ports.TransactionManager
}

func NewUserService(r repositories.UserRepository, p repositories.ProfileRepository, m repositories.ProfileMemberRepository, txManager ports.TransactionManager) services.UserService {
	return &UserServiceImpl{userRepository: r, profileRepository: p, profileMemberRepository: m, txManager: txManager}
}

func (s *UserServiceImpl) FindByEmailAndPass(email, pass string) (*user.User, error) {
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

func (s *UserServiceImpl) RegisterUser(newUser request.CreateUserRequest) (string, error) {
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

func (s *UserServiceImpl) UpdateOnboarding(ctx context.Context, userId string, dto request.UpdateUserOnboardingRequest) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		userFounded, updatedUser, err := s.prepareUserUpdate(txCtx, userId, dto)
		if err != nil {
			return err
		}

		partner, err := s.findPartner(txCtx, dto.PartnerEmail)
		if err != nil {
			return err
		}

		profileID, err := s.createImmigrationProfile(txCtx, userId, userFounded.FirstName, partner)
		if err != nil {
			return err
		}

		if err := s.addProfileMembers(txCtx, profileID, userId, partner); err != nil {
			return err
		}

		if err := s.userRepository.Update(txCtx, userId, updatedUser); err != nil {
			utils.Log.Errorf("error updating user onboarding: %v", err)
			return utils.ErrInternal
		}

		return nil
	})
}

func (s *UserServiceImpl) SearchPartnerByName(userID string, name string) (response.PartnersListResponse, error) {
	users, err := s.userRepository.SearchByName(userID, name)
	if err != nil {
		utils.Log.Errorf("error searching users by name: %v", err)
		return response.PartnersListResponse{}, utils.ErrInternal
	}

	partners := make([]response.Partner, 0, len(users))
	for _, u := range users {
		p := response.Partner{
			Name:  u.FirstName + " " + u.LastName,
			Email: u.Email,
		}
		partners = append(partners, p)
	}

	return response.PartnersListResponse{Partners: partners}, nil
}

func (s *UserServiceImpl) GetCompleteUser(ctx context.Context, userID string) (*response.UserWithCoupleResponse, error) {
	var userFounded response.UserWithCoupleResponse

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		u, err := s.userRepository.FindByID(txCtx, userID)
		if err != nil {
			utils.Log.Errorf("error finding user by id: %v", err)
			return utils.ErrInternal
		}

		if u == nil {
			return utils.ErrRecordNotFound
		}

		userFounded = response.UserWithCoupleResponse{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: &u.FirstName,
			LastName:  &u.LastName,
			Phone:     u.Phone,
		}

		// Find profile members for this user
		members, err := s.profileMemberRepository.FindByUserID(txCtx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			utils.Log.Errorf("error finding profile members by user id: %v", err)
			return utils.ErrInternal
		}

		// If user has profile memberships, use the first one (typically should have one main profile)
		if len(members) > 0 && members[0].Profile != nil {
			profile := members[0].Profile
			userFounded.CoupleID = &profile.ID // Use ProfileID as CoupleID for backward compatibility

			// Find other members of this profile to get partner info
			profileMembers, err := s.profileMemberRepository.FindByProfileID(txCtx, profile.ID)
			if err != nil {
				utils.Log.Errorf("error finding profile members: %v", err)
				return utils.ErrInternal
			}

			for _, member := range profileMembers {
				if member.UserID != userID {
					partner, err := s.userRepository.FindByID(txCtx, member.UserID)
					if err != nil {
						utils.Log.Errorf("error finding partner by id: %v", err)
						continue
					}
					userFounded.PartnerEmail = &partner.Email
					userFounded.PartnerFirstName = &partner.FirstName
					userFounded.PartnerLastName = &partner.LastName
					userFounded.PartnerId = &partner.ID
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &userFounded, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *UserServiceImpl) prepareUserUpdate(ctx context.Context, userId string, dto request.UpdateUserOnboardingRequest) (*user.User, *user.User, error) {
	userFounded, err := s.userRepository.FindByID(ctx, userId)
	if err != nil {
		utils.Log.Errorf("error finding user by id: %v", err)
		return nil, nil, utils.ErrInternal
	}

	updatedUser, err := user.NewFromUpdateOnboardingDTO(dto, userFounded)
	if err != nil {
		utils.Log.Errorf("error creating user from update onboarding dto: %v", err)
		return nil, nil, utils.ErrInternal
	}

	return userFounded, updatedUser, nil
}

func (s *UserServiceImpl) findPartner(ctx context.Context, partnerEmail string) (*user.User, error) {
	if partnerEmail == "" {
		return nil, nil
	}

	partner, err := s.userRepository.FindByEmail(ctx, partnerEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Log.Infof("partner not found, finishing onboarding without partner email=%s", partnerEmail)
			return nil, nil
		}
		utils.Log.Errorf("error finding partner by email: %v", err)
		return nil, utils.ErrInternal
	}

	return partner, nil
}

func (s *UserServiceImpl) createImmigrationProfile(ctx context.Context, userId, userFirstName string, partner *user.User) (string, error) {
	var profileName string
	if partner != nil {
		profileName = fmt.Sprintf("%s & %s", userFirstName, partner.FirstName)
	} else {
		profileName = fmt.Sprintf("%s's Profile", userFirstName)
	}

	p := &profile.ImmigrationProfile{
		UserID: userId,
		Name:   profileName,
	}

	if err := s.profileRepository.Save(ctx, p); err != nil {
		utils.Log.Errorf("error creating profile: %v", err)
		return "", utils.ErrInternal
	}

	return p.ID, nil
}

func (s *UserServiceImpl) addProfileMembers(ctx context.Context, profileID, userId string, partner *user.User) error {
	// Add primary member
	primaryMember, err := profile.NewProfileMember(profileID, userId, profile.RolePrimary)
	if err != nil {
		utils.Log.Errorf("error creating primary member: %v", err)
		return utils.ErrInternal
	}

	if err := s.profileMemberRepository.Save(ctx, primaryMember); err != nil {
		utils.Log.Errorf("error saving primary member: %v", err)
		return utils.ErrInternal
	}

	// Add spouse member (if exists)
	if partner != nil {
		spouseMember, err := profile.NewProfileMember(profileID, partner.ID, profile.RoleSpouse)
		if err != nil {
			utils.Log.Errorf("error creating spouse member: %v", err)
			return utils.ErrInternal
		}

		if err := s.profileMemberRepository.Save(ctx, spouseMember); err != nil {
			utils.Log.Errorf("error saving spouse member: %v", err)
			return utils.ErrInternal
		}
	}

	return nil
}
