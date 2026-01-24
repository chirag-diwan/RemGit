package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chirag-diwan/RemGit/tui"
)

func main() {
	p := tea.NewProgram(tui.NewManager())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
}
