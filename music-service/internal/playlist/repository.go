package playlist

import (
	"music-service/internal/track"

	"gorm.io/gorm"
)

type Repository interface {
	Create(playlist *Playlist) error
	GetByID(id string) (*Playlist, error)
	AddTrack(playlistID string, track *track.Track) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(playlist *Playlist) error {
	return r.db.Create(playlist).Error
}

func (r *repository) GetByID(id string) (*Playlist, error) {
	var playlist Playlist
	// Preload Tracks and their Artists/Genres for full details
	err := r.db.Preload("Tracks.Artists").
		Preload("Tracks.Genres").
		First(&playlist, "id = ?", id).Error
	return &playlist, err
}

func (r *repository) AddTrack(playlistID string, track *track.Track) error {
	// Association Mode: Append track to playlist
	// GORM automatically handles the "playlist_tracks" table
	var playlist Playlist
	if err := r.db.First(&playlist, "id = ?", playlistID).Error; err != nil {
		return err
	}
	
	return r.db.Model(&playlist).Association("Tracks").Append(track)
}