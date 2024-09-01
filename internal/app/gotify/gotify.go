package gotify

import (
	"fmt"
	"log"

	"github.com/MJDevelops/gotify/internal/pkg/envs"
)

func StartClient() {
	log.Println("Loading envs")
	env := envs.LoadEnv()

	log.Println("Envs loaded")
	fmt.Println(env.GotifyClientID, env.GotifyClientSecret)
}
