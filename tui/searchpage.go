package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/githubapi"
	"github.com/chirag-diwan/RemGit/utils"
)

const (
	NavigationMode int = iota
	SearchMode
)

const (
	UserMode int = iota
	RepoMode
)

const (
	ItemsPerPage = 4
)

var (
	colAccent  = lipgloss.Color("212")
	colSurface = lipgloss.Color("235")

	styleSearchBar = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colAccent).
			Padding(0, 1).
			Width(60)

	styleSearchDimmed = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colSubtle).
				Padding(0, 1).
				Width(60)

	styleTabActive = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(colAccent).
			Foreground(colAccent).
			Bold(true).
			Padding(0, 2)

	styleTabInactive = lipgloss.NewStyle().
				Foreground(colSubtle).
				Padding(0, 2)

	styleCardActive = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colAccent).
			Padding(0, 1).
			MarginBottom(1).
			Width(60)

	styleCardInactive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colSubtle).
				Padding(0, 1).
				MarginBottom(1).
				Width(60)

	styleName  = lipgloss.NewStyle().Bold(true).Foreground(colText)
	styleDesc  = lipgloss.NewStyle().Italic(true).Foreground(colSubtle)
	styleStats = lipgloss.NewStyle().Foreground(colAccent)
)

type SearchPageModel struct {
	Mode       int
	SearchType int
	SearchBar  textinput.Model
	Result     utils.SearchResult

	Cursor      int
	WindowStart int

	Width  int
	Height int
}

func NewSearchPageModel() SearchPageModel {
	ti := textinput.New()
	ti.Placeholder = "Search GitHub..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return SearchPageModel{
		Mode:        SearchMode,
		SearchType:  RepoMode,
		SearchBar:   ti,
		Result:      utils.SearchResult{},
		Cursor:      0,
		WindowStart: 0,
		Width:       20,
		Height:      20,
	}
}

func (m SearchPageModel) getQueryResult(query string) tea.Cmd {
	return func() tea.Msg {

		if m.SearchType == UserMode {
			return utils.SearchResult{Users: githubapi.GetUsers(query)}
		}
		return utils.SearchResult{Repos: githubapi.GetRepos(query)}
	}
}

func (m SearchPageModel) getListLength() int {
	if m.SearchType == UserMode {
		return len(m.Result.Users.Items)
	}
	return len(m.Result.Repos.Items)
}

func (m *SearchPageModel) moveCursor(step int) {
	total := m.getListLength()
	if total == 0 {
		return
	}

	newCursor := m.Cursor + step

	if newCursor < 0 {
		newCursor = 0
	} else if newCursor >= total {
		newCursor = total - 1
	}
	m.Cursor = newCursor

	if m.Cursor >= m.WindowStart+ItemsPerPage {
		m.WindowStart = m.Cursor - ItemsPerPage + 1
	}

	if m.Cursor < m.WindowStart {
		m.WindowStart = m.Cursor
	}
}

func (m SearchPageModel) renderRepoCard(repo githubapi.Repository, isActive bool) string {
	style := styleCardInactive
	if isActive {
		style = styleCardActive
	}

	desc := ""
	if repo.Description != nil {
		desc = *repo.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
	} else {
		desc = "No description provided."
	}

	lang := "Text"
	if repo.Language != nil {
		lang = *repo.Language
	}

	header := lipgloss.JoinHorizontal(lipgloss.Center,
		styleName.Render(repo.FullName),
		lipgloss.NewStyle().MarginLeft(2).Render(styleStats.Render(fmt.Sprintf("★ %d", repo.StargazersCount))),
	)

	body := styleDesc.Render(desc)
	footer := lipgloss.NewStyle().Foreground(colSubtle).Render(fmt.Sprintf("%s • Updated %s", lang, repo.UpdatedAt.Format("02 Jan")))

	content := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
	return style.Render(content)
}

func (m SearchPageModel) renderUserCard(user githubapi.UserSummary, isActive bool) string {
	style := styleCardInactive
	if isActive {
		style = styleCardActive
	}

	content := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.NewStyle().Foreground(colAccent).Render("(o) "),
		styleName.Render(user.Login),
		lipgloss.NewStyle().MarginLeft(2).Foreground(colSubtle).Render("View Profile →"),
	)

	return style.Render(content)
}

func (m SearchPageModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case utils.SearchResult:
		m.Result = msg
		m.Cursor = 0
		m.WindowStart = 0
		return m, nil

	case tea.KeyMsg:

		if m.Mode == NavigationMode {
			switch msg.String() {
			case "j", "down":
				m.moveCursor(1)
			case "k", "up":
				m.moveCursor(-1)
			case "i", "/":
				m.Mode = SearchMode
				m.SearchBar.Focus()
				return m, textinput.Blink
			case "tab":
				if m.SearchType == UserMode {
					m.SearchType = RepoMode
				} else {
					m.SearchType = UserMode
				}
				m.Cursor = 0
				m.WindowStart = 0
			case "enter":
				if m.SearchType == UserMode {
					return m, func() tea.Msg {
						return NavMsg{to: UserPage, userdata: m.Result.Users.Items[m.Cursor]}
					}
				} else {
					return m, func() tea.Msg {
						return NavMsg{to: RepoPage, repodata: m.Result.Repos.Items[m.Cursor]}
					}
				}
			}
		}

		if m.Mode == SearchMode {
			switch msg.String() {
			case "enter":
				m.Mode = NavigationMode
				m.SearchBar.Blur()
				return m, m.getQueryResult(m.SearchBar.Value())
			case "esc":
				m.Mode = NavigationMode
				m.SearchBar.Blur()
			case "tab":

				if m.SearchType == UserMode {
					m.SearchType = RepoMode
				} else {
					m.SearchType = UserMode
				}
				m.Cursor = 0
			}
			m.SearchBar, cmd = m.SearchBar.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m SearchPageModel) View() string {

	var tUser, tRepo string
	if m.SearchType == UserMode {
		tUser = styleTabActive.Render("Users")
		tRepo = styleTabInactive.Render("Repositories")
	} else {
		tUser = styleTabInactive.Render("Users")
		tRepo = styleTabActive.Render("Repositories")
	}
	header := lipgloss.JoinHorizontal(lipgloss.Top, tUser, tRepo)

	var sBar string
	if m.Mode == SearchMode {
		sBar = styleSearchBar.Render(m.SearchBar.View())
	} else {
		sBar = styleSearchDimmed.Render(m.SearchBar.View())
	}

	var listItems []string

	if m.getListLength() > 0 {

		endIndex := m.WindowStart + ItemsPerPage
		if endIndex > m.getListLength() {
			endIndex = m.getListLength()
		}

		if m.SearchType == UserMode {
			subset := m.Result.Users.Items[m.WindowStart:endIndex]
			for i, item := range subset {

				isSelected := (m.WindowStart + i) == m.Cursor
				listItems = append(listItems, m.renderUserCard(item, isSelected))
			}
		} else {
			subset := m.Result.Repos.Items[m.WindowStart:endIndex]
			for i, item := range subset {
				isSelected := (m.WindowStart + i) == m.Cursor
				listItems = append(listItems, m.renderRepoCard(item, isSelected))
			}
		}
	} else {
		listItems = append(listItems, lipgloss.NewStyle().Padding(2).Foreground(colSubtle).Render("No results found."))
	}

	listView := lipgloss.JoinVertical(lipgloss.Left, listItems...)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().MarginBottom(1).Render(header),
		lipgloss.NewStyle().MarginBottom(1).Render(sBar),
		listView,
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().PaddingTop(2).Render(content),
	)
}
