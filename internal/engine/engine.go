package engine

import (
	"fmt"
	"os"

	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/Azure/InnovationEngine/internal/utils"
	"github.com/charmbracelet/lipgloss"
)

var (
	scriptHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Bold(true)
	scriptText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
)

const (
	EnvironmentsLocal = "local"
	EnvironmentsCI    = "ci"
	EnvironmentsOCD   = "ocd"
)

func IsValidEnvironment(environment string) bool {
	return environment == EnvironmentsLocal || environment == EnvironmentsCI || environment == EnvironmentsOCD
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
func NewEngine(configuration EngineConfiguration) *Engine {
	return &Engine{
		Configuration: configuration,
	}
}

// Executes a deployment scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	if e.Configuration.Subscription != "" {
		command := fmt.Sprintf("az account set --subscription %s", e.Configuration.Subscription)
		_, err := shells.ExecuteBashCommand(command, map[string]string{}, true, false)
		if err != nil {
			logging.GlobalLogger.Error("Failed to set subscription", err)
			return err
		}
		logging.GlobalLogger.Infof("Set subscription to %s", e.Configuration.Subscription)
	}

	// Store the current directory so we can restore it later
	currentDirectory, err := os.Getwd()
	if err != nil {
		logging.GlobalLogger.Error("Failed to get current directory", err)
		return err
	}

	// Change working directory if specified
	if e.Configuration.WorkingDirectory != "" {
		err = os.Chdir(e.Configuration.WorkingDirectory)
		if err != nil {
			logging.GlobalLogger.Error("Failed to change working directory", err)
			return err
		}

		logging.GlobalLogger.Infof("Changed working directory to %s", e.Configuration.WorkingDirectory)
	}

	// Execute the steps
	fmt.Println(scenarioTitleStyle.Render(scenario.Name))
	e.ExecuteAndRenderSteps(scenario.Steps, utils.CopyMap(scenario.Environment))
	fmt.Printf(scriptHeader.Render("# Generated bash to replicate what just happened:")+"\n%s\n", scriptText.Render(scenario.ToShellScript()))

	// Restore original directory
	err = os.Chdir(currentDirectory)
	if err != nil {
		logging.GlobalLogger.Error("Failed to restore the original working directory", err)
		return err
	}

	return nil
}

// Validates a deployment scenario.
func (e *Engine) TestScenario(scenario *Scenario) error {
	if e.Configuration.Subscription != "" {
		command := fmt.Sprintf("az account set --subscription %s", e.Configuration.Subscription)
		_, err := shells.ExecuteBashCommand(command, map[string]string{}, true, false)
		if err != nil {
			logging.GlobalLogger.Error("Failed to set subscription", err)
			return err
		}
		logging.GlobalLogger.Infof("Set subscription to %s", e.Configuration.Subscription)
	}

	fmt.Println(scenarioTitleStyle.Render(scenario.Name))
	e.TestSteps(scenario.Steps, utils.CopyMap(scenario.Environment))
	return nil
}
