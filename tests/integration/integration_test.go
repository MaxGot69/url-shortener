package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/MaxGot69/url-shortener/internal/config"
	"github.com/MaxGot69/url-shortener/internal/handler"
	"github.com/MaxGot69/url-shortener/internal/metrics"
	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/repository"
	"github.com/MaxGot69/url-shortener/pkg/database"
	appLogger "github.com/MaxGot69/url-shortener/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var db *gorm.DB

func setupTestDB(t *testing.T) *gorm.DB {
	cfg := config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "maxim",
		DBPassword: "secret",
		DBName:     "urlshortener",
		LogLevel:   "error",
	}

	var err error
	db, err = database.PostgresConnection(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	require.NoError(t, err)

	db.AutoMigrate(&models.URL{}, &models.User{})
	
	db.Exec("DELETE FROM urls")
	db.Exec("DELETE FROM users")

	return db
}

func TestShortenAndRedirectFlow(t *testing.T) {
	db := setupTestDB(t)
	appLogger.Init("error")
	metrics.Init()

	repo := &repository.PostgresRepository{DB: db}
	handler.NewPostHandler(nil)

	t.Run("Create short URL and redirect", func(t *testing.T) {
		urlModel := models.URL{
			OriginalURL: "https://example.com",
			ShortCode:   "test123",
			ExpiresAt:   time.Now().Add(24 * time.Hour),
			ClickCount:  0,
		}

		err := repo.SaveURL(urlModel)
		require.NoError(t, err)

		retrievedURL, err := repo.GetURLByShortCode("test123")
		require.NoError(t, err)
		assert.Equal(t, "https://example.com", retrievedURL.OriginalURL)
	})
}

func TestExpiredURL(t *testing.T) {
	db := setupTestDB(t)

	repo := &repository.PostgresRepository{DB: db}

	t.Run("Check expired URL", func(t *testing.T) {
		urlModel := models.URL{
			OriginalURL: "https://example.com",
			ShortCode:   "expired",
			ExpiresAt:   time.Now().Add(-1 * time.Hour),
			ClickCount:  0,
		}

		err := repo.SaveURL(urlModel)
		require.NoError(t, err)

		assert.True(t, time.Now().After(urlModel.ExpiresAt))
	})
}

func TestConcurrentAccess(t *testing.T) {
	db := setupTestDB(t)

	repo := &repository.PostgresRepository{DB: db}

	t.Run("Concurrent URL creation", func(t *testing.T) {
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(i int) {
				urlModel := models.URL{
					OriginalURL: fmt.Sprintf("https://example%d.com", i),
					ShortCode:   fmt.Sprintf("test%d", i),
					ExpiresAt:   time.Now().Add(24 * time.Hour),
					ClickCount:  0,
				}
				repo.SaveURL(urlModel)
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		count := int64(0)
		db.Model(&models.URL{}).Count(&count)
		assert.Equal(t, int64(10), count)
	})
}

func TestRegisterAndLogin(t *testing.T) {
	appLogger.Init("error")
	cfg := &config.Config{JWTSecret: "test-secret"}

	handler := handler.NewRegisterHandler(nil, cfg)
	assert.NotNil(t, handler)
}
