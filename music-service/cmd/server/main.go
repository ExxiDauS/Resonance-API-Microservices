package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify/v2"
)

type TrackResponse struct {
	Title       string   `json:"title"`
	Artist      string   `json:"artist"`
	Image       string   `json:"image"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

func main() {

	r := gin.Default()
	r.GET("/tracks/random", getRandomTracks)

	log.Println("Server running on :8080")
	r.Run(":8080")

}

func getRandomTracks(c *gin.Context) {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     "1956424920024e97bd4c6847d5af05dc",
		ClientSecret: "9e82d5748fa344779609117824584aa0",
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	var randomSearchQueries = []string{
		"year:2020-2024", "year:2015-2019", "year:2010-2014",
		"genre:pop", "genre:rock", "genre:electronic", "genre:indie",
		"genre:hip-hop", "genre:r-n-b", "genre:jazz", "genre:alternative",
		"a", "e", "i", "o", "the", "love", "night", "dream", "sun",
	}

	query := randomSearchQueries[rand.Intn(len(randomSearchQueries))]

	offset := rand.Intn(100)

	results, err := client.Search(
		ctx,
		query,
		spotify.SearchTypeTrack,
		spotify.Limit(10),
		spotify.Offset(offset),
	)

	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	log.Printf("Query: %s", query)

	var tracks []TrackResponse

	for _, t := range results.Tracks.Tracks {
		image := ""
		if len(t.Album.Images) > 0 {
			image = t.Album.Images[0].URL
		}

		tracks = append(tracks, TrackResponse{
			Title:       t.Name,
			Artist:      t.Artists[0].Name,
			Image:       image,
			Tags:        []string{"Spotify", query},
			Description: "Discovered via Spotify",
		})
	}

	c.JSON(http.StatusOK, tracks)
}
