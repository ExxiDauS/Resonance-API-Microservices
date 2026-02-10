package track

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetTrackByID(id string) (*Track, error)
	GetRandomTracks(limit int) ([]Track, error)
	SaveTrack(track *Track) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetTrackByID(id string) (*Track, error) {
	var track Track
	err := r.db.Preload("Artists").Preload("Genres").First(&track, "id = ?", id).Error
	return &track, err
}

// SaveTrack handles "Upserting" a track and its relationships
func (r *repository) SaveTrack(track *Track) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Handle Genres: Ensure they exist, prevent duplicates
		var cleanGenres []Genre
		for _, g := range track.Genres {
			var dbGenre Genre
			// Find by Name, or Create if missing
			if err := tx.Where(Genre{Name: g.Name}).FirstOrCreate(&dbGenre).Error; err != nil {
				return err
			}
			cleanGenres = append(cleanGenres, dbGenre)
		}
		track.Genres = cleanGenres

		// 2. Handle Artists: Ensure they exist
		var cleanArtists []Artist
		for _, a := range track.Artists {
			var dbArtist Artist
			// Find by ID, or Create if missing (Update attributes if found)
			if err := tx.Assign(Artist{Name: a.Name, SpotifyURL: a.SpotifyURL}).FirstOrCreate(&dbArtist, Artist{ID: a.ID}).Error; err != nil {
				return err
			}
			cleanArtists = append(cleanArtists, dbArtist)
		}
		track.Artists = cleanArtists

		// 3. Save Track (Upsert)
		// OnConflict tells Postgres to update columns if the ID already exists
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(track).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) GetRandomTracks(limit int) ([]Track, error) {
	var tracks []Track
	// PostgreSQL specific: ORDER BY RANDOM()
	err := r.db.Preload("Artists").Preload("Genres").
		Order("RANDOM()").
		Limit(limit).
		Find(&tracks).Error
	return tracks, err
}