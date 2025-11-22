package users

import (
	userDomain "github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	repoPort "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repoPort.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Save(u *userDomain.User) (string, error) {
	if err := r.db.Create(u).Error; err != nil {
		return "", err
	}
	return u.ID, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*userDomain.User, error) {
	var u userDomain.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
