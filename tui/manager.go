package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/githubapi"
)

var (
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

const (
	HomePage int = iota
	SearchPage
	RepoPage
	UserPage
)

type RepoLoaded struct {
	Data githubapi.Repository
}

type UserLoaded struct {
	Data githubapi.UserSummary
}

type NavMsg struct {
	to       int
	from     int
	userdata githubapi.UserSummary
	repodata githubapi.Repository
}

type Manager struct {
	page   tea.Model
	Width  int
	Height int
}

func NewManager() Manager {
	return Manager{
		page: NewHomePageModel(),
	}
}

func (m Manager) Init() tea.Cmd {
	return nil
}

func (m Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case NavMsg:
		switch msg.to {
		case HomePage:
			m.page = NewHomePageModel()

		case SearchPage:
			m.page = NewSearchPageModel()
			return m, m.page.Init()
		case RepoPage:
			m.page = NewRepoPageModel(msg.repodata, msg.userdata, msg.from)

		case UserPage:
			m.page = NewUserPageModel(msg.userdata, msg.from)
			return m, m.page.Init()
		}
		return m, nil
	}

	m.page, cmd = m.page.Update(msg)
	return m, cmd
}

func (m Manager) View() string {
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		m.page.View(),
	)
}
