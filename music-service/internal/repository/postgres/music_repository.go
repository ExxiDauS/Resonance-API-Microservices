package postgres

import (
	"time"

	"gorm.io/gorm"
)

type Track struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null;index"`
	DurationMs  int    `gorm:"not null"`
	PreviewURL  *string
	SpotifyURL  string         `gorm:"not null"`
	Popularity  int            `gorm:"index"`
	ImageURL    string         // Album cover (640x640)
	AlbumID     string         `gorm:"index"`
	Album       Album          `gorm:"foreignKey:AlbumID"`
	Artists     []Artist       `gorm:"many2many:track_artists;"`
	Genres      []Genre        `gorm:"many2many:track_genres;"`
	IsExplicit  bool           `gorm:"default:false"`
	ReleaseDate time.Time      `gorm:"index"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Artist stores artist information
type Artist struct {
	ID         string         `gorm:"primaryKey"`
	Name       string         `gorm:"not null;index"`
	SpotifyURL string         `gorm:"not null"`
	Tracks     []Track        `gorm:"many2many:track_artists;"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Album stores album information
type Album struct {
	ID          string    `gorm:"primaryKey"`
	Name        string    `gorm:"not null;index"`
	AlbumType   string    `gorm:"type:varchar(20)"` // album, single, compilation
	ImageURL    string    // 640x640 cover art
	SpotifyURL  string    `gorm:"not null"`
	ReleaseDate time.Time `gorm:"index"`
	TotalTracks int
	Tracks      []Track        `gorm:"foreignKey:AlbumID"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Genre stores music genres/tags (e.g., Electronic, Dream Pop, Indie)
type Genre struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"uniqueIndex;not null"` // Electronic, Dream Pop, etc.
	Slug        string         `gorm:"uniqueIndex;not null"` // electronic, dream-pop
	Description string         `gorm:"type:text"`
	ImageURL    string         // Genre artwork for UI
	Tracks      []Track        `gorm:"many2many:track_genres;"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
