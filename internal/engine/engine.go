package engine

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/utils"
	"github.com/charmbracelet/lipgloss"
)

var (
	scriptHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Bold(true)
	scriptText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
)

type EngineConfiguration struct {
	Verbose          bool
	ResourceTracking bool
	DoNotDelete      bool
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

// Executes a deployment scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	fmt.Println(titleStyle.Render(scenario.Name))
	e.ExecuteAndRenderSteps(scenario.Steps, utils.CopyMap(scenario.Environment))
	fmt.Printf(scriptHeader.Render("# Generated bash to replicate what just happened:")+"\n%s", scriptText.Render(scenario.ToShellScript()))
	return nil
}

// Validates a deployment scenario.
func (e *Engine) TestScenario(scenario *Scenario) error {
	fmt.Println(titleStyle.Render(scenario.Name))
	e.TestSteps(scenario.Steps, utils.CopyMap(scenario.Environment))
	return nil
}
