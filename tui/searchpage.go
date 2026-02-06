package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	HeaderHeight    = 3
	SearchBarHeight = 3
	FooterHeight    = 4
	ChromeHeight    = HeaderHeight + SearchBarHeight + FooterHeight
)

var (
	ItemsPerPage = 4
)

var (
	styleSearchBar    lipgloss.Style
	styleSearchDimmed lipgloss.Style
	styleTabActive    lipgloss.Style
	styleTabInactive  lipgloss.Style
	styleCardActive   lipgloss.Style
	styleCardInactive lipgloss.Style
	styleName         lipgloss.Style
	styleDesc         lipgloss.Style
	styleStats        lipgloss.Style
)

type progressMsg float64

type progressWriter struct {
	ch chan float64
}

func (m *SearchPageModel) Resize(width, height int) {
	m.Width = width
	m.Height = height

	m.ProgressBar.Width = width - 4

	styleSearchBar = styleSearchBar.Width(width - 4)
	styleSearchDimmed = styleSearchDimmed.Width(width - 4)

	vpHeight := m.Height - ChromeHeight
	if vpHeight < 5 {
		vpHeight = 5
	}
	m.Viewport.Height = vpHeight
	m.Viewport.Width = m.Width

	ItemsPerPage = m.Viewport.Height / 5
	if ItemsPerPage < 1 {
		ItemsPerPage = 1
	}

	m.Viewport.SetContent(m.renderContentString())
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	if pw.ch != nil {
		pw.ch <- 0.05
	}
	return len(p), nil
}

func waitForProgress(ch chan float64) tea.Cmd {
	return func() tea.Msg {
		if ch == nil {
			return nil
		}
		p, ok := <-ch
		if !ok {
			return nil
		}
		return progressMsg(p)
	}
}

type SearchPageModel struct {
	Mode       int
	SearchType int
	SearchBar  textinput.Model
	Result     utils.SearchResult

	Cursor      int
	WindowStart int

	Width   int
	Height  int
	Spinner spinner.Model
	Loading bool

	CloneProgress io.Writer

	ProgressBar  progress.Model
	IsCloning    bool
	progressChan chan float64

	Viewport viewport.Model
}

func NewSearchPageModel() SearchPageModel {

	styleSearchBar = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlight).
		Padding(0, 1)

	styleSearchDimmed = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle).
		Padding(0, 1)

	styleTabActive = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(special).
		Foreground(special).
		Bold(true).
		Padding(0, 2)

	styleTabInactive = lipgloss.NewStyle().
		Foreground(subtle).
		Padding(0, 2)

	styleCardActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlight).
		Padding(0, 1).
		MarginBottom(0)

	styleCardInactive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle).
		Padding(0, 1).
		MarginBottom(0)

	styleName = lipgloss.NewStyle().Bold(true).Foreground(text)
	styleDesc = lipgloss.NewStyle().Italic(true).Foreground(subtle)
	styleStats = lipgloss.NewStyle().Foreground(special)

	ti := textinput.New()
	ti.Placeholder = "Search GitHub..."
	ti.Focus()
	ti.CharLimit = 156

	ti.TextStyle = lipgloss.NewStyle().Foreground(text)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(subtle)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(highlight)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(highlight)

	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	viewport := viewport.New(0, 0)

	return SearchPageModel{
		Mode:        SearchMode,
		SearchType:  RepoMode,
		SearchBar:   ti,
		Result:      utils.SearchResult{},
		Cursor:      0,
		WindowStart: 0,
		Width:       20,
		Height:      20,
		Spinner:     s,
		Loading:     false,
		ProgressBar: prog,
		IsCloning:   false,
		Viewport:    viewport,
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

func (m SearchPageModel) renderRepoCard(repo githubapi.Repository, isActive bool, width int) string {
	style := styleCardInactive
	if isActive {
		style = styleCardActive
	}

	innerStyle := style.Copy().Width(width - 4)

	desc := ""
	if repo.Description != nil {
		desc = *repo.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
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
	footer := lipgloss.NewStyle().Foreground(subtle).Render(fmt.Sprintf("%s • Updated %s", lang, repo.UpdatedAt.Format("02 Jan")))

	content := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
	return innerStyle.Render(content)
}

func (m SearchPageModel) renderUserCard(user githubapi.UserSummary, isActive bool, width int) string {
	style := styleCardInactive
	if isActive {
		style = styleCardActive
	}
	innerStyle := style.Copy().Width(width - 4)

	content := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.NewStyle().Foreground(special).Render("(o) "),
		styleName.Render(user.Login),
		lipgloss.NewStyle().MarginLeft(2).Foreground(subtle).Render("View Profile →"),
	)

	return innerStyle.Render(content)
}

func (m SearchPageModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.Spinner.Tick)
}

func (m SearchPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		m.ProgressBar.Width = msg.Width - 4
		styleSearchBar = styleSearchBar.Width(msg.Width - 4)
		styleSearchDimmed = styleSearchDimmed.Width(msg.Width - 4)

		vpHeight := m.Height - ChromeHeight
		if vpHeight < 5 {
			vpHeight = 5
		}

		m.Viewport.Height = vpHeight
		m.Viewport.Width = m.Width

		ItemsPerPage = m.Viewport.Height / 5
		if ItemsPerPage < 1 {
			ItemsPerPage = 1
		}

	case progress.FrameMsg:
		progressModel, cmd := m.ProgressBar.Update(msg)
		m.ProgressBar = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

	case progressMsg:
		if m.ProgressBar.Percent() >= 1.0 {

			m.IsCloning = false
			m.ProgressBar.SetPercent(0)
			return m, nil
		}

		cmd := m.ProgressBar.IncrPercent(float64(msg))
		cmds = append(cmds, tea.Batch(cmd, waitForProgress(m.progressChan)), cmd)

	case spinner.TickMsg:
		if m.Loading {
			m.Spinner, cmd = m.Spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	case utils.SearchResult:
		m.Loading = false
		m.Result = msg
		m.Cursor = 0
		m.WindowStart = 0

		ItemsPerPage = m.Viewport.Height / 5
		return m, nil

	case tea.KeyMsg:
		if m.Mode == NavigationMode {
			switch msg.String() {
			case "j", "down":
				m.moveCursor(1)

				m.Viewport.ScrollDown(1)
			case "k", "up":
				m.moveCursor(-1)
				m.Viewport.ScrollUp(1)
			case "i", "/":
				m.Mode = SearchMode
				m.SearchBar.Focus()
				return m, textinput.Blink
			case "m":
				return m, func() tea.Msg {
					return NavMsg{
						to:   CreateRepoPage,
						from: SearchPage,
					}
				}
			case "c":
				if m.SearchType == RepoMode && !m.IsCloning {
					m.IsCloning = true
					m.ProgressBar.SetPercent(0)

					m.progressChan = make(chan float64)
					pw := &progressWriter{ch: m.progressChan}
					m.CloneProgress = pw
					name := m.Result.Repos.Items[m.Cursor].Name
					go githubapi.CloneURL(m.Result.Repos.Items[m.Cursor].CloneURL, name, &m.CloneProgress)

					cmds = append(cmds, waitForProgress(m.progressChan))
				}
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
						return NavMsg{to: UserPage, from: SearchPage, userdata: m.Result.Users.Items[m.Cursor]}
					}
				} else {
					return m, func() tea.Msg {
						return NavMsg{to: RepoPage, from: SearchPage, repodata: m.Result.Repos.Items[m.Cursor]}
					}
				}
			}
		}

		if m.Mode == SearchMode {
			switch msg.String() {
			case "enter":
				m.Mode = NavigationMode
				m.SearchBar.Blur()
				m.Loading = true
				cmds = append(cmds, tea.Batch(m.getQueryResult(m.SearchBar.Value()), m.Spinner.Tick))
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

	content := m.renderContentString()
	m.Viewport.SetContent(content)
	m.Viewport.SetYOffset(0)

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m SearchPageModel) renderContentString() string {
	var listItems []string

	if m.getListLength() > 0 {
		if m.WindowStart >= m.getListLength() {
			m.WindowStart = 0
		}

		endIndex := m.WindowStart + ItemsPerPage
		if endIndex > m.getListLength() {
			endIndex = m.getListLength()
		}

		if m.SearchType == UserMode {
			subset := m.Result.Users.Items[m.WindowStart:endIndex]
			for i, item := range subset {
				isSelected := (m.WindowStart + i) == m.Cursor

				listItems = append(listItems, m.renderUserCard(item, isSelected, m.Viewport.Width))
			}
		} else {
			subset := m.Result.Repos.Items[m.WindowStart:endIndex]
			for i, item := range subset {
				isSelected := (m.WindowStart + i) == m.Cursor

				listItems = append(listItems, m.renderRepoCard(item, isSelected, m.Viewport.Width))
			}
		}
	} else {
		listItems = append(listItems, lipgloss.NewStyle().Padding(2).Foreground(subtle).Render("No results found."))
	}

	return lipgloss.JoinVertical(lipgloss.Left, listItems...)
}

func (m SearchPageModel) View() string {

	if m.Loading {
		if m.SearchType == UserMode {
			return lipgloss.Place(
				m.Width, m.Height,
				lipgloss.Center, lipgloss.Center,
				fmt.Sprintf("%s Loading Users ...", m.Spinner.View()),
			)
		} else {
			return lipgloss.Place(
				m.Width, m.Height,
				lipgloss.Center, lipgloss.Center,
				fmt.Sprintf("%s Loading Repos ...", m.Spinner.View()),
			)
		}
	}

	var tUser, tRepo string
	if m.SearchType == UserMode {
		tUser = styleTabActive.Render("Users")
		tRepo = styleTabInactive.Render("Repositories")
	} else {
		tUser = styleTabInactive.Render("Users")
		tRepo = styleTabActive.Render("Repositories")
	}

	header := lipgloss.NewStyle().Height(HeaderHeight).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, tUser, tRepo),
	)

	var sBar string
	if m.Mode == SearchMode {
		sBar = styleSearchBar.Render(m.SearchBar.View())
	} else {
		sBar = styleSearchDimmed.Render(m.SearchBar.View())
	}
	sBar = lipgloss.NewStyle().Height(SearchBarHeight).Render(sBar)

	var footerStatus string
	if m.IsCloning {
		footerStatus = lipgloss.JoinVertical(lipgloss.Left,
			styleDesc.Render("Cloning Repository..."),
			m.ProgressBar.View(),
		)
	} else {
		page := (m.WindowStart / ItemsPerPage) + 1
		footerStatus = styleDesc.Render(fmt.Sprintf("Page %d", page))
	}

	footerStatus = lipgloss.NewStyle().Height(FooterHeight).Render(footerStatus)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		sBar,
		m.Viewport.View(),
		footerStatus,
	)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
