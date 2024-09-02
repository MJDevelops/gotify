package envs

import (
	"log"
	"os"
)

type GotifyEnv struct {
	GotifyClientID string
}

func LoadEnv() *GotifyEnv {
	clientId := os.Getenv("GOTIFY_CLIENT_ID")

	if clientId == "" {
		log.Fatal("client_id not set")
	}

	return &GotifyEnv{
		GotifyClientID: clientId,
	}
}
