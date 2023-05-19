package engine

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/charmbracelet/lipgloss"
)

var (
	scriptHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Bold(true)
	scriptText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
)

type EngineConfiguration struct {
	Verbose bool
}

type Engine struct {
	Configuration EngineConfiguration
}

// / Create a new engine instance.
func NewEngine(configuration EngineConfiguration) *Engine {
	return &Engine{
		Configuration: configuration,
	}
}

// / Executes a scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	fmt.Println(titleStyle.Render(scenario.Name))
	ExecuteAndRenderSteps(scenario.Steps, scenario.Environment, e.Configuration.Verbose)
	shells.ResetStoredEnvironmentVariables()
	fmt.Printf(scriptHeader.Render("# Generated bash replicate the deployment:")+"\n%s", scriptText.Render(scenario.ToShellScript()))
	return nil
}
