package az

import (
	"fmt"
	"regexp"

	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/shells"
)

func LoginWithMSI() error {
	// Login
	command := "az login --identity"
	logging.GlobalLogger.Info("Logging into the azure cli.")
	output, err := shells.ExecuteBashCommand(command, shells.BashCommandConfiguration{EnvironmentVariables: map[string]string{}, InteractiveCommand: false, WriteToHistory: false, InheritEnvironment: true})

	logging.GlobalLogger.Debugf("Login stdout: %s", output.StdOut)
	logging.GlobalLogger.Debugf("Login stderr: %s", output.StdErr)

	if err != nil {
		logging.GlobalLogger.Errorf("Failed to login %s", err)
		return err
	}

	logging.GlobalLogger.Info("Login successful.")
	return nil
}

func SetSubscription(subscription string) error {
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

type AzureTokenProvider struct {
	Resource string
	Regex    *regexp.Regexp
}

var KeyVaultProvider = AzureTokenProvider{
	Resource: "https://vault.azure.net",
	Regex:    regexp.MustCompile("az keyvault"),
}

var AzureTokenProviders = []AzureTokenProvider{
	KeyVaultProvider,
}

func GetAccessToken(provider AzureTokenProvider) (string, error) {
	command := fmt.Sprintf("az account get-access-token --resource %s --query accessToken -o tsv", provider.Resource)
	output, err := shells.ExecuteBashCommand(command, shells.BashCommandConfiguration{EnvironmentVariables: map[string]string{}, InteractiveCommand: false, WriteToHistory: false, InheritEnvironment: true})

	if err != nil {
		logging.GlobalLogger.Errorf("Failed to get access token: %s", err)
		return "", err
	}

	return output.StdOut, nil
}
