package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/MJDevelops/gotify/internal/app/spotifyflow"
	"github.com/MJDevelops/gotify/internal/pkg/logs"
)

type playlistRequest struct {
	*spotifyflow.SpotifyAuthorizationCode
}

type simplifiedPlaylistObject struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalUrls  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href   string `json:"href"`
	Id     string `json:"id"`
	Images []struct {
		Url    string      `json:"url"`
		Height json.Number `json:"height"`
		Width  json.Number `json:"width"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Href  string `json:"href"`
			Total int    `json:"total"`
		} `json:"followers"`
		Href        string `json:"href"`
		Id          string `json:"id"`
		Type        string `json:"type"`
		Uri         string `json:"uri"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	Public     bool   `json:"public"`
	SnapshotId string `json:"snapshot_id"`
	Tracks     []struct {
		Href  string      `json:"href"`
		Total json.Number `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	Uri  string `json:"uri"`
}

type currentUserPlaylistResponse struct {
	Href     string                     `json:"href"`
	Limit    json.Number                `json:"limit"`
	Next     string                     `json:"next"`
	Offset   json.Number                `json:"offset"`
	Previous string                     `json:"previous"`
	Total    json.Number                `json:"total"`
	Items    []simplifiedPlaylistObject `json:"items"`
}

const apiURL string = "https://api.spotify.com/v1"

var logger = logs.GetLoggerInstance()

func InitPlaylistRequest(s *spotifyflow.SpotifyAuthorizationCode) (*playlistRequest, error) {
	return &playlistRequest{
		s,
	}, nil
}

func (u *playlistRequest) GetCurrentUserPlaylists(limit int, offset int) (*currentUserPlaylistResponse, error) {
	var userRes *currentUserPlaylistResponse

	if limit > 50 {
		logger.Println("Maximum limit for request is 50")
		return nil, errors.New("limit is over 50")
	} else if limit < 1 {
		logger.Println("Minimum limit for request is 1")
		return nil, errors.New("limit is below 1")
	}

	client := &http.Client{}
	urlVal := &url.Values{}

	urlVal.Add("limit", fmt.Sprint(limit))
	urlVal.Add("offset", fmt.Sprint(offset))

	req, err := http.NewRequest("GET", apiURL+"/me/playlists?"+urlVal.Encode(), nil)

	if err != nil {
		logger.Printf("Couldn't create request: %v\n", err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+u.AccessToken)
	fmt.Printf("%v\n", req)
	res, err := client.Do(req)

	if err != nil {
		logger.Printf("Couldn't perform GET request: %v\n", err)
		return nil, err
	}

	resBytes, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Printf("Coudln't read bytes from response body: %v\n", err)
		return nil, err
	}

	if err = json.Unmarshal(resBytes, userRes); err != nil {
		logger.Printf("Couldn't parse JSON data from response body: %v\n", err)
		fmt.Print(string(resBytes))
		return nil, err
	}

	return userRes, nil
}
