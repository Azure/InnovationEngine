package engine

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/lib/fs"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/Azure/InnovationEngine/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// Configuration for the engine.
type EngineConfiguration struct {
	Verbose          bool
	DoNotDelete      bool
	CorrelationId    string
	Subscription     string
	Environment      string
	WorkingDirectory string
	RenderValues     bool
}

type Engine struct {
	Configuration EngineConfiguration
}

// / Create a new engine instance.
func NewEngine(configuration EngineConfiguration) (*Engine, error) {
	return &Engine{
		Configuration: configuration,
	}, nil
}

// Executes a deployment scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	return fs.UsingDirectory(e.Configuration.WorkingDirectory, func() error {
		az.SetCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

		// Execute the steps
		fmt.Println(ui.ScenarioTitleStyle.Render(scenario.Name))
		err := e.ExecuteAndRenderSteps(scenario.Steps, lib.CopyMap(scenario.Environment))
		return err
	})
}

// Validates a deployment scenario.
func (e *Engine) TestScenario(scenario *Scenario) error {
	return fs.UsingDirectory(e.Configuration.WorkingDirectory, func() error {
		az.SetCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

		// Test the steps
		fmt.Println(ui.ScenarioTitleStyle.Render(scenario.Name))
		err := e.TestSteps(scenario.Steps, lib.CopyMap(scenario.Environment))
		return err
	})
}

// Executes a Scenario in interactive mode. This mode goes over each codeblock
// step by step and allows the user to interact with the codeblock.
func (e *Engine) InteractWithScenario(scenario *Scenario) error {
	return fs.UsingDirectory(e.Configuration.WorkingDirectory, func() error {
		az.SetCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

		stepsToExecute := filterDeletionCommands(scenario.Steps, e.Configuration.DoNotDelete)

		model, err := NewInteractiveModeModel(e, stepsToExecute, lib.CopyMap(scenario.Environment))
		if err != nil {
			return err
		}

		program = tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
		_, err = program.Run()

		switch e.Configuration.Environment {
		case environments.EnvironmentsAzure, environments.EnvironmentsOCD:
			logging.GlobalLogger.Info(
				"Cleaning environment variable file located at /tmp/env-vars",
			)
			err := shells.CleanEnvironmentStateFile()
			if err != nil {
				logging.GlobalLogger.Errorf("Error cleaning environment variables: %s", err.Error())
				return err
			}

		default:
			shells.ResetStoredEnvironmentVariables()
		}

		if err != nil {
			logging.GlobalLogger.Errorf("Failed to run program %s", err)
			return err
		}

		return nil
	})
}
