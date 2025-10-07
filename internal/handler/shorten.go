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

// Оюработчик POST
func PostHandler(service *service.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var req URL

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
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
				ExpiresAt:   time.Now().Add(24 * time.Hour),
				ClickCount:  0,
			}

			response := map[string]string{
				"short_url": "http://localhost:8081/" + shortCode,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
