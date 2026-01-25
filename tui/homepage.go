package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	heading lipgloss.Style
)

var Logo = "██████  ███████ ███    ███  ██████  ██ ████████\n██   ██ ██      ████  ████ ██       ██    ██   \n██████  █████   ██ ████ ██ ██   ███ ██    ██   \n██   ██ ██      ██  ██  ██ ██    ██ ██    ██   \n██   ██ ███████ ██      ██  ██████  ██    ██   "

type HomePageModel struct {
	Logo   string
	Width  int
	Height int
}

func NewHomePageModel() HomePageModel {
	heading = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(highlight).
		Foreground(highlight).
		Align(lipgloss.Center)

	return HomePageModel{}
}

func (model HomePageModel) Init() tea.Cmd {
	return nil
}

func (model HomePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.Height = msg.Height - 2
		model.Width = msg.Width - 3
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return model, tea.Quit
		case "s":
			return model, func() tea.Msg {
				return NavMsg{to: SearchPage, from: HomePage}
			}
		case "c":
		}
	}
	return model, nil
}

func (model HomePageModel) View() string {
	message := `

Press s to search user or repository :)

Press h for help

You can disable this screen in config files (~/.remgit.conf)
	`
	return lipgloss.JoinVertical(
		lipgloss.Center,
		window.Height(model.Height).Width(model.Width).Render(heading.Render(Logo)+message),
	)
}
