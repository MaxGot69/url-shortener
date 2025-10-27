package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MaxGot69/url-shortener/internal/metrics"
	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/service"
	"github.com/MaxGot69/url-shortener/pkg/cache"
	appLogger "github.com/MaxGot69/url-shortener/pkg/logger"
)

type URLRequest struct {
	OriginalURL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type PostHandler struct {
	urlService *service.URLService
	redisClient *cache.RedisClient
}

func NewPostHandler(urlService *service.URLService) *PostHandler {
	redisClient, err := cache.NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		appLogger.Logger.Warn("Failed to connect to Redis", "error", err)
	}
	return &PostHandler{
		urlService:  urlService,
		redisClient: redisClient,
	}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics.URLShortensTotal.Inc()

	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appLogger.Logger.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !h.urlService.IsValidURL(req.OriginalURL) {
		appLogger.Logger.Warn("Invalid URL format", "url", req.OriginalURL)
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	shortCode, err := h.urlService.GenerateRandomShortString(6)
	if err != nil {
		appLogger.Logger.Error("Error generating short code", "error", err)
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	urlModel := models.URL{
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		ClickCount:  0,
	}

	err = h.urlService.SaveURL(urlModel)
	if err != nil {
		appLogger.Logger.Error("Failed to save URL", "error", err)
		http.Error(w, "Error saving URL", http.StatusInternalServerError)
		return
	}

	if h.redisClient != nil {
		err = h.redisClient.Set(shortCode, req.OriginalURL, 24*time.Hour)
		if err != nil {
			appLogger.Logger.Warn("Failed to cache in Redis", "error", err)
		}
	}

	response := ShortenResponse{
		ShortURL: "http://localhost:8081/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		appLogger.Logger.Error("Failed to encode response", "error", err)
	}

	appLogger.Logger.Info("URL shortened", "short_code", shortCode, "original_url", req.OriginalURL)
}