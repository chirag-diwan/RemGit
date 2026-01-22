package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/utils"
)

var (
	searchBar    = lipgloss.NewStyle()
	searchList   = lipgloss.NewStyle()
	searchWindow = lipgloss.NewStyle()
)

const (
	RepoMode   = 0
	UserMode   = 1
	Query      = 2
	Navigation = 3
)

type searchPage struct {
	Query              string
	SearchMode         int
	QueryResult        utils.SearchResult
	Mode               int
	Cursor             int
	Selected           int
	SearchWindowHeight int
	SearchWindowWidth  int
}

func (model searchPage) Init() tea.Cmd {
	return nil
}

func (model searchPage) Update(msg tea.Msg) {
	switch model.Mode {
	case Query:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			model.SearchWindowHeight = msg.Height
			model.SearchWindowWidth = msg.Width
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				model.Mode = Navigation
			case "backspace":
				model.Query = model.Query[:len(model.Query)-1]
			default:
				model.Query += msg.String()
			}
		}
	case Navigation:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			model.SearchWindowHeight = msg.Height
			model.SearchWindowWidth = msg.Width
		case tea.KeyMsg:
			switch msg.String() {
			case "j":
				if model.Cursor < len(model.QueryResult.Repos.Items)-1 || model.Cursor < len(model.QueryResult.Users.Items)-1 {
					model.Cursor++
				}
			case "k":
				if model.Cursor > 0 {
					model.Cursor--
				}
			case "enter":
				//call open repo page or open user page based on the flag that will be set during init
			}
		}
	}
}

func (model searchPage) View() {
}
