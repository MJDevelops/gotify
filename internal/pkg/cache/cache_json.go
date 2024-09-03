package cache

import (
	"encoding/json"
	"os"

	"github.com/MJDevelops/gotify/internal/app/spotifyflow"
	"github.com/MJDevelops/gotify/internal/pkg/logs"
)

var authCacheFile = "gotify_cache.json"
var logger = logs.GetLoggerInstance()

func CacheSpotifyAuthCode(s *spotifyflow.SpotifyAuthorizationCode) error {
	file, _ := os.OpenFile(authCacheFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	data, err := json.Marshal(s)

	if err != nil {
		logger.Println("Couldn't cache auth data")
		return err
	}

	if _, err = file.Write(data); err != nil {
		logger.Println("Couldn't write to cache file")
		return err
	}

	file.Close()
	return nil
}
