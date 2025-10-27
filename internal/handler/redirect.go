package handler

import (
	"net/http"
	"time"

	"github.com/MaxGot69/url-shortener/internal/metrics"
	"github.com/MaxGot69/url-shortener/internal/repository"
	"github.com/MaxGot69/url-shortener/pkg/cache"
	appLogger "github.com/MaxGot69/url-shortener/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type RedirectHandler struct {
	repo        repository.UserRepository
	redisClient *cache.RedisClient
}

func NewRedirectHandler(repo repository.UserRepository) *RedirectHandler {
	redisClient, err := cache.NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		appLogger.Logger.Warn("Failed to connect to Redis", "error", err)
	}
	return &RedirectHandler{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	metrics.URLRedirectsTotal.Inc()
	shortCode := chi.URLParam(r, "shortCode")

	appLogger.Logger.Info("Redirect request", "short_code", shortCode)

	if h.redisClient != nil {
		cacheURL, err := h.redisClient.Get(shortCode)
		if err == nil && cacheURL != "" {
			appLogger.Logger.Info("Redis cache hit", "short_code", shortCode)
			http.Redirect(w, r, cacheURL, http.StatusFound)
			return
		}
	}

	urlRepo, ok := h.repo.(repository.URLRepository)
	if !ok {
		appLogger.Logger.Error("Repository does not implement URLRepository")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	url, err := urlRepo.GetURLByShortCode(shortCode)
	if err != nil {
		appLogger.Logger.Warn("URL not found", "short_code", shortCode, "error", err)
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	if time.Now().After(url.ExpiresAt) {
		appLogger.Logger.Info("URL expired", "short_code", shortCode)
		http.Error(w, "Short URL has expired", http.StatusGone)
		return
	}

	if h.redisClient != nil {
		err = h.redisClient.Set(shortCode, url.OriginalURL, 24*time.Hour)
		if err != nil {
			appLogger.Logger.Warn("Failed to cache in Redis", "error", err)
		}
	}

	url.ClickCount++
	err = urlRepo.UpdateClickCount(*url)
	if err != nil {
		appLogger.Logger.Error("Failed to update click count", "error", err)
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
	appLogger.Logger.Info("Redirect successful", "short_code", shortCode, "url", url.OriginalURL, "clicks", url.ClickCount)
}