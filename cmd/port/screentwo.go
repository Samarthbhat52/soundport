package port

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

func (m *portModel) updatePortProgress(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case playlistSelected:
		// Get all the tracks present in source playlist.
		tracks, err := m.src.GetPlaylistTracks(m.selectedPlId)
		if err != nil {
			glbLogger.Println("err: ", err.Error())
			m.quitting = true
			return m, tea.Quit
		}

		if ok := m.dst.AddTracks(m.createdPlId, tracks); ok {
			m.successful = true
			m.quitting = true
			// Add a something that will show that playlist successfully created.
			return m, tea.Quit
		}

		glbLogger.Println("unable to add tracks")
		m.quitting = true
		return m, tea.Quit

	default:
		return m, nil
	}
}

func (m *portModel) viewPortProgress() string {
	s := "\n"
	s += strings.Repeat(" ", padding)

	s += "Finding songs..."

	if m.successful {
		s = "Completed porting tracks."
	}

	if m.quitting {
		s += "\n"
	}

	return s
}
