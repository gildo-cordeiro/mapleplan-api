package repository

import (
	"strconv"

	user "github.com/gildo-cordeiro/mapleplan-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmailAndPass(email string, pass string) (*user.User, error)
	Save(u *user.User) (string, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) Repository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Save(u *user.User) (string, error) {
	if err := r.db.Create(u).Error; err != nil {
		return "", err
	}
	return strconv.Itoa(int(u.ID)), nil
}

func (r *gormUserRepository) FindByEmailAndPass(email string, pass string) (*user.User, error) {
	var u user.User
	err := r.db.Where("email = ? AND password = ?", email, pass).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
