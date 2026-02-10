package playlist

import (
	"context"
	"fmt"
	"music-service/internal/track" // Use music types

	"github.com/zmb3/spotify/v2"
)

// MusicProvider lets us talk to the Music Service
type MusicProvider interface {
	GetTrack(id string) (*track.Track, error)
}

type Service struct {
	repo          Repository
	musicProvider MusicProvider
}

func NewService(repo Repository, mp MusicProvider) *Service {
	return &Service{
		repo:          repo,
		musicProvider: mp,
	}
}

func (s *Service) CreatePlaylist(name, userID string) (*PlaylistResponse, error) {
	p := &Playlist{
		Name:   name,
		UserID: userID,
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return &PlaylistResponse{ID: p.ID, Name: p.Name, CreatedAt: p.CreatedAt}, nil
}

func (s *Service) AddTrack(ctx context.Context, playlistID, trackID string) error {
	// 1. Verify Track Exists via Music Provider
    // The Music Service's GetTrackByID method (exposed via interface) 
    // should handle fetching from DB.
	track, err := s.musicProvider.GetTrack(trackID)
	if err != nil {
		return fmt.Errorf("track not found: %v", err)
	}

	// 2. Add to Playlist
	return s.repo.AddTrack(playlistID, track)
}

func (s *Service) GetPlaylist(id string) (*PlaylistResponse, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Map to Response
	var trackResponses []track.TrackResponse
	for _, t := range p.Tracks {
		// Basic mapping (simplify as needed)
        var artists []string
        for _, a := range t.Artists { artists = append(artists, a.Name) }
        
		trackResponses = append(trackResponses, track.TrackResponse{
			ID:       spotify.ID(t.ID),
			Name:     t.Name,
			Artists:  artists,
			ImageURL: t.ImageURL,
		})
	}

	return &PlaylistResponse{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		Tracks:    trackResponses,
	}, nil
}