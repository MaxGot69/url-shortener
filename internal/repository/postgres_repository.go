package repository

import (
	"github.com/MaxGot69/url-shortener/internal/models"
	"gorm.io/gorm"
)

type URLRepository interface {
	SaveURL(url models.URL) error
	GetURLByShortCode(shortCode string) (*models.URL, error)
	UpdateClickCount(url models.URL) error
}

type PostgresRepository struct {
	DB *gorm.DB
}

func (r *PostgresRepository) SaveURL(url models.URL) error {
	return r.DB.Create(&url).Error
}

func (r *PostgresRepository) GetURLByShortCode(shortCode string) (*models.URL, error) {
	var url models.URL
	if err := r.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *PostgresRepository) UpdateClickCount(url models.URL) error {
	return r.DB.Model(&url).Update("click_count", url.ClickCount).Error
}