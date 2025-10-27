package models

import (
	"time"
)

type URL struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OriginalURL string    `gorm:"not null" json:"original_url"`
	ShortCode   string    `gorm:"uniqueIndex;not null" json:"short_code"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt   time.Time `gorm:"not null" json:"expires_at"`
	ClickCount  int       `gorm:"default:0" json:"click_count"`
}
