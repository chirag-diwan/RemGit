package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	colSubtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	colHighlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	colText      = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#E0E0E0"}
	colDim       = lipgloss.AdaptiveColor{Light: "#A0A0A0", Dark: "#555555"}

	window = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(colSubtle).
		Align(lipgloss.Center).
		Align(lipgloss.Center).
		Padding(2)
	heading = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(colHighlight).
		Foreground(colHighlight).
		Align(lipgloss.Center)
)

var Logo = "██████  ███████ ███    ███  ██████  ██ ████████\n██   ██ ██      ████  ████ ██       ██    ██   \n██████  █████   ██ ████ ██ ██   ███ ██    ██   \n██   ██ ██      ██  ██  ██ ██    ██ ██    ██   \n██   ██ ███████ ██      ██  ██████  ██    ██   "

type HomePage struct {
	Logo   string
	Width  int
	Height int
}

func (model HomePage) Init() tea.Cmd {
	return nil
}

func (model HomePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.Height = msg.Height - 2
		model.Width = msg.Width - 3
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return model, tea.Quit
		case "s":
			//search mode
		case "c":
			//config mode
		}
	}
	return model, nil
}

func (model HomePage) View() string {
	message := `

Press S to search user or repository :)

Press h for help

You can disable this screen in config files (~/.remgit.conf)
	`
	return lipgloss.Place(
		model.Width,
		model.Height,
		lipgloss.Center,
		lipgloss.Center,
		window.Height(model.Height).Width(model.Width).Render(heading.Render(Logo)+message),
	)
}
