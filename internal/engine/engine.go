package engine

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/lib/fs"
	"github.com/Azure/InnovationEngine/internal/logging"
)

const (
	EnvironmentsLocal = "local"
	EnvironmentsCI    = "ci"
	EnvironmentsOCD   = "ocd"
	EnvironmentsAzure = "azure"
)

// Check if the environment is valid.
func IsValidEnvironment(environment string) bool {
	switch environment {
	case EnvironmentsLocal, EnvironmentsCI, EnvironmentsOCD, EnvironmentsAzure:
		return true
	default:
		return false
	}

}

type EngineConfiguration struct {
	Verbose          bool
	DoNotDelete      bool
	CorrelationId    string
	Subscription     string
	Environment      string
	WorkingDirectory string
}

type Engine struct {
	Configuration EngineConfiguration
}

// / Create a new engine instance.
func NewEngine(configuration EngineConfiguration) (*Engine, error) {
	err := az.SetSubscription(configuration.Subscription)
	if err != nil {
		logging.GlobalLogger.Errorf("Invalid Config: Failed to set subscription: %s", err)
		return nil, err
	}

	return &Engine{
		Configuration: configuration,
	}, nil
}

type AzureTokens struct {
	Tokens []string `json:"tokens"`
}

// Executes a deployment scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	return fs.UsingDirectory(e.Configuration.WorkingDirectory, func() error {
		az.SetCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

		// Execute the steps
		fmt.Println(scenarioTitleStyle.Render(scenario.Name))
		err := e.ExecuteAndRenderSteps(scenario.Steps, lib.CopyMap(scenario.Environment))
		return err
	})
}

// Validates a deployment scenario.
func (e *Engine) TestScenario(scenario *Scenario) error {
	return fs.UsingDirectory(e.Configuration.WorkingDirectory, func() error {
		az.SetCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

		// Test the steps
		fmt.Println(scenarioTitleStyle.Render(scenario.Name))
		err := e.TestSteps(scenario.Steps, lib.CopyMap(scenario.Environment))
		return err
	})
}
