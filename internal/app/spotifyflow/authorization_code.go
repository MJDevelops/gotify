package spotifyflow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/MJDevelops/gotify/internal/pkg/envs"
	"github.com/MJDevelops/gotify/internal/pkg/logs"
	"github.com/MJDevelops/gotify/pkg/browser"
	"github.com/google/go-querystring/query"
)

type SpotifyAuthorizationCode struct {
	AccessToken  string      `json:"access_token"`
	TokenType    string      `json:"token_type"`
	Scope        string      `json:"scope"`
	ExpiresIn    json.Number `json:"expires_in"`
	RefreshToken string      `json:"refresh_token"`
}

type spotifyExchangeCodeRequest struct {
	ClientID            string `url:"client_id"`
	ResponseType        string `url:"response_type"`
	RedirectUri         string `url:"redirect_uri"`
	Scope               string `url:"scope"`
	State               string `url:"state"`
	CodeChallengeMethod string `url:"code_challenge_method"`
	CodeChallenge       string `url:"code_challenge"`
}

type spotifyExchangeCodeResponse struct {
	Code string
}

type spotifyAuthorizationCodeRequest struct {
	Code         string `url:"code"`
	GrantType    string `url:"grant_type"`
	RedirectUri  string `url:"redirect_uri"`
	ClientID     string `url:"client_id"`
	CodeVerifier string `url:"code_verifier"`
}

const scopes = "user-read-playback-state user-modify-playback-state " +
	"user-read-currently-playing app-remote-control " +
	"streaming playlist-read-private " +
	"playlist-read-collaborative playlist-modify-private " +
	"playlist-modify-public user-library-modify " +
	"user-library-read"

const redirectUri = "http://localhost:8888/callback"
const spotifyAuthorizeURL = "https://accounts.spotify.com/authorize?"
const spotifyTokenReqURL = "https://accounts.spotify.com/api/token"

var closeWg sync.WaitGroup
var resCh = make(chan url.Values)
var env = envs.LoadEnv()
var logger = logs.GetLoggerInstance()

func (s *SpotifyAuthorizationCode) Authorize() error {
	req, codeVerifier, err := newExchangeCodeRequest()

	if err != nil {
		logger.Printf("Couldn't create code request: %v\n", err)
		return err
	}

	urlVals, err := query.Values(*req)

	if err != nil {
		logger.Println("Couldn't create URL params from auth struct")
		return err
	}

	urlStr := fmt.Sprintf(spotifyAuthorizeURL+"%s", urlVals.Encode())

	browser.Open(urlStr)
	go waitForExchangeCode()
	urlVals = <-resCh

	if urlVals == nil {
		logger.Fatal("Couldn't get url params from callback: value is nil\n")
	}

	exchangeCodeResponse := spotifyExchangeCodeResponse{
		Code: urlVals.Get("code"),
	}

	authCodeRequest, err := newAuthorizationCodeRequest(exchangeCodeResponse, codeVerifier)

	if err != nil {
		logger.Printf("Couldn't create Authorization Code Request in KickstartAuthorizationCodeRequest(): %v\n", err)
		return err
	}

	err = requestAuthorizationCode(authCodeRequest, s)

	if err != nil {
		logger.Printf("Error during Authorization Code Request: %v\n", err)
		return err
	}

	return nil
}

func newExchangeCodeRequest() (*spotifyExchangeCodeRequest, string, error) {
	codeVerifier, err := randomBytesInHex(32)

	if err != nil {
		logger.Printf("Couldn't generate code verifier: %v\n", err)
		return nil, "", err
	}

	sha := sha256.New()
	io.WriteString(sha, codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha.Sum(nil))
	state, err := randomBytesInHex(24)

	return &spotifyExchangeCodeRequest{
		RedirectUri:         redirectUri,
		ResponseType:        "code",
		ClientID:            env.GotifyClientID,
		Scope:               scopes,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: "S256",
	}, codeVerifier, err
}

func newAuthorizationCodeRequest(exchangeCodeReq spotifyExchangeCodeResponse, verifier string) (*spotifyAuthorizationCodeRequest, error) {
	urlVals, err := query.Values(exchangeCodeReq)
	code := urlVals.Get("Code")

	if err != nil {
		logger.Printf("Couldn't parse response object in newAuthorizationCodeRequest(): %v\n", err)
		return nil, err
	}

	return &spotifyAuthorizationCodeRequest{
		Code:         code,
		RedirectUri:  redirectUri,
		GrantType:    "authorization_code",
		CodeVerifier: verifier,
		ClientID:     env.GotifyClientID,
	}, nil
}

func requestAuthorizationCode(authReq *spotifyAuthorizationCodeRequest, s *SpotifyAuthorizationCode) error {
	urlVals, err := query.Values(*authReq)
	client := &http.Client{}

	if err != nil {
		logger.Printf("Couldn't parse auth code request in requestAuthorizationCode(): %v\n", err)
		return err
	}

	httpReq, err := http.NewRequest("POST", spotifyTokenReqURL, strings.NewReader(urlVals.Encode()))

	if err != nil {
		logger.Printf("Couldn't create request in requestAuthorizationCode(): %v\n", err)
		return err
	}

	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(httpReq)

	if err != nil {
		logger.Printf("Couldn't perform 'POST' request in requestAuthorizationCode(): %v\n", err)
		return err
	}

	if res.StatusCode != 200 {
		logger.Printf("Authorization Code Response is not 200: %d\n", res.StatusCode)
		return errors.New("response code is not 200")
	}

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Printf("Couldn't read from response body in requestAuthorizationCode(): %v\n", err)
		return err
	}

	if err = json.Unmarshal(resBody, s); err != nil {
		logger.Printf("Couldn't parse JSON from response body in requestAuthorizationCode(): %v\n", err)
		return err
	}

	return nil
}

func waitForExchangeCode() {
	srv := &http.Server{Addr: ":8888"}
	http.HandleFunc("/callback", handleCallback)
	closeWg.Add(1)

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalf("ListenAndServe() %v\n", err)
	}

	logger.Println("Waiting for request...")

	closeWg.Wait()

	logger.Println("Request received")
	srv.Close()
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	defer closeWg.Done()

	urlVal, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		logger.Printf("Something went wrong: %v\n", err)
		resCh <- nil
	}

	resCh <- urlVal
	io.WriteString(w, "Authorized, you can close this tab")
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("could not generate %d random bytes: %v", count, err)
	}

	return hex.EncodeToString(buf), nil
}
