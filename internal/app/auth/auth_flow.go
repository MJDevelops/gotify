package auth

type SpotifyFlow interface {
	Authorize() error
}
