package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/zmb3/spotify/v2"

	"music-service/internal/infrastructure/spotify_client"
)

// formatDuration converts a time.Duration to MM:SS format
func formatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// TODO: Search music, Get Random Songs, Get Top Songs, Get Recommendations
type Tracks struct {
	ID       spotify.ID `json:"id"`
	Name     string     `json:"name"`
	Artists  []string   `json:"artists"`
	ImageURL string     `json:"image_url"`
	Genres   []string   `json:"genres"`
	Duration string     `json:"duration"`
}

func GetRandomTracks(ctx context.Context) ([]Tracks, error) {
	client, err := spotify_client.NewSpotifyClient()
	if err != nil {
		return nil, err
	}
	var randomSearchQueries = []string{
		"year:2020-2024", "year:2015-2019", "year:2010-2014",
		"genre:pop", "genre:rock", "genre:electronic", "genre:indie",
		"genre:hip-hop", "genre:r-n-b", "genre:jazz", "genre:alternative",
		"a", "e", "i", "o", "the", "love", "night", "dream", "sun",
	}

	query := randomSearchQueries[rand.Intn(len(randomSearchQueries))]
	result, err := client.Search(
		ctx,
		query,
		spotify.SearchTypeTrack,
	)
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}
	var tracks []Tracks

	for _, t := range result.Tracks.Tracks {
		// Get all artists on the track
		var artistNames []string
		for _, artist := range t.Artists {
			artistNames = append(artistNames, artist.Name)
		}

		// Get all genres on the track
		genres, err := GetGenresForTrack(ctx, client, t.ID)
		if err != nil {
			log.Fatalf("failed to get genres: %v", err)
		}

		tracks = append(tracks, Tracks{
			ID:       t.ID,
			Name:     t.Name,
			Artists:  artistNames,
			ImageURL: t.Album.Images[0].URL,
			Genres:   genres,
			Duration: formatDuration(t.TimeDuration()),
		})
	}
	return tracks, nil
}
