package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/chirag-diwan/RemGit/config"
	"github.com/chirag-diwan/RemGit/tui"
)

func main() {
	tokens := config.Lexer("/home/chirag/.remgit.conf")
	praser := config.NewPraser(tokens)
	obj := praser.Prase()
	Manager := tui.NewManager(obj)
	p := tea.NewProgram(Manager)
	if _, err := p.Run(); err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
}
