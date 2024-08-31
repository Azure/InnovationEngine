package interactive

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type interactiveModeComponents struct {
	paginator        paginator.Model
	stepViewport     viewport.Model
	outputViewport   viewport.Model
	azureCLIViewport viewport.Model
}

// Initializes the viewports for the interactive mode model.
func initializeComponents(model InteractiveModeModel, width, height int) interactiveModeComponents {
	// paginator setup
	p := paginator.New()
	p.TotalPages = len(model.codeBlockState)
	p.Type = paginator.Dots
	// Dots
	p.ActiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).
		Render("•")
	p.InactiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).
		Render("•")

	p.KeyMap.PrevPage = model.commands.previous
	p.KeyMap.NextPage = model.commands.next

	stepViewport := viewport.New(width, 4)
	outputViewport := viewport.New(width, 2)
	azureCLIViewport := viewport.New(width, height)

	components := interactiveModeComponents{
		paginator:        p,
		stepViewport:     stepViewport,
		outputViewport:   outputViewport,
		azureCLIViewport: azureCLIViewport,
	}

	components.updateViewportHeight(height)
	return components
}

func (components *interactiveModeComponents) updateViewportHeight(terminalHeight int) {
	stepViewportPercent := 0.4
	outputViewportPercent := 0.2
	stepViewportHeight := int(float64(terminalHeight) * stepViewportPercent)
	outputViewportHeight := int(float64(terminalHeight) * outputViewportPercent)

	if stepViewportHeight < 4 {
		stepViewportHeight = 4
	}

	if outputViewportHeight < 2 {
		outputViewportHeight = 2
	}

	components.stepViewport.Height = stepViewportHeight
	components.outputViewport.Height = outputViewportHeight
	components.azureCLIViewport.Height = terminalHeight - 1
}

func updateComponents(
	components interactiveModeComponents,
	currentBlock int,
	message tea.Msg,
) (interactiveModeComponents, []tea.Cmd) {
	var commands []tea.Cmd
	var command tea.Cmd

	components.paginator.Page = currentBlock

	components.stepViewport, command = components.stepViewport.Update(message)
	commands = append(commands, command)

	components.outputViewport, command = components.outputViewport.Update(message)
	commands = append(commands, command)

	components.azureCLIViewport, command = components.azureCLIViewport.Update(message)
	commands = append(commands, command)

	return components, commands
}
