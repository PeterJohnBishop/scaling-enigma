package tui

import (
	"log"
	"os"
	admin "scaling-enigma/go-scaling-enigma/main.go/tui/Admin"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

func StartCLI() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := tea.NewProgram(admin.InitAdminModel())
		if _, err := p.Run(); err != nil {
			log.Fatal("TUI Failed to start: " + err.Error())
			os.Exit(1)
		}
	}()
	wg.Wait()

}
