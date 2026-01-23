package tui

import tea "github.com/charmbracelet/bubbletea"

func Run() {
	p := tea.NewProgram(NewSearchPage())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
