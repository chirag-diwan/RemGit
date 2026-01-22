package utils

import (
	"fmt"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/githubapi"
)

var Logo = `
  _____               _____ _ _
 |  __ \             / ____(_) |
 | |__) |___ _ __ __| |  __ _| |_
 |  _  
 | | \ \  __/ | | | | |__| | | |_
 |_|  \_\___|_| |_| |_|\_____|_|\__|
`

var (
	colSubtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	colHighlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	colText      = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#E0E0E0"}
	colDim       = lipgloss.AdaptiveColor{Light: "#A0A0A0", Dark: "#555555"}

	styleWindow = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colSubtle).
			Align(lipgloss.Center)

	styleInputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(60).
			Align(lipgloss.Left)

	styleItem = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(colText)

	styleSelected = lipgloss.NewStyle().
			PaddingLeft(0).
			Foreground(colHighlight).
			Bold(true).
			SetString("> ")

	styleMenu = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colHighlight).
			Padding(1, 2).
			Align(lipgloss.Center)
)

type SearchResult struct {
	Users githubapi.UserSearchResponse
	Repos githubapi.RepoSearchResponse
}

type Menu struct {
	Active  bool
	Options []string
	Cursor  int
}

type Model struct {
	AppState     int
	Width        int
	Height       int
	SearchText   string
	InputMode    int
	SearchType   int
	Results      SearchResult
	ResultCursor int
	ActionMenu   Menu
}

func getMenuOptions(searchType int) []string {
	if searchType == SearchRepo {
		return []string{"Open Repo", "Clone Repo", "Copy Link", "Cancel"}
	}
	return []string{"View Profile", "Show Repositories", "Follow User", "Cancel"}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:

		if msg.String() == "tab" && m.AppState != StateMenu {
			if m.SearchType == SearchRepo {
				m.SearchType = SearchUser
			} else {
				m.SearchType = SearchRepo
			}

			m.ResultCursor = 0
			m.AppState = StateStart
			return m, nil
		}

		if m.AppState == StateMenu {
			switch msg.String() {
			case "j", "down":
				if m.ActionMenu.Cursor < len(m.ActionMenu.Options)-1 {
					m.ActionMenu.Cursor++
				}
			case "k", "up":
				if m.ActionMenu.Cursor > 0 {
					m.ActionMenu.Cursor--
				}
			case "enter":
				selected := m.ActionMenu.Options[m.ActionMenu.Cursor]

				if selected == "Cancel" {
					m.AppState = StateResults
				} else {
					m.AppState = StateResults
				}
			case "esc":
				m.AppState = StateResults
			}
			return m, nil
		}

		if m.InputMode == ModeInsert {
			switch msg.String() {
			case "enter":

				m.Results = PerformSearch(m.SearchText, m.SearchType)
				m.ResultCursor = 0

				hasRepos := m.SearchType == SearchRepo && len(m.Results.Repos.Items) > 0
				hasUsers := m.SearchType == SearchUser && len(m.Results.Users.Items) > 0

				if hasRepos || hasUsers {
					m.AppState = StateResults
					m.InputMode = ModeNormal
				}
			case "esc":
				m.InputMode = ModeNormal
			case "backspace":
				if len(m.SearchText) > 0 {
					m.SearchText = m.SearchText[:len(m.SearchText)-1]
				}
			default:
				if len(msg.Runes) == 1 && unicode.IsPrint(msg.Runes[0]) {
					m.SearchText += string(msg.Runes[0])
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "i":
			m.InputMode = ModeInsert
		case "q", "ctrl+c":
			return m, tea.Quit

		case "j", "down":
			if m.AppState == StateResults {
				maxIndex := 0
				if m.SearchType == SearchRepo {
					maxIndex = len(m.Results.Repos.Items) - 1
				} else {
					maxIndex = len(m.Results.Users.Items) - 1
				}

				if m.ResultCursor < maxIndex {
					m.ResultCursor++
				}
			}

		case "k", "up":
			if m.AppState == StateResults && m.ResultCursor > 0 {
				m.ResultCursor--
			}

		case "enter":

			if m.AppState == StateResults {
				m.AppState = StateMenu
				m.ActionMenu.Options = getMenuOptions(m.SearchType)
				m.ActionMenu.Cursor = 0
			}

		case "esc":
			if m.AppState == StateResults {
				m.AppState = StateStart

			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Width == 0 {
		return "loading..."
	}

	tabRepo := styleItem.Copy().Foreground(colDim).Render("Repository")
	tabUser := styleItem.Copy().Foreground(colDim).Render("User")

	if m.SearchType == SearchRepo {
		tabRepo = styleSelected.Render("Repository")
	} else {
		tabUser = styleSelected.Render("User")
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top, tabRepo, "  ", tabUser)

	var content string
	availableHeight := m.Height - 10

	if m.AppState == StateStart {
		content = lipgloss.Place(m.Width-4, availableHeight, lipgloss.Center, lipgloss.Center, Logo)
	} else {
		var rows []string

		if m.SearchType == SearchRepo {
			for i, repo := range m.Results.Repos.Items {
				cursor := "  "
				style := styleItem
				if m.ResultCursor == i {
					cursor = "> "
					style = styleSelected
				}

				desc := ""
				if repo.Description != nil {

					if len(*repo.Description) > 50 {
						desc = (*repo.Description)[:47] + "..."
					} else {
						desc = *repo.Description
					}
				}

				row := fmt.Sprintf("%s%s %s", cursor, style.Render(repo.FullName), lipgloss.NewStyle().Foreground(colDim).Render(desc))
				rows = append(rows, row)
			}
		} else {

			for i, user := range m.Results.Users.Items {
				cursor := "  "
				style := styleItem
				if m.ResultCursor == i {
					cursor = "> "
					style = styleSelected
				}
				row := fmt.Sprintf("%s%s", cursor, style.Render(user.Login))
				rows = append(rows, row)
			}
		}

		if len(rows) == 0 {
			rows = append(rows, "No results found.")
		}

		listContent := lipgloss.JoinVertical(lipgloss.Left, rows...)
		content = lipgloss.NewStyle().Padding(1, 2).Render(listContent)
	}

	borderCol := colSubtle
	if m.InputMode == ModeInsert {
		borderCol = colHighlight
	}

	barText := m.SearchText
	if m.InputMode == ModeInsert {
		barText += "█"
	}

	searchBar := styleInputBox.Copy().
		BorderForeground(borderCol).
		Render(barText)

	baseView := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().PaddingTop(1).Render(header),
		lipgloss.Place(m.Width-4, availableHeight, lipgloss.Left, lipgloss.Top, content),
		searchBar,
	)

	if m.AppState == StateMenu {
		var menuRows []string
		for i, opt := range m.ActionMenu.Options {
			prefix := "  "
			style := styleItem
			if i == m.ActionMenu.Cursor {
				prefix = "> "
				style = styleSelected
			}
			menuRows = append(menuRows, fmt.Sprintf("%s%s", prefix, style.Render(opt)))
		}

		menuBlock := styleMenu.Render(lipgloss.JoinVertical(lipgloss.Left, menuRows...))

		contentWithMenu := lipgloss.Place(m.Width-4, availableHeight, lipgloss.Center, lipgloss.Center, menuBlock)

		baseView = lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().PaddingTop(1).Render(header),
			contentWithMenu,
			searchBar,
		)
	}

	return styleWindow.
		Width(m.Width - 2).
		Height(m.Height).
		Render(baseView)
}
