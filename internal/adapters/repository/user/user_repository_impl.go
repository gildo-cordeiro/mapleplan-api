package user

import (
	"context"
	"errors"
	"strings"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/user"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/repositories"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repositories.UserRepository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return r.db
	}
	if tx, ok := ctx.Value("tx_key").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *RepositoryImpl) Save(user *user.User) (string, error) {
	if err := r.db.Create(user).Error; err != nil {
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
	return user.ID, nil
}

func (r *RepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.getDB(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *RepositoryImpl) FindByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.getDB(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *RepositoryImpl) Update(ctx context.Context, id string, u *user.User) error {
	if err := r.getDB(ctx).Where("id = ?", id).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *RepositoryImpl) SearchByName(userId string, name string) ([]*user.User, error) {
	var users []user.User
	pattern := "%" + strings.ToLower(name) + "%"
	err := r.db.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? AND id != ?", pattern, pattern, userId).Find(&users).Error
	if err != nil {
		return nil, err
	}

	var result []*user.User
	for _, u := range users {
		result = append(result, &u)
	}
	return result, nil
}
