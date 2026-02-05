package services

import (
	"context"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/user/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/user/response"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
)

type UserService interface {
	FindByEmailAndPass(email, pass string) (*user.User, error)
	RegisterUser(request.CreateUserRequest) (string, error)
	UpdateOnboarding(ctx context.Context, userId string, dto request.UpdateUserOnboardingRequest) error
	SearchPartnerByName(userID string, name string) (response.PartnersListResponse, error)
	GetCompleteUser(ctx context.Context, userID string) (*response.UserWithCoupleResponse, error)
}
