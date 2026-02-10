package track

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/zmb3/spotify/v2"
)

type Service struct {
	repo          Repository
	spotifyClient *spotify.Client
}

// GetTrack implements [playlist.MusicProvider].
func (s *Service) GetTrack(id string) (*Track, error) {
	panic("unimplemented")
}

func NewService(repo Repository, spotifyClient *spotify.Client) *Service {
	return &Service{
		repo:          repo,
		spotifyClient: spotifyClient,
	}
}

// Update the Interface or Method Signature
func (s *Service) GetSuggestions(ctx context.Context, source string) ([]TrackResponse, error) {
	if source == "db" {
		return s.getRandomFromDB(ctx)
	}
	// Default to Spotify (True Random)
	return s.getRandomFromSpotify(ctx)
}

func (s *Service) getRandomFromDB(ctx context.Context) ([]TrackResponse, error) {
	tracks, err := s.repo.GetRandomTracks(10) // Fetch 10 random
	if err != nil {
		return nil, err
	}

	// Map DB entity to Response
	var response []TrackResponse
	for _, t := range tracks {
		// Convert Artists/Genres to simple strings for response
		var artists []string
		for _, a := range t.Artists {
			artists = append(artists, a.Name)
		}
		var genres []string
		for _, g := range t.Genres {
			genres = append(genres, g.Name)
		}

		response = append(response, TrackResponse{
			ID:       spotify.ID(t.ID),
			Name:     t.Name,
			Artists:  artists,
			ImageURL: t.ImageURL,
			Genres:   genres,
			Duration: "0:00", // DB might not store duration, handle accordingly
		})
	}
	return response, nil
}

func (s *Service) getRandomFromSpotify(ctx context.Context) ([]TrackResponse, error) {
	// 1. Pick a random query
	queries := []string{
		"year:2020-2024", "year:2015-2019", "year:2010-2014",
		"genre:pop", "genre:rock", "genre:electronic", "genre:indie",
		"genre:hip-hop", "genre:r-n-b", "genre:jazz", "genre:alternative",
	}
	query := queries[rand.Intn(len(queries))]

	// 2. Call Spotify
	result, err := s.spotifyClient.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		return nil, fmt.Errorf("spotify search failed: %w", err)
	}

	var tracks []TrackResponse

	// 3. Process & Save
	for _, t := range result.Tracks.Tracks {
		// --- A. Data Preparation ---

		// Get detailed Genres
		genreNames, err := s.getGenresForTrack(ctx, t)
		if err != nil {
			log.Printf("failed to get genres for track %s: %v", t.ID, err)
			genreNames = []string{}
		}

		// Prepare Artist objects for DB
		var dbArtists []Artist
		var artistNames []string
		for _, a := range t.Artists {
			artistNames = append(artistNames, a.Name)
			dbArtists = append(dbArtists, Artist{
				ID:         string(a.ID),
				Name:       a.Name,
				SpotifyURL: "https://open.spotify.com/artist/" + string(a.ID),
			})
		}

		// Prepare Genre objects for DB
		var dbGenres []Genre
		for _, name := range genreNames {
			dbGenres = append(dbGenres, Genre{Name: name})
		}

		// Parse Release Date (Spotify returns "YYYY" or "YYYY-MM" or "YYYY-MM-DD")
		releaseDate := parseSpotifyDate(t.Album.ReleaseDate)

		// --- B. Map to Database Entity ---
		dbTrack := &Track{
			ID:          string(t.ID),
			Name:        t.Name,
			ImageURL:    t.Album.Images[0].URL,
			ReleaseDate: releaseDate,
			Artists:     dbArtists,
			Genres:      dbGenres,
		}

		// --- C. Save to Database ---
		// We run this in a goroutine or just sequentially.
		// Sequentially is safer for now to ensure data integrity.
		if err := s.repo.SaveTrack(dbTrack); err != nil {
			log.Printf("Failed to save track %s: %v", t.Name, err)
			// We continue even if save fails, to still return the result to user
		}

		// --- D. Map to API Response ---
		tracks = append(tracks, TrackResponse{
			ID:       t.ID,
			Name:     t.Name,
			Artists:  artistNames,
			ImageURL: t.Album.Images[0].URL,
			Genres:   genreNames,
			Duration: s.formatDuration(t.TimeDuration()),
		})
	}

	return tracks, nil
}

// --- Helpers ---

func parseSpotifyDate(dateStr string) time.Time {
	// Spotify release dates vary in precision
	layouts := []string{"2006-01-02", "2006-01", "2006"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t
		}
	}
	return time.Now() // Fallback
}

func (s *Service) getGenresForTrack(ctx context.Context, track spotify.FullTrack) ([]string, error) {
	genreMap := make(map[string]bool)
	for _, simpleArtist := range track.Artists {
		fullArtist, err := s.spotifyClient.GetArtist(ctx, simpleArtist.ID)
		if err != nil {
			continue
		}
		for _, g := range fullArtist.Genres {
			genreMap[g] = true
		}
	}
	var genres []string
	for g := range genreMap {
		genres = append(genres, g)
	}
	return genres, nil
}

func (s *Service) formatDuration(d time.Duration) string {
	return fmt.Sprintf("%d:%02d", int(d.Seconds())/60, int(d.Seconds())%60)
}
