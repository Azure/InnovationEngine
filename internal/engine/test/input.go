package test

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Handle user input for Test mode.
func handleUserInput(
	model TestModeModel,
	message tea.KeyMsg,
) (TestModeModel, []tea.Cmd) {
	var commands []tea.Cmd

	switch {
	case key.Matches(message, model.commands.quit):
		commands = append(commands, tea.Quit)
	}

	return model, commands
}
