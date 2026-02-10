package playlist

import (
	"music-service/internal/track" // Import music to link entities
	"time"
)

// --- API DTOs ---

type CreatePlaylistRequest struct {
	Name   string `json:"name" binding:"required"`
	UserID string `json:"user_id" binding:"required"` // In real app, get this from JWT
}

type AddTrackRequest struct {
	TrackID string `json:"track_id" binding:"required"`
}

type PlaylistResponse struct {
	ID        uint                `json:"id"`
	Name      string              `json:"name"`
	CreatedAt time.Time           `json:"created_at"`
	Tracks    []track.TrackResponse `json:"tracks,omitempty"`
}

// --- DB Entity ---

type Playlist struct {
	ID        uint          `gorm:"primaryKey"`
	Name      string        `gorm:"not null"`
	UserID    string        `gorm:"index"`
	Tracks    []track.Track `gorm:"many2many:playlist_tracks;"` // GORM handles the join table
	CreatedAt time.Time
	UpdatedAt time.Time
}