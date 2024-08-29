package gotify

import (
	"fmt"
	"log"

	"github.com/MJDevelops/gotify/internal/pkg/envs"
)

func StartClient() {
	log.Println("Loading envs")
	env, err := envs.LoadEnv()

	if err != nil {
		log.Fatal("Couldn't load client secret and/or client id, exiting")
	}

	log.Println("Envs loaded")
	fmt.Println(env.GotifyClientID, env.GotitifyClientSecret)
}
