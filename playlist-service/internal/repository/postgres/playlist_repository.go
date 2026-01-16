package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	id                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	email               string    `gorm:"unique;not null"`
	username            string    `gorm:"unique;not null"`
	password            string    `gorm:"not null"`
	display_name        string
	bio                 string
	profile_picture_url string
	spotify_user_id     string
	created_at          time.Time
	updated_at          time.Time
	last_active_at      time.Time
}
