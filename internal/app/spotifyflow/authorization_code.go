package spotifyflow

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/MJDevelops/gotify/internal/pkg/envs"
	"github.com/MJDevelops/gotify/pkg/browser"
	"github.com/google/go-querystring/query"
)

type SpotifyAuthorizationCode struct {
	AccessToken string
}

type spotifyAuthorizationCodeRequest struct {
	ClientID     string `url:"client_id"`
	ResponseType string `url:"response_type"`
	RedirectUri  string `url:"redirect_uri"`
	Scope        string `url:"scope"`
}

type spotifyAuthorizationCodeResponse struct {
	Code  string
	State string
}

const scopes = "user-read-playback-state user-modify-playback-state " +
	"user-read-currently-playing app-remote-control " +
	"streaming playlist-read-private " +
	"playlist-read-collaborative playlist-modify-private " +
	"playlist-modify-public user-library-modify " +
	"user-library-read"

const spotifyAuthorizeURL = "https://accounts.spotify.com/authorize?"

var closeWg sync.WaitGroup
var resCh = make(chan url.Values)

func KickstartAuthorizationCodeRequest() error {
	req, err := newAuthorizationCodeRequest()

	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't initialize request")
		return err
	}

	urlVals, err := query.Values(*req)

	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't create URL params from auth struct")
		return err
	}

	urlStr := fmt.Sprintf(spotifyAuthorizeURL+"%s", urlVals.Encode())

	browser.Open(urlStr)
	go waitForAuth()
	urlVals = <-resCh

	if urlVals == nil {
		log.Fatal("Couldn't get url params from callback: value is nil")
	}

	fmt.Printf("%v", urlVals)

	return nil
}

func newAuthorizationCodeRequest() (*spotifyAuthorizationCodeRequest, error) {
	envs, err := envs.LoadEnv()

	return &spotifyAuthorizationCodeRequest{
		RedirectUri:  "localhost:8080/callback",
		ResponseType: "code",
		ClientID:     envs.GotifyClientID,
		Scope:        scopes,
	}, err
}

func waitForAuth() {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", handleCallback)
	closeWg.Add(1)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe() %v", err)
	}

	closeWg.Wait()
	srv.Close()
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	defer closeWg.Done()

	urlVal, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong: %v", err)
		resCh <- nil
	}

	resCh <- urlVal
	io.WriteString(w, "Authorized, you can close this tab")
}
