package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/service"
)

// map для хранения в памяти
type URL struct {
	OrginalURL string `json:"url"`
}

var urlMap = make(map[string]models.URL)

var req URL

// Оюработчик POST
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid Json", http.StatusBadRequest)
			return
		}
		if !service.IsValidURL(req.OrginalURL) {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		shortCode, err := service.GenerateRandomShortString(6)
		if err != nil {
			http.Error(w, "Error generating short URL", http.StatusInternalServerError)
			return
		}
		urlMap[shortCode] = models.URL{
			OriginalURL: req.OrginalURL,
			ShortCode:   shortCode,
			CreatedAt:   time.Now(),
			ExpiresAt:   time.Now().Add(24 * time.Hour), // Ссылка истекает через 24 часа
			ClickCount:  0,
		}

		// Формируем ответ
		response := map[string]string{
			"short_url": "http://localhost:8080/" + shortCode,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
