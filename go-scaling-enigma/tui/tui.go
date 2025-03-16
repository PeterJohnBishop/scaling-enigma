package tui

import (
	"log"
	"os"

	"scaling-enigma/go-scaling-enigma/main.go/websocket"

	tea "github.com/charmbracelet/bubbletea"
)

func StartCLI() {
	p := tea.NewProgram(initialFiller())
	websocket.ServeWebsocketClient()
	if _, err := p.Run(); err != nil {
		log.Fatal("TUI Failed to start: " + err.Error())
		os.Exit(1)

	}
}

type Filler struct {
	message string
}

func initialFiller() Filler {
	return Filler{
		message: "Hi",
	}
}

func (m Filler) Init() tea.Cmd {
	return nil
}

func (m Filler) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Filler) View() string {
	return "Hella World"
}
