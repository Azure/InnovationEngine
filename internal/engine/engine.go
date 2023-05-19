package engine

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/shells"
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
	fmt.Printf("---Generated script---\n%s", scenario.ToShellScript())
	return nil
}
