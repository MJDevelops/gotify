package envs

import (
	"errors"
	"os"
)

type Env struct {
	ClientSecret string
	ClientID     string
}

func LoadEnv() (Env, error) {
	clientSecret, clientId := os.Getenv("GOTIFY_CLIENT_SECRET"), os.Getenv("GOTIFY_CLIENT_ID")
	var err error

	if clientSecret == "" || clientId == "" {
		err = errors.New("client id and/or client secret are/is not set, exiting")
	}
	
	return Env{
		ClientSecret: clientSecret,
		ClientID: clientId,
	}, err
}
