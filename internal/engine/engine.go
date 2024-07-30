package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/engine/interactive"
	"github.com/Azure/InnovationEngine/internal/engine/test"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/lib/fs"
	"github.com/Azure/InnovationEngine/internal/logging"
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
	GenerateReport   string
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

// Executes a markdown scenario.
func (e *Engine) ExecuteScenario(scenario *common.Scenario) error {
	return fs.UsingDirectory(e.Configuration.WorkingDirectory, func() error {
		az.SetCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

		// Execute the steps
		fmt.Println(ui.ScenarioTitleStyle.Render(scenario.Name))
		err := e.ExecuteAndRenderSteps(scenario.Steps, lib.CopyMap(scenario.Environment))
		return err
	})
}

// Executes a scenario in testing moe. This mode goes over each code block
// and executes it without user interaction.
func (e *Engine) TestScenario(scenario *common.Scenario) error {
	return fs.UsingDirectory(e.Configuration.WorkingDirectory, func() error {
		az.SetCorrelationId(e.Configuration.CorrelationId, scenario.Environment)
		stepsToExecute := filterDeletionCommands(scenario.Steps, e.Configuration.DoNotDelete)

		initialEnvironmentVariables := lib.GetEnvironmentVariables()

		model, err := test.NewTestModeModel(
			scenario.Name,
			e.Configuration.Subscription,
			e.Configuration.Environment,
			stepsToExecute,
			lib.CopyMap(scenario.Environment),
		)
		if err != nil {
			return err
		}

		var flags []tea.ProgramOption
		if environments.EnvironmentsGithubAction == e.Configuration.Environment {
			flags = append(
				flags,
				tea.WithoutRenderer(),
				tea.WithOutput(os.Stdout),
				tea.WithInput(os.Stdin),
			)
		} else {
			flags = append(flags, tea.WithAltScreen(), tea.WithMouseCellMotion())
		}

		common.Program = tea.NewProgram(model, flags...)

		var finalModel tea.Model
		finalModel, err = common.Program.Run()

		// TODO(vmarcella): After testing is complete, we should generate a report.

		model, ok := finalModel.(test.TestModeModel)

		if !ok {
			err = errors.Join(err, fmt.Errorf("failed to cast tea.Model to TestModeModel"))
			return err
		}

		allEnvironmentVariables, err := lib.LoadEnvironmentStateFile(
			lib.DefaultEnvironmentStateFile,
		)
		if err != nil {
			logging.GlobalLogger.Errorf("Failed to load environment state file: %s", err)
			return err
		}

		variablesDeclaredByScenario := lib.DiffMaps(
			allEnvironmentVariables,
			initialEnvironmentVariables,
		)

		if e.Configuration.GenerateReport != "" {
			report := common.BuildReport(scenario.Name)
			report.
				WithProperties(scenario.Properties).
				WithEnvironmentVariables(variablesDeclaredByScenario).
				WriteToJSONFile(e.Configuration.GenerateReport)
		}

		err = errors.Join(err, model.GetFailure())

		fmt.Println(strings.Join(model.CommandLines, "\n"))

		return err
	})
}

// Executes a Scenario in interactive mode. This mode goes over each codeblock
// step by step and allows the user to interact with the codeblock.
func (e *Engine) InteractWithScenario(scenario *common.Scenario) error {
	return fs.UsingDirectory(e.Configuration.WorkingDirectory, func() error {
		az.SetCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

		stepsToExecute := filterDeletionCommands(scenario.Steps, e.Configuration.DoNotDelete)

		model, err := interactive.NewInteractiveModeModel(
			scenario.Name,
			e.Configuration.Subscription,
			e.Configuration.Environment,
			stepsToExecute,
			lib.CopyMap(scenario.Environment),
		)
		if err != nil {
			return err
		}

		common.Program = tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

		var finalModel tea.Model
		var ok bool
		finalModel, err = common.Program.Run()

		model, ok = finalModel.(interactive.InteractiveModeModel)

		if environments.EnvironmentsAzure == e.Configuration.Environment {
			if !ok {
				return fmt.Errorf("failed to cast tea.Model to InteractiveModeModel")
			}

			logging.GlobalLogger.Info("Writing session output to stdout")
			fmt.Println(strings.Join(model.CommandLines, "\n"))
		}

		switch e.Configuration.Environment {
		case environments.EnvironmentsAzure, environments.EnvironmentsOCD:
			logging.GlobalLogger.Info(
				"Cleaning environment variable file located at /tmp/env-vars",
			)
			err := lib.CleanEnvironmentStateFile(lib.DefaultEnvironmentStateFile)
			if err != nil {
				logging.GlobalLogger.Errorf("Error cleaning environment variables: %s", err.Error())
				return err
			}

		default:
			lib.DeleteEnvironmentStateFile(lib.DefaultEnvironmentStateFile)
		}

		if err != nil {
			logging.GlobalLogger.Errorf("Failed to run program %s", err)
			return err
		}

		return nil
	})
}
