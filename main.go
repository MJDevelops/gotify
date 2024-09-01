package main

import (
	"log"

	"github.com/MJDevelops/gotify/internal/app/spotifyflow"
)

func main() {
	authorizationCode := &spotifyflow.SpotifyAuthorizationCode{}
	err := authorizationCode.Authorize()
	if err != nil {
		log.Fatalf("%v", err)
	}
}
