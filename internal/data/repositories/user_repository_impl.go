package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/user"
	ports "github.com/gildo-cordeiro/mapleplan-api/internal/ports/repositories"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) ports.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return r.db
	}
	if tx, ok := ctx.Value("tx_key").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *UserRepositoryImpl) Save(user *user.User) (string, error) {
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

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.getDB(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.getDB(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, id string, u *user.User) error {
	if err := r.getDB(ctx).Where("id = ?", id).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) SearchByName(userId string, name string) ([]*user.User, error) {
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
