package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/githubapi"
)

var (
	docStyle  = lipgloss.NewStyle().Padding(1, 2)
	descStyle = lipgloss.NewStyle().Foreground(text).Italic(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight)

	statusBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(highlight).
			Padding(0, 1).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1).
			MarginRight(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Bold(true)

	statValStyle = lipgloss.NewStyle().
			Foreground(special).
			Bold(true)

	readmeBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)
)

type ReadmeMsg string

type RepoPageModel struct {
	Width         int
	Height        int
	CurrentRepo   githubapi.Repository
	CameFrom      int
	UserData      githubapi.UserSummary
	Viewport      viewport.Model
	ReadmeText    string
	LoadingReadme bool
}

func NewRepoPageModel(data githubapi.Repository, userdata githubapi.UserSummary, camefrom int, width int, height int) RepoPageModel {

	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle().Align(lipgloss.Left)

	return RepoPageModel{
		CurrentRepo:   data,
		CameFrom:      camefrom,
		UserData:      userdata,
		Viewport:      vp,
		LoadingReadme: true,
	}
}

func fetchReadmeCmd(repo githubapi.Repository) tea.Cmd {
	return func() tea.Msg {
		content := githubapi.GetReadme(repo)
		return ReadmeMsg(content)
	}
}

func (m RepoPageModel) Init() tea.Cmd {
	m.Viewport.SetContent(m.renderFullPage())

	return fetchReadmeCmd(m.CurrentRepo)
}

func (m RepoPageModel) renderFullPage() string {
	fullName := titleStyle.Render(m.CurrentRepo.FullName)
	visibility := "Public"
	if m.CurrentRepo.Private {
		visibility = "Private"
	}
	visBadge := statusBadge.Render(visibility)
	lang := safeStr(m.CurrentRepo.Language, "Unknown")
	langBadge := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(lang)
	header := lipgloss.JoinHorizontal(lipgloss.Center, fullName, "  ", visBadge, "  ", langBadge)

	description := descStyle.Width(50).Render(safeStr(m.CurrentRepo.Description, "No description provided."))

	var topicBlock string
	if len(m.CurrentRepo.Topics) > 0 {
		limit := 5
		if len(m.CurrentRepo.Topics) < 5 {
			limit = len(m.CurrentRepo.Topics)
		}
		topicBlock = "\n" + labelStyle.Render("Topics: ") + strings.Join(m.CurrentRepo.Topics[:limit], ", ")
	}

	leftCol := lipgloss.JoinVertical(lipgloss.Left,
		description, topicBlock, "\n",
		labelStyle.Render("Created: ")+formatDate(m.CurrentRepo.CreatedAt),
		labelStyle.Render("Updated: ")+formatDate(m.CurrentRepo.UpdatedAt),
	)
	leftBox := boxStyle.Width(55).Height(12).Render(leftCol)

	makeStat := func(label, icon string, val int) string {
		return fmt.Sprintf("%s %s  %s", icon, statValStyle.Render(fmt.Sprintf("%d", val)), labelStyle.Render(label))
	}
	statsContent := lipgloss.JoinVertical(lipgloss.Left,
		makeStat("Stars", "★", m.CurrentRepo.StargazersCount),
		makeStat("Forks", "⑂", m.CurrentRepo.ForksCount),
		makeStat("Issues", "◎", m.CurrentRepo.OpenIssuesCount),
		makeStat("Watchers", "👁", m.CurrentRepo.WatchersCount),
		"\n",
		labelStyle.Render("Size: ")+fmt.Sprintf("%d KB", m.CurrentRepo.Size),
		labelStyle.Render("Branch: ")+m.CurrentRepo.DefaultBranch,
	)
	rightBox := boxStyle.Width(25).Height(12).Render(statsContent)
	middleSection := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	cloneLink := lipgloss.NewStyle().Foreground(lipgloss.Color("#43BF6D")).Render(m.CurrentRepo.CloneURL)
	footer := lipgloss.JoinVertical(lipgloss.Left, labelStyle.Render("HTTP Clone:"), cloneLink)
	footerBox := boxStyle.Width(82).Render(footer)

	var readmeBlock string
	if m.LoadingReadme {
		readmeBlock = lipgloss.NewStyle().Foreground(subtle).Padding(2).Render("Loading README...")
	} else {
		readmeBlock = fmt.Sprintf("\n%s\n\n%s",
			labelStyle.Render("README.md"),
			m.ReadmeText,
		)
	}
	readmeSection := readmeBoxStyle.Width(82).Render(readmeBlock)

	content := lipgloss.JoinVertical(lipgloss.Center,
		header, "\n", middleSection, footerBox, readmeSection,
		"\n"+lipgloss.NewStyle().Foreground(subtle).Render("(Scroll with j/k • backspace to go back)"),
	)
	return lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, content)
}

func (m RepoPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var err error
	switch msg := msg.(type) {

	case ReadmeMsg:
		m.LoadingReadme = false
		m.ReadmeText, err = glamour.Render(string(msg), "dark")
		if err != nil {
			panic(err)
		}
		m.Viewport.SetContent(m.renderFullPage())

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		m.Viewport.Width = m.Width
		m.Viewport.Height = m.Height

		m.Viewport.Width = m.Width
		m.Viewport.Height = m.Height
		m.Viewport.SetContent(m.renderFullPage())

	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			return m, func() tea.Msg {
				return NavMsg{to: m.CameFrom, from: RepoPage, repodata: m.CurrentRepo, userdata: m.UserData}
			}

		case "j", "down":
			m.Viewport.LineDown(1)
		case "k", "up":
			m.Viewport.LineUp(1)
		case "d", "ctrl+d":
			m.Viewport.HalfViewDown()
		case "u", "ctrl+u":
			m.Viewport.HalfViewUp()
		}
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func safeStr(s *string, fallback string) string {
	if s == nil || *s == "" {
		return fallback
	}
	return *s
}

func formatDate(t time.Time) string {
	return t.Format("02 Jan 2006")
}

func (m RepoPageModel) View() string {
	return m.Viewport.View()
}
