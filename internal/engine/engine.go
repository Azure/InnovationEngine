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
func NewEngine(configuration EngineConfiguration) (*Engine, error) {
	// Temporarily disable until login code worked out.
	// err := refreshAccessToken()
	// if err != nil {
	// logging.GlobalLogger.Errorf("Invalid Config: Failed to login: %s", err)
	// return nil, err
	// }

	err := setSubscription(configuration.Subscription)
	if err != nil {
		logging.GlobalLogger.Errorf("Invalid Config: Failed to set subscription: %s", err)
		return nil, err
	}

	return &Engine{
		Configuration: configuration,
	}, nil
}

func refreshAccessToken() error {
	// Login
	command := "az account get-access-token > ~/.azure/accessTokens.json"
	logging.GlobalLogger.Info("Logging into the azure cli.")
	output, err := shells.ExecuteBashCommand(command, shells.BashCommandConfiguration{EnvironmentVariables: map[string]string{}, InteractiveCommand: true, WriteToHistory: false, InheritEnvironment: false})

	logging.GlobalLogger.Debugf("Login stdout: %s", output.StdOut)
	logging.GlobalLogger.Debugf("Login stderr: %s", output.StdErr)

	if err != nil {
		logging.GlobalLogger.Errorf("Failed to login %s", err)
		return err
	}

	logging.GlobalLogger.Info("Login successful.")
	return nil
}

func setSubscription(subscription string) error {
	if subscription != "" {
		command := fmt.Sprintf("az account set --subscription %s", subscription)
		_, err := shells.ExecuteBashCommand(command, shells.BashCommandConfiguration{EnvironmentVariables: map[string]string{}, InteractiveCommand: false, WriteToHistory: false, InheritEnvironment: false})

		if err != nil {
			logging.GlobalLogger.Errorf("Failed to set subscription: %s", err)
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

// If the correlation ID is set, we need to set the AZURE_HTTP_USER_AGENT
// environment variable so that the Azure CLI will send the correlation ID
// with Azure Resource Manager requests.
func setCorrelationId(correlationId string, env map[string]string) {
	if correlationId != "" {
		env["AZURE_HTTP_USER_AGENT"] = fmt.Sprintf("innovation-engine-%s", correlationId)
		logging.GlobalLogger.Info("Resource tracking enabled. Tracking ID: " + env["AZURE_HTTP_USER_AGENT"])
	}
}

// Executes a deployment scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	// Store the current directory so we can restore it later
	originalDirectory, err := os.Getwd()
	if err != nil {
		logging.GlobalLogger.Error("Failed to get current directory", err)
		return err
	}

	setWorkingDirectory(e.Configuration.WorkingDirectory)
	setCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

	// Execute the steps
	fmt.Println(scenarioTitleStyle.Render(scenario.Name))
	e.ExecuteAndRenderSteps(scenario.Steps, utils.CopyMap(scenario.Environment))
	fmt.Printf(scriptHeader.Render("# Generated bash to replicate what just happened:")+"\n%s\n", scriptText.Render(scenario.ToShellScript()))

	setWorkingDirectory(originalDirectory)

	return nil
}

// Validates a deployment scenario.
func (e *Engine) TestScenario(scenario *Scenario) error {
	// Store the current directory so we can restore it later
	originalDirectory, err := os.Getwd()
	if err != nil {
		logging.GlobalLogger.Error("Failed to get current directory", err)
		return err
	}

	setWorkingDirectory(e.Configuration.WorkingDirectory)
	setCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

	fmt.Println(scenarioTitleStyle.Render(scenario.Name))
	e.TestSteps(scenario.Steps, utils.CopyMap(scenario.Environment))

	setWorkingDirectory(originalDirectory)
	return nil
}
