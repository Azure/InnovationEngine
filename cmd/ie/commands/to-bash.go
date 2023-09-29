package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/spf13/cobra"
)

type OcdScript struct {
	Script string `json:"script"`
}

var toBashCommand = &cobra.Command{
	Use:   "to-bash",
	Short: "Convert a markdown scenario into a bash script.",
	Run: func(cmd *cobra.Command, args []string) {
		markdownFile := args[0]
		if markdownFile == "" {
			logging.GlobalLogger.Errorf("Error: No markdown file specified.")
			cmd.Help()
			os.Exit(1)
		}

		environment, _ := cmd.Flags().GetString("environment")
		environmentVariables, _ := cmd.Flags().GetStringArray("var")

		// Parse the environment variables
		cliEnvironmentVariables := make(map[string]string)
		for _, environmentVariable := range environmentVariables {
			keyValuePair := strings.SplitN(environmentVariable, "=", 2)
			if len(keyValuePair) != 2 {
				logging.GlobalLogger.Errorf("Error: Invalid environment variable format: %s", environmentVariable)
				fmt.Printf("Error: Invalid environment variable format: %s", environmentVariable)
				cmd.Help()
				os.Exit(1)
			}

			cliEnvironmentVariables[keyValuePair[0]] = keyValuePair[1]
		}

		// Parse the markdown file and create a scenario
		scenario, err := engine.CreateScenarioFromMarkdown(
			markdownFile,
			[]string{"bash", "azurecli", "azurecli-interactive", "terraform"},
			cliEnvironmentVariables)

		if err != nil {
			logging.GlobalLogger.Errorf("Error creating scenario: %s", err)
			fmt.Printf("Error creating scenario: %s", err)
			os.Exit(0)
		}

		if environment == "ocd" {
			script := OcdScript{Script: scenario.ToShellScript()}
			scriptJson, err := json.Marshal(script)

			if err != nil {
				logging.GlobalLogger.Errorf("Error converting to json: %s", err)
				fmt.Printf("Error converting to json: %s", err)
				os.Exit(1)
			}

			fmt.Printf("ie_us%sie_ue\n", scriptJson)
		} else {
			fmt.Printf("%s", scenario.ToShellScript())
		}

	},
}

func init() {
	rootCommand.AddCommand(toBashCommand)
	toBashCommand.PersistentFlags().StringArray("var", []string{}, "Sets an environment variable for the scenario. Format: --var <key>=<value>")
}
