package service

import (
	"context"
	"log"
	"math/rand"

	"github.com/zmb3/spotify/v2"

	"music-service/internal/infrastructure/spotify_client"
)

// TODO: Search music, Get Random Songs, Get Top Songs, Get Recommendations
type Tracks struct {
	ID         spotify.ID `json:"id"`
	Name       string     `json:"name"`
	Artists    []string   `json:"artists"`
	TrackImage string     `json:"track_image"`
	Genres     []string   `json:"genres"`
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
			ID:         t.ID,
			Name:       t.Name,
			Artists:    artistNames,
			TrackImage: t.Album.Images[0].URL,
			Genres:     genres,
		})
	}
	return tracks, nil
}
