package models

import (
	"time"
)

type User struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
