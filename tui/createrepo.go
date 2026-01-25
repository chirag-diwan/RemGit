package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chirag-diwan/RemGit/githubapi"
)

const (
	ModeNav int = iota
	ModeEdit
)

const (
	fieldRepoName = iota
	fieldDescription
	fieldPrivate
	fieldIssues
	fieldProjects
	fieldWiki
	fieldSubmit
)

type createResultMsg struct {
	Success bool
}

func createRepoCmd(githubPAT string, req githubapi.RepoRequest) tea.Cmd {
	return func() tea.Msg {
		result := githubapi.CreateRepo(githubPAT, req)
		return createResultMsg{
			Success: result,
		}
	}
}

type createRepoPage struct {
	mode        int
	focusIndex  int
	inputs      []textinput.Model
	toggles     []bool
	statusMsg   string
	statusColor lipgloss.Color

	width     int
	height    int
	githubPAT string
}

func NewCreateRepoPage(githubPAT string, width, height int) tea.Model {
	m := createRepoPage{
		mode:       ModeNav,
		focusIndex: 0,
		toggles:    []bool{false, true, true, true},
		width:      width,
		height:     height,
		githubPAT:  githubPAT,
	}

	m.inputs = make([]textinput.Model, 2)

	tName := textinput.New()
	tName.Placeholder = "Project Name"
	tName.CharLimit = 100
	tName.Width = 50
	tName.Prompt = ""
	m.inputs[fieldRepoName] = tName

	tDesc := textinput.New()
	tDesc.Placeholder = "Short description of your repository"
	tDesc.CharLimit = 200
	tDesc.Width = 50
	tDesc.Prompt = ""
	m.inputs[fieldDescription] = tDesc

	return m
}

func (m createRepoPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m createRepoPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case createResultMsg:
		if msg.Success {
			m.statusMsg = "Repository Created Successfully!"
			m.statusColor = special
		} else {
			m.statusMsg = "Failed to Create Repository."
			m.statusColor = warning
		}
		m.mode = ModeNav

	case tea.KeyMsg:

		if m.mode == ModeEdit {
			switch msg.String() {
			case "esc":
				m.mode = ModeNav
				m.inputs[m.focusIndex].Blur()
				return m, nil
			case "enter":
				m.mode = ModeNav
				m.inputs[m.focusIndex].Blur()
				return m, nil
			}

			return m, m.updateInputs(msg)
		}

		switch msg.String() {
		case "q":

			return m, func() tea.Msg {
				return NavMsg{to: SearchPage, from: CreateRepoPage}
			}

		case "backspace":
			return m, nil

		case "esc":
			return m, func() tea.Msg {
				return NavMsg{to: SearchPage, from: CreateRepoPage}
			}

		case "tab", "shift+tab", "up", "down", "k", "j":
			s := msg.String()

			if s == "up" || s == "shift+tab" || s == "k" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > fieldSubmit {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = fieldSubmit
			}

		case "enter", " ":

			if m.focusIndex == fieldRepoName || m.focusIndex == fieldDescription {
				m.mode = ModeEdit
				m.inputs[m.focusIndex].Focus()

				cmd := m.inputs[m.focusIndex].Cursor.BlinkCmd()
				return m, cmd
			}

			if m.focusIndex >= fieldPrivate && m.focusIndex <= fieldWiki {
				idx := m.focusIndex - 2
				m.toggles[idx] = !m.toggles[idx]
			}

			if m.focusIndex == fieldSubmit {
				m.statusMsg = "Creating repository..."
				m.statusColor = text

				req := githubapi.RepoRequest{
					Name:        m.inputs[fieldRepoName].Value(),
					Description: m.inputs[fieldDescription].Value(),
					Private:     m.toggles[fieldPrivate-2],
					HasIssues:   m.toggles[fieldIssues-2],
					HasProjects: m.toggles[fieldProjects-2],
					HasWiki:     m.toggles[fieldWiki-2],
				}
				return m, createRepoCmd(m.githubPAT, req)
			}
		}
	}

	for i := range m.inputs {
		if i == m.focusIndex {

			if m.mode == ModeEdit {
				m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(text)
				m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(highlight)
			} else {
				m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(highlight)
				m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(highlight)
			}
		} else {
			m.inputs[i].Blur()
			m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(subtle)
			m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(subtle)
		}
	}

	return m, nil
}

func (m *createRepoPage) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m createRepoPage) View() string {

	activeStyle := lipgloss.NewStyle().Foreground(highlight)
	inactiveStyle := lipgloss.NewStyle().Foreground(subtle)
	checkboxStyle := lipgloss.NewStyle().Foreground(special)

	modeStr := "NAV"
	modeColor := subtle
	if m.mode == ModeEdit {
		modeStr = "EDIT"
		modeColor = special
	}
	modeIndicator := lipgloss.NewStyle().Foreground(modeColor).Bold(true).Render("-- " + modeStr + " --")

	renderCheckbox := func(label string, isChecked, isFocused bool) string {
		check := "[ ]"
		if isChecked {
			check = "[x]"
		}

		style := inactiveStyle
		if isFocused {
			style = activeStyle
		}

		checkDisplay := style.Render(check)
		if isChecked {
			checkDisplay = checkboxStyle.Render(check)
		}

		return fmt.Sprintf("%s %s", checkDisplay, style.Render(label))
	}

	renderInput := func(label string, input textinput.Model, isFocused bool) string {
		style := inactiveStyle
		if isFocused {
			style = activeStyle
		}

		border := ""
		if isFocused && m.mode == ModeEdit {
			border = " ✐"
		}

		return lipgloss.JoinVertical(lipgloss.Left,
			style.Render(label+border),
			input.View(),
		)
	}

	nameField := renderInput("Repository Name", m.inputs[fieldRepoName], m.focusIndex == fieldRepoName)
	descField := renderInput("Description", m.inputs[fieldDescription], m.focusIndex == fieldDescription)

	settingPrivate := renderCheckbox("Private Repository", m.toggles[0], m.focusIndex == fieldPrivate)
	settingIssues := renderCheckbox("Enable Issues", m.toggles[1], m.focusIndex == fieldIssues)
	settingProjects := renderCheckbox("Enable Projects", m.toggles[2], m.focusIndex == fieldProjects)
	settingWiki := renderCheckbox("Enable Wiki", m.toggles[3], m.focusIndex == fieldWiki)

	submitBtn := "[ Create Repository ]"
	if m.focusIndex == fieldSubmit {
		submitBtn = activeStyle.Copy().Bold(true).Render(submitBtn)
	} else {
		submitBtn = inactiveStyle.Render(submitBtn)
	}

	statusDisplay := ""
	if m.statusMsg != "" {
		statusDisplay = lipgloss.NewStyle().
			Foreground(m.statusColor).
			Bold(true).
			Render(m.statusMsg)
	}

	formContent := lipgloss.JoinVertical(lipgloss.Left,
		heading.Render("Create New Repository"),
		"",
		modeIndicator,
		"",
		nameField,
		"",
		descField,
		"",
		settingPrivate,
		settingIssues,
		settingProjects,
		settingWiki,
		"",
		submitBtn,
		"",
		statusDisplay,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		window.Render(formContent),
	)
}
