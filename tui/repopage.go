package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/githubapi"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	text      = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#FFFFFF"}
	warning   = lipgloss.AdaptiveColor{Light: "#F1F1F1", Dark: "#CD5C5C"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	docStyle  = lipgloss.NewStyle().Padding(1, 2)
	descStyle = lipgloss.NewStyle().Foreground(colText).Italic(true)

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
)

type RepoPageModel struct {
	Width       int
	Height      int
	CurrentRepo githubapi.Repository
}

func (m RepoPageModel) SetRepoData(data githubapi.Repository) RepoPageModel {
	m.CurrentRepo = data
	return m
}

func NewRepoPageModel(data githubapi.Repository) RepoPageModel {
	return RepoPageModel{
		CurrentRepo: data,
	}
}

func (m RepoPageModel) Init() tea.Cmd {
	return nil
}

func (m RepoPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			return m, func() tea.Msg {
				return NavMsg{to: SearchPage}
			}
		}
	}
	return m, nil
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

	leftCol := lipgloss.JoinVertical(
		lipgloss.Left,
		description,
		topicBlock,
		"\n",
		labelStyle.Render("Created: ")+formatDate(m.CurrentRepo.CreatedAt),
		labelStyle.Render("Updated: ")+formatDate(m.CurrentRepo.UpdatedAt),
	)

	leftBox := boxStyle.Width(55).Height(12).Render(leftCol)

	makeStat := func(label, icon string, val int) string {
		return fmt.Sprintf("%s %s  %s", icon, statValStyle.Render(fmt.Sprintf("%d", val)), labelStyle.Render(label))
	}

	statsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		makeStat("Stars", "★", m.CurrentRepo.StargazersCount),
		makeStat("Forks", "⑂", m.CurrentRepo.ForksCount),
		makeStat("Issues", "◎", m.CurrentRepo.OpenIssuesCount),
		makeStat("Watchers", "👁", m.CurrentRepo.WatchersCount),
		"\n",
		labelStyle.Render("Size: ")+fmt.Sprintf("%d KB", m.CurrentRepo.Size),
		labelStyle.Render("Branch: ")+m.CurrentRepo.DefaultBranch,
	)

	rightBox := boxStyle.Width(25).Height(12).Render(statsContent)

	cloneLink := lipgloss.NewStyle().Foreground(lipgloss.Color("#43BF6D")).Render(m.CurrentRepo.CloneURL)
	footer := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("HTTP Clone:"),
		cloneLink,
	)
	footerBox := boxStyle.Width(82).Render(footer)

	middleSection := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	ui := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"\n",
		middleSection,
		footerBox,
		"\n"+lipgloss.NewStyle().Foreground(subtle).Render("(backspace to go back)"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		ui,
	)
}
