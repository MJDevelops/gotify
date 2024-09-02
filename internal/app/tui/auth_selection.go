package tui

import (
	"fmt"

	"github.com/MJDevelops/gotify/internal/app/spotifyflow"
	tea "github.com/charmbracelet/bubbletea"
)

type authSelect struct {
	choices    []string
	prevCursor int
	cursor     int
	selected   map[int]struct{}
}

func InitialAuthSelect() *authSelect {
	return &authSelect{
		choices:  []string{"Authorization Code\n", "Client Credential (restricted access only)\n"},
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

func (m *authSelect) Init() tea.Cmd {
	m.selected[m.cursor] = struct{}{}
	return tea.SetWindowTitle("gotify")
}

func (m *authSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.prevCursor = m.cursor
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.prevCursor = m.cursor
				m.cursor++
			}
		case "enter", " ":
			m.selected[m.cursor] = struct{}{}
			delete(m.selected, m.prevCursor)
		case "c":
			for key := range m.selected {
				if _, ok := m.selected[key]; ok {
					var authCode spotifyflow.SpotifyFlow
					if key == 0 {
						authCode = &spotifyflow.SpotifyAuthorizationCode{}
					} else {
						authCode = &spotifyflow.SpotifyClientCredential{}
					}

					authCode.Authorize()
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
}

func (m *authSelect) View() string {
	s := "Select authentication flow\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s", cursor, checked, choice)
	}

	s += "\nq to quit, c to confirm\n"

	return s
}
