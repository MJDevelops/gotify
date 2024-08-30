package envs

import (
	"errors"
	"os"
)

type GotifyEnv struct {
	GotifyClientSecret string
	GotifyClientID     string
}

func LoadEnv() (GotifyEnv, error) {
	clientSecret, clientId := os.Getenv("GOTIFY_CLIENT_SECRET"), os.Getenv("GOTIFY_CLIENT_ID")
	var err error

	if clientSecret == "" || clientId == "" {
		err = errors.New("client id and/or client secret are/is not set, exiting")
	}

	return GotifyEnv{
		GotifyClientSecret: clientSecret,
		GotifyClientID:     clientId,
	}, err
}
