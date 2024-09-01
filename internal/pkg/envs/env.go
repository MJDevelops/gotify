package envs

import (
	"log"
	"os"
)

type GotifyEnv struct {
	GotifyClientSecret string
	GotifyClientID     string
}

func LoadEnv() *GotifyEnv {
	clientSecret, clientId := os.Getenv("GOTIFY_CLIENT_SECRET"), os.Getenv("GOTIFY_CLIENT_ID")

	if clientSecret == "" || clientId == "" {
		log.Fatal("client_secret and/or client_id not set")
	}

	return &GotifyEnv{
		GotifyClientSecret: clientSecret,
		GotifyClientID:     clientId,
	}
}
