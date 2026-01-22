package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chirag-diwan/RemGit/utils"
)

func main() {
	p := tea.NewProgram(utils.Model{AppState: utils.StateStart, SearchType: utils.SearchRepo})
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
}
