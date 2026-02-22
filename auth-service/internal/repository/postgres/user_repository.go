package postgres

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                string `gorm:"primaryKey"`
	Email             string `gorm:"uniqueIndex;not null"`
	Username          string `gorm:"uniqueIndex;not null"`
	PasswordHash      string `gorm:"not null"`
	DisplayName       string
	Bio               string
	ProfilePictureURL string
	SpotifyUserID     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RefreshToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	TokenHash string
	ExpiresAt time.Time `gorm:"index"`
} // ยังไม่ได้ใช้
