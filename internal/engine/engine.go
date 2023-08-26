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

func setSubscription(subscription string) error {
	if subscription != "" {
		command := fmt.Sprintf("az account set --subscription %s", subscription)
		_, err := shells.ExecuteBashCommand(command, shells.BashCommandConfiguration{EnvironmentVariables: map[string]string{}, InteractiveCommand: false, WriteToHistory: false, InheritEnvironment: false})

		if err != nil {
			logging.GlobalLogger.Error("Failed to set subscription", err)
			return err
		}

		logging.GlobalLogger.Infof("Set subscription to %s", subscription)
	}

	return nil
}

func setWorkingDirectory(directory string) error {
	// Change working directory if specified
	if directory != "" {
		err := os.Chdir(directory)
		if err != nil {
			logging.GlobalLogger.Error("Failed to change working directory", err)
			return err
		}

		logging.GlobalLogger.Infof("Changed directory to %s", directory)
	}
	return nil
}

// Executes a deployment scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	err := setSubscription(e.Configuration.Subscription)
	if err != nil {
		return err
	}

	// Store the current directory so we can restore it later
	originalDirectory, err := os.Getwd()
	if err != nil {
		logging.GlobalLogger.Error("Failed to get current directory", err)
		return err
	}

	setWorkingDirectory(e.Configuration.WorkingDirectory)

	// Execute the steps
	fmt.Println(scenarioTitleStyle.Render(scenario.Name))
	e.ExecuteAndRenderSteps(scenario.Steps, utils.CopyMap(scenario.Environment))
	fmt.Printf(scriptHeader.Render("# Generated bash to replicate what just happened:")+"\n%s\n", scriptText.Render(scenario.ToShellScript()))

	setWorkingDirectory(originalDirectory)

	return nil
}

// Validates a deployment scenario.
func (e *Engine) TestScenario(scenario *Scenario) error {
	err := setSubscription(e.Configuration.Subscription)
	if err != nil {
		return err
	}

	fmt.Println(scenarioTitleStyle.Render(scenario.Name))
	e.TestSteps(scenario.Steps, utils.CopyMap(scenario.Environment))
	return nil
}
