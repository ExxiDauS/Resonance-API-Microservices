package service

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
)

func GetGenresForTrack(
	ctx context.Context,
	client *spotify.Client,
	trackID spotify.ID,
) ([]string, error) {
	// Get track details
	track, err := client.GetTrack(ctx, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	// Get all unique genres from all artists on the track
	genreMap := make(map[string]bool)

	for _, artist := range track.Artists {
		// Get full artist details (simple artist doesn't include genres)
		fullArtist, err := client.GetArtist(ctx, artist.ID)
		if err != nil {
			log.Printf("failed to get artist %s: %v", artist.ID, err)
			continue
		}

		// Add genres to map
		for _, genre := range fullArtist.Genres {
			genreMap[genre] = true
		}
	}

	// Convert map to slice
	genres := make([]string, 0, len(genreMap))
	for genre := range genreMap {
		genres = append(genres, genre)
	}

	return genres, nil
}
