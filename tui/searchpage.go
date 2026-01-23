package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/githubapi"
	"github.com/chirag-diwan/RemGit/utils"
)

/*
States{
	Navigation
	User
	Repository
	Search
}
User Repository
(Search bar)TextInput
space to show search result
*/

const (
	NavigationMode int = iota
	SearchMode
)

const (
	UserMode int = iota
	RepoMode
)

var (
	colWhite = lipgloss.Color("255")

	searchBar = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(colText).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colHighlight).
			Padding(0, 1)

	tabBar = lipgloss.NewStyle().
		PaddingTop(1).
		PaddingBottom(1).
		Align(lipgloss.Center)

	activeTab = lipgloss.NewStyle().
			Foreground(colHighlight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colHighlight).
			Padding(0, 2).
			Bold(true)

	inactiveTab = lipgloss.NewStyle().
			Foreground(colText).
			Border(lipgloss.HiddenBorder()).
			Padding(0, 2)

	tabGap = lipgloss.NewStyle().Width(2)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Border(lipgloss.HiddenBorder(), false, false, false, true)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(colHighlight).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(colHighlight)

	repoNameStyle = lipgloss.NewStyle().Bold(true)
	descStyle     = lipgloss.NewStyle().Foreground(colText).Italic(true)
	statStyle     = lipgloss.NewStyle().Foreground(colText).MarginRight(2)
)

type searchPage struct {
	Mode       int
	SearchType int
	SearchBar  textinput.Model
	Result     utils.SearchResult
	Cursor     int
	Width      int
	Height     int
}

func NewSearchPage() searchPage {
	txtinput := textinput.New()
	txtinput.Placeholder = "Search Github..."
	txtinput.Focus()
	txtinput.Cursor.Blink = true
	txtinput.Width = 40
	return searchPage{
		SearchMode,
		UserMode,
		txtinput,
		utils.SearchResult{},
		0,
		10,
		10,
	}
}

func (model searchPage) getQueryResult(query string) tea.Cmd {
	return func() tea.Msg {

		var res utils.SearchResult
		if model.SearchType == UserMode {
			res = utils.SearchResult{Users: githubapi.GetUsers(query), Repos: githubapi.RepoSearchResponse{}}
		} else {
			res = utils.SearchResult{Users: githubapi.UserSearchResponse{}, Repos: githubapi.GetRepos(query)}
		}

		return res
	}
}

func (model searchPage) showResult() string {
	var listContent []string

	safeStr := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	if model.SearchType == UserMode {
		users := model.Result.Users.Items
		if len(users) == 0 {
			return lipgloss.NewStyle().Padding(2).Foreground(colText).Render("No users found.")
		}

		for i, user := range users {

			style := itemStyle
			cursorStr := "  "
			if i == model.Cursor {
				style = selectedItemStyle
				cursorStr = "> "
			}

			row := fmt.Sprintf("%s%s", cursorStr, user.Login)
			listContent = append(listContent, style.Render(row))
		}

	} else {
		repos := model.Result.Repos.Items
		if len(repos) == 0 {
			return lipgloss.NewStyle().Padding(2).Foreground(colText).Render("No repositories found.")
		}

		for i, repo := range repos {

			containerStyle := itemStyle
			if i == model.Cursor {
				containerStyle = selectedItemStyle
			}

			name := repoNameStyle.Render(repo.FullName)

			desc := safeStr(repo.Description)
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			description := descStyle.Render(desc)

			lang := safeStr(repo.Language)
			if lang == "" {
				lang = "Text"
			}
			stats := statStyle.Render(fmt.Sprintf("★ %d  •  %s", repo.StargazersCount, lang))

			block := lipgloss.JoinVertical(
				lipgloss.Left,
				name,
				description,
				stats,
			)

			listContent = append(listContent, containerStyle.Render(block))
			listContent = append(listContent, " ")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, listContent...)
}

func (model searchPage) getTabBarContent() string {
	var userTab, repoTab string

	if model.SearchType == UserMode {
		userTab = activeTab.Render("User")
		repoTab = inactiveTab.Render("Repo")
	} else {
		userTab = inactiveTab.Render("User")
		repoTab = activeTab.Render("Repo")
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		userTab,
		tabGap.Render(""),
		repoTab,
	)
}

func (model searchPage) GetListLength() int {
	if model.SearchType == UserMode {
		return len(model.Result.Users.Items)
	}
	return len(model.Result.Repos.Items)
}

func (model searchPage) clampCursor() {
	maxLength := model.GetListLength()
	if maxLength == 0 {
		model.Cursor = 0
		return
	}
	if model.Cursor >= maxLength {
		model.Cursor = maxLength - 1
	} else if model.Cursor < 0 {
		model.Cursor = 0
	}
}

func (model searchPage) switchMode() searchPage {
	if model.SearchType == UserMode {
		model.SearchType = RepoMode
	} else {
		model.SearchType = UserMode
	}
	model.Cursor = 0
	return model
}

func (model searchPage) Init() tea.Cmd {
	return textinput.Blink
}

func (model searchPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		model.Width = msg.Width - 2
		model.Height = msg.Height - 2

	case utils.SearchResult:
		model.Result = msg
		model.Cursor = 0
		return model, nil

	case tea.KeyMsg:
		switch model.Mode {

		case NavigationMode:
			switch msg.String() {
			case "i":
				model.Mode = SearchMode
				model.SearchBar.Focus()
				return model, textinput.Blink

			case "tab":
				model = model.switchMode()

			case "j", "down":
				model.Cursor++
				model.clampCursor()

			case "k", "up":
				model.Cursor--
				model.clampCursor()

			case "enter":
				if model.SearchType == UserMode {

				} else {
					showRepoPage(model.Result.Repos.Items[model.Cursor])
				}
			case "q", "ctrl+c":
				return model, tea.Quit
			}

		case SearchMode:
			switch msg.String() {
			case "tab":
				model = model.switchMode()
				return model, nil

			case "esc":
				model.Mode = NavigationMode
				model.SearchBar.Blur()
				return model, nil

			case "enter":
				query := model.SearchBar.Value()

				searchCmd := model.getQueryResult(query)
				cmds = append(cmds, searchCmd)

				model.SearchBar.Blur()
				model.Mode = NavigationMode
			}
			model.SearchBar, cmd = model.SearchBar.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return model, tea.Batch(cmds...)
}

func (model searchPage) View() string {
	searchView := searchBar.Render(model.SearchBar.View())
	tabBarContent := model.getTabBarContent()

	listView := model.showResult()

	ui := lipgloss.JoinVertical(
		lipgloss.Center,
		tabBar.Render(tabBarContent),
		searchView,
		"\n",
		listView,
	)

	return lipgloss.Place(
		model.Width,
		model.Height,
		lipgloss.Center,
		lipgloss.Top,
		ui,
	)
}
