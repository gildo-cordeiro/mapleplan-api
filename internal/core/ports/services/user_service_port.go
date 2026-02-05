package services

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
)

type UserService interface {
	FindByEmailAndPass(email, pass string) (*user.User, error)
	RegisterUser(contract.CreateNewUserDto) (string, error)
	UpdateOnboarding(ctx context.Context, userId string, dto contract.UpdateUserOnboardingDto) error
	SearchPartnerByName(userID string, name string) (contract.PartnersListDto, error)
	GetCompleteUser(ctx context.Context, userID string) (*contract.UserDTO, error)
}
