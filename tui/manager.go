package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	p := tea.NewProgram(HomePage{})
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
}
