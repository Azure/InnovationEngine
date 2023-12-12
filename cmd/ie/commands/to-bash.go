package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/spf13/cobra"
)

type AzureScript struct {
	Script string `json:"script"`
}

var toBashCommand = &cobra.Command{
	Use:   "to-bash",
	Short: "Convert a markdown scenario into a bash script.",
	RunE: func(cmd *cobra.Command, args []string) error {
		markdownFile := args[0]
		if markdownFile == "" {
			return errors.New("no markdown file specified")
		}

		environment, _ := cmd.Flags().GetString("environment")
		environmentVariables, _ := cmd.Flags().GetStringArray("var")

		// Parse the environment variables
		cliEnvironmentVariables := make(map[string]string)
		for _, environmentVariable := range environmentVariables {
			keyValuePair := strings.SplitN(environmentVariable, "=", 2)
			if len(keyValuePair) != 2 {
				return fmt.Errorf("invalid environment variable format: %s", environmentVariable)
			}

			cliEnvironmentVariables[keyValuePair[0]] = keyValuePair[1]
		}

		// Parse the markdown file and create a scenario
		scenario, err := engine.CreateScenarioFromMarkdown(
			markdownFile,
			[]string{"bash", "azurecli", "azurecli-interactive", "terraform"},
			cliEnvironmentVariables)
		if err != nil {
			return fmt.Errorf("creating scenario: %w", err)
		}

		// If within cloudshell, we need to wrap the script in a json object to
		// communicate it to the portal.
		if environments.IsAzureEnvironment(environment) {
			script := AzureScript{Script: scenario.ToShellScript()}
			scriptJson, err := json.Marshal(script)
			if err != nil {
				return fmt.Errorf("converting to json: %w", err)
			}

			fmt.Printf("ie_us%sie_ue\n", scriptJson)
		} else {
			fmt.Printf("%s", scenario.ToShellScript())
		}

		return nil
	},
}

func init() {
	rootCommand.AddCommand(toBashCommand)
	toBashCommand.PersistentFlags().
		StringArray("var", []string{}, "Sets an environment variable for the scenario. Format: --var <key>=<value>")
}
