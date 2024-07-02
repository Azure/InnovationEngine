package test

import "github.com/charmbracelet/bubbles/viewport"

// Components used for test mode.
type testModeComponents struct {
	commandViewport viewport.Model
}

// Initializes the viewports for the interactive mode model.
func initializeComponents(model TestModeModel, width, height int) testModeComponents {
	commandViewport := viewport.New(width, height)

	components := testModeComponents{
		commandViewport: commandViewport,
	}

	components.updateViewportSizing(width, height)
	return components
}

// Update the viewport height for the test mode components.
func (components *testModeComponents) updateViewportSizing(terminalWidth int, terminalHeight int) {
	components.commandViewport.Width = terminalWidth
	components.commandViewport.Height = terminalHeight - 1
}
