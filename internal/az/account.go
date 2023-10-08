package az

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/shells"
)

func SetSubscription(subscription string) error {
	if subscription != "" {
		command := fmt.Sprintf("az account set --subscription %s", subscription)
		_, err := shells.ExecuteBashCommand(
			command,
			shells.BashCommandConfiguration{
				EnvironmentVariables: map[string]string{},
				InteractiveCommand:   false,
				WriteToHistory:       false,
				InheritEnvironment:   false,
			},
		)

		if err != nil {
			logging.GlobalLogger.Errorf("Failed to set subscription: %s", err)
			return err
		}

		logging.GlobalLogger.Infof("Set subscription to %s", subscription)
	}

	return nil
}
