package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/MaxGot69/url-shortener/pkg/cache"
	"github.com/go-chi/chi/v5"
)

var RedisClient *cache.RedisClient

func init() {
	var err error
	RedisClient, err = cache.NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		log.Printf("Redis connection failed: %v", err)
	} else {
		log.Println("Redis client initialized")
	}
}

// Get обработчик
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode") // Получаем короткую из url
	// Достаем url из urlMap
	url, ok := urlMap[shortCode]
	if !ok {
		http.Error(w, "Short URL not found", http.StatusFound) // если нету то 404
	}
	// Проверяем не истек ли срок действия ссылки
	if time.Now().After(url.ExpiresAt) {
		http.Error(w, "Short URL has expired", http.StatusGone) // 410
		// Увеличиваем счетчик кликов
		url.ClickCount++
		urlMap[shortCode] = url
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound) // Перенаправлем на оригинальную ссылку
	log.Printf("Redirected %s to %s", url.ShortCode, url.OriginalURL)

}
