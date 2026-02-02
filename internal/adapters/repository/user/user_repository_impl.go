package user

import (
	"errors"
	"strings"

	userDomain "github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	repoPort "github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repoPort.UserRepository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Save(u *userDomain.User) (string, error) {
	if err := r.db.Create(u).Error; err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return "", utils.ErrAlreadyExists
		}
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "duplicate key value") || strings.Contains(lower, "sqlstate 23505") {
			return "", utils.ErrAlreadyExists
		}
		return "", err
	}
	return u.ID, nil
}

func (r *RepositoryImpl) FindByEmail(email string) (*userDomain.User, error) {
	var u userDomain.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *RepositoryImpl) FindByID(id string) (*userDomain.User, error) {
	var u userDomain.User
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *RepositoryImpl) Update(id string, u *userDomain.User) error {
	if err := r.db.Where("id = ?", id).Updates(u).Error; err != nil {
		return err
	}
	return nil
}
