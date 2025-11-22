package services

import (
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract"
	userDomain "github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
)

type UserService interface {
	FindByEmailAndPass(email, pass string) (*userDomain.User, error)
	RegisterUser(contract.CreateNewUserDto) (string, error)
}
