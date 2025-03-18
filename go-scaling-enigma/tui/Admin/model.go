package admin

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Admin struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func InitAdminModel() Admin {
	return Admin{
		cursor:   0,
		choices:  []string{"All Users", "Users Online", "Users Offline"},
		selected: make(map[int]struct{}),
	}
}

func (m Admin) Init() tea.Cmd {
	return nil
}
