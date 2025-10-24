package repository

import (
	"errors"

	"github.com/MaxGot69/url-shortener/internal/models"
	"gorm.io/gorm"
)

type UserReposytory interface {
	FindByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

func (r *PostgresRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *PostgresRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}
