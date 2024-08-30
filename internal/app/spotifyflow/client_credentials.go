package spotifyflow

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/MJDevelops/gotify/internal/pkg/envs"
	"github.com/MJDevelops/gotify/internal/pkg/jsons"
)

type SpotifyClientCredential struct {
	AccessToken string
}

func (s *SpotifyClientCredential) RequestAccessToken() error {
	var err error
	jsonMap := make(map[string]string)
	envs, err := envs.LoadEnv()

	if err != nil {
		log.Println("Couldn't load envs")
		return err
	}

	req, err := buildRequest(envs)

	if err != nil {
		log.Println("Couldn't build request")
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Println("Couldn't 'POST' request")
		return err
	} else if res.StatusCode != 200 {
		err = fmt.Errorf("response contains a status code of %v", res.StatusCode)
		return err
	}

	resBytes, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("Error trying to read response body")
		return err
	}

	if err = jsons.ParseJSON(resBytes, &jsonMap); err != nil {
		log.Println("Error trying to parse JSON")
		return err
	}

	s.AccessToken = jsonMap["access_token"]
	return nil
}

func buildRequest(e *envs.GotifyEnv) (*http.Request, error) {
	data := &url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", e.GotifyClientID)
	data.Add("client_secret", e.GotifyClientSecret)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}
