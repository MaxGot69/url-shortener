package models

import (
	"time"
)

type URL struct {
	ID          string    `json:"ID"`
	OriginalURL string    `json:"originalUrl"`
	ShortCode   string    `json:"shortCode"`
	CreatedAt   time.Time // ДАТА СОЗДАНИЯ
	ExpiresAt   time.Time // ДАТА ИСТЕЧЕНИЯ СРОКА ССЫЛКИ
	ClickCount  int       // СЧЕТЧИК КЛИКОВ
}
