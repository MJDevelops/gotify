package tui

import tea "github.com/charmbracelet/bubbletea"

const (
  playlists = iota
  albums
)

type gotifyHome struct {
  currentTab int
  cursor int
  tabs []string
}

func InitialGotifyHome() *gotifyHome {
  return &gotifyHome{
    currentTab: playlists,
    cursor: playlists,
    tabs: []string{"Playlists" , "Albums"},
  }
}


func (g *gotifyHome) Init() tea.Cmd {
    return tea.Batch(tea.ClearScreen, tea.SetWindowTitle("gotify"))
}
