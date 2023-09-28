package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/lib/fs"
	"github.com/Azure/InnovationEngine/internal/logging"
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
	if configuration.Environment == "ocd" {
		err := az.LoginWithMSI()
		if err != nil {
			logging.GlobalLogger.Errorf("Invalid Config: Failed to login: %s", err)
			return nil, err
		}
	}

	err := az.SetSubscription(configuration.Subscription)
	if err != nil {
		logging.GlobalLogger.Errorf("Invalid Config: Failed to set subscription: %s", err)
		return nil, err
	}

	return &Engine{
		Configuration: configuration,
	}, nil
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

type AzureTokens struct {
	Tokens []string `json:"tokens"`
}

// Executes a deployment scenario.
func (e *Engine) ExecuteScenario(scenario *Scenario) error {
	// Store the current directory so we can restore it later
	originalDirectory, err := os.Getwd()
	if err != nil {
		logging.GlobalLogger.Error("Failed to get current directory", err)
		return err
	}

	fs.SetWorkingDirectory(e.Configuration.WorkingDirectory)
	setCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

	if e.Configuration.Environment == "ocd" {
		tokenURLs := AzureTokens{
			Tokens: []string{},
		}

		for _, step := range scenario.Steps {
			for _, codeblock := range step.CodeBlocks {
				for _, provider := range az.AzureTokenProviders {
					if provider.Regex.MatchString(codeblock.Content) {
						// accessToken, err := az.GetAccessToken(provider)
						// if err != nil {
						// 	logging.GlobalLogger.Errorf("Failed to get access token: %s", err)
						// 	return err
						// }

						tokenURLs.Tokens = append(tokenURLs.Tokens, provider.Resource)

					}
				}
			}
		}

		json, err := json.Marshal(tokenURLs)
		if err != nil {
			logging.GlobalLogger.Errorf("Failed to marshal tokens: %s", err)
			return err
		}
		fmt.Printf("ie_us%sie_ue\n", string(json))
	}

	// Execute the steps
	fmt.Println(scenarioTitleStyle.Render(scenario.Name))
	e.ExecuteAndRenderSteps(scenario.Steps, lib.CopyMap(scenario.Environment))
	fmt.Printf(scriptHeader.Render("# Generated bash to replicate what just happened:")+"\n%s\n", scriptText.Render(scenario.ToShellScript()))

	fs.SetWorkingDirectory(originalDirectory)

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

	fs.SetWorkingDirectory(e.Configuration.WorkingDirectory)
	setCorrelationId(e.Configuration.CorrelationId, scenario.Environment)

	fmt.Println(scenarioTitleStyle.Render(scenario.Name))
	e.TestSteps(scenario.Steps, lib.CopyMap(scenario.Environment))

	fs.SetWorkingDirectory(originalDirectory)
	return nil
}
