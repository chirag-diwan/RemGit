package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/githubapi"
)

type UserReposMsg struct {
	Login string
	Repos []githubapi.Repository
}

var (
	styleTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(0, 1)

	styleRepoCard = lipgloss.NewStyle().
			PaddingLeft(2).
			Border(lipgloss.HiddenBorder(), false, false, false, true)

	styleRepoActive = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("212")).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("212"))

	styleMeta = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type UserPageModel struct {
	Width           int
	Height          int
	currentUserData githubapi.UserSummary

	repos       []githubapi.Repository
	cursor      int
	windowStart int
	loading     bool
	spinner     spinner.Model
}

func NewUserPageModel(data githubapi.UserSummary) UserPageModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return UserPageModel{
		currentUserData: data,
		repos:           []githubapi.Repository{},
		loading:         true,
		spinner:         s,
	}
}

func (model UserPageModel) SetUserData(data githubapi.UserSummary) UserPageModel {
	model.currentUserData = data
	model.repos = []githubapi.Repository{}
	model.cursor = 0
	model.windowStart = 0
	model.loading = true
	return model
}

func fetchReposCmd(username, url string) tea.Cmd {
	return func() tea.Msg {

		data := githubapi.GetRepoFromUrl(url)

		return UserReposMsg{Login: username, Repos: data}
	}
}

func (model UserPageModel) Init() tea.Cmd {

	return tea.Batch(
		model.spinner.Tick,
		fetchReposCmd(model.currentUserData.Login, model.currentUserData.ReposURL),
	)
}

func (model UserPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		if model.loading {
			model.spinner, cmd = model.spinner.Update(msg)
			return model, cmd
		}

	case UserReposMsg:

		if msg.Login != model.currentUserData.Login {
			return model, nil
		}
		model.loading = false
		model.repos = msg.Repos
		return model, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if len(model.repos) > 0 && model.cursor < len(model.repos)-1 {
				model.cursor++

				if model.cursor >= model.windowStart+5 {
					model.windowStart++
				}
			}
		case "k", "up":
			if model.cursor > 0 {
				model.cursor--
				if model.cursor < model.windowStart {
					model.windowStart--
				}
			}
		}
	}
	return model, nil
}

func (model UserPageModel) View() string {

	header := lipgloss.JoinVertical(
		lipgloss.Center,
		styleTitle.Render(model.currentUserData.Login),
		styleMeta.Render(model.currentUserData.HTMLURL),
		" ",
	)

	var content string

	if model.loading {
		content = fmt.Sprintf("%s Loading repositories...", model.spinner.View())
	} else if len(model.repos) == 0 {
		content = styleMeta.Render("No public repositories found.")
	} else {

		var listItems []string
		itemsPerPage := 5

		endIndex := model.windowStart + itemsPerPage
		if endIndex > len(model.repos) {
			endIndex = len(model.repos)
		}

		for i := model.windowStart; i < endIndex; i++ {
			repo := model.repos[i]

			isActive := i == model.cursor
			style := styleRepoCard
			pointer := "  "
			if isActive {
				style = styleRepoActive
				pointer = "> "
			}

			desc := "No description"
			if repo.Description != nil {
				if len(*repo.Description) > 50 {
					desc = (*repo.Description)[:47] + "..."
				} else {
					desc = *repo.Description
				}
			}

			itemBlock := lipgloss.JoinVertical(
				lipgloss.Left,
				fmt.Sprintf("%s%s", pointer, lipgloss.NewStyle().Bold(true).Render(repo.Name)),
				styleMeta.PaddingLeft(2).Render(desc),
				styleMeta.PaddingLeft(2).Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("★ %d", repo.StargazersCount)),
			)

			listItems = append(listItems, style.Render(itemBlock))
			listItems = append(listItems, "")
		}
		content = lipgloss.JoinVertical(lipgloss.Left, listItems...)
	}

	layout := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		content,
	)

	return lipgloss.Place(
		model.Width,
		model.Height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.NewStyle().PaddingTop(2).Render(layout),
	)
}
