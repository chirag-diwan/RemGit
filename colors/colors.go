package colors

import "github.com/charmbracelet/lipgloss"

const (
	Accent = lipgloss.Color("#7aa2f7")
	Text   = lipgloss.Color("#c0caf5")
	Muted  = lipgloss.Color("#565f89")
	Danger = lipgloss.Color("#f7768e")
)

func SetColors() {
	lipgloss.SetColorProfile(lipgloss.ColorProfile())
}
