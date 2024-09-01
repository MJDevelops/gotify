package spotifyflow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/MJDevelops/gotify/internal/pkg/envs"
)

type SpotifyClientCredential struct {
	AccessToken string `json:"access_token"`
}

func (s *SpotifyClientCredential) Authorize() error {
	envs := envs.LoadEnv()

	req, err := buildRequest(envs)

	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build request")
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't 'POST' request")
		return err
	} else if sc := res.StatusCode; sc != 200 {
		err = fmt.Errorf("response contains a status code of %v", sc)
		return err
	}

	resBytes, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error trying to read response body")
		return err
	}

	if err = json.Unmarshal(resBytes, s); err != nil {
		fmt.Fprint(os.Stderr, "Error trying to parse JSON")
		return err
	}

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
