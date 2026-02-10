package track

import (
	"time"

	"github.com/zmb3/spotify/v2"
)

// --- API Response Structures (DTOs) ---

// TrackResponse is the clean JSON we return to the frontend
type TrackResponse struct {
	ID       spotify.ID `json:"id"`
	Name     string     `json:"name"`
	Artists  []string   `json:"artists"`
	ImageURL string     `json:"image_url"`
	Genres   []string   `json:"genres"`
	Duration string     `json:"duration"`
}

// --- Database Entities (GORM) ---

type Track struct {
	ID          string    `gorm:"primaryKey"`
	Name        string    `gorm:"not null;index"`
	ImageURL    string    // Album cover (640x640)
	Artists     []Artist  `gorm:"many2many:track_artists;"`
	Genres      []Genre   `gorm:"many2many:track_genres;"`
	ReleaseDate time.Time `gorm:"index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Artist struct {
	ID         string    `gorm:"primaryKey"`
	Name       string    `gorm:"not null;index"`
	SpotifyURL string    `gorm:"not null"`
	Tracks     []Track   `gorm:"many2many:track_artists;"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

type Genre struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	Tracks      []Track   `gorm:"many2many:track_genres;"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}