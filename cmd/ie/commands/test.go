package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(testCommand)
	testCommand.PersistentFlags().
		Bool("verbose", false, "Enable verbose logging & standard output.")
	testCommand.PersistentFlags().
		String("subscription", "", "Sets the subscription ID used by a scenarios azure-cli commands. Will rely on the default subscription if not set.")
	testCommand.PersistentFlags().
		String("working-directory", ".", "Sets the working directory for innovation engine to operate out of. Restores the current working directory when finished.")

	testCommand.PersistentFlags().
		StringArray("var", []string{}, "Sets an environment variable for the scenario. Format: --var <key>=<value>")
}

var testCommand = &cobra.Command{
	Use:   "test",
	Args:  cobra.MinimumNArgs(1),
	Short: "Test document commands against their expected outputs.",
	Run: func(cmd *cobra.Command, args []string) {
		markdownFile := args[0]
		if markdownFile == "" {
			cmd.Help()
			return
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		subscription, _ := cmd.Flags().GetString("subscription")
		workingDirectory, _ := cmd.Flags().GetString("working-directory")
		environment, _ := cmd.Flags().GetString("environment")

		environmentVariables, _ := cmd.Flags().GetStringArray("var")

		// Parse the environment variables from the command line into a map
		cliEnvironmentVariables := make(map[string]string)
		for _, environmentVariable := range environmentVariables {
			keyValuePair := strings.SplitN(environmentVariable, "=", 2)
			if len(keyValuePair) != 2 {
				logging.GlobalLogger.Errorf(
					"Error: Invalid environment variable format: %s",
					environmentVariable,
				)
				fmt.Printf("Error: Invalid environment variable format: %s", environmentVariable)
				cmd.Help()
				os.Exit(1)
			}

			cliEnvironmentVariables[keyValuePair[0]] = keyValuePair[1]
		}

		innovationEngine, err := engine.NewEngine(engine.EngineConfiguration{
			Verbose:          verbose,
			DoNotDelete:      false,
			Subscription:     subscription,
			CorrelationId:    "",
			WorkingDirectory: workingDirectory,
			Environment:      environment,
		})
		if err != nil {
			logging.GlobalLogger.Errorf("Error creating engine %s", err)
			fmt.Printf("Error creating engine %s", err)
			os.Exit(1)
		}

		scenario, err := common.CreateScenarioFromMarkdown(
			markdownFile,
			[]string{"bash", "azurecli", "azurecli-interactive", "terraform"},
			cliEnvironmentVariables,
		)
		if err != nil {
			logging.GlobalLogger.Errorf("Error creating scenario %s", err)
			fmt.Printf("Error creating engine %s", err)
			os.Exit(1)
		}

		err = innovationEngine.TestScenario(scenario)
		if err != nil {
			logging.GlobalLogger.Errorf("Error testing scenario: %s", err)
			fmt.Printf("Error testing scenario: %s\n", err)
			os.Exit(1)
		}
	},
}
