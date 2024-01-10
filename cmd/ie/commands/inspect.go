package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/ui"
	"github.com/spf13/cobra"
)

// Register the command with our command runner.
func init() {
	rootCommand.AddCommand(inspectCommand)

	// String flags
	inspectCommand.PersistentFlags().
		String("correlation-id", "", "Adds a correlation ID to the user agent used by a scenarios azure-cli commands.")
	inspectCommand.PersistentFlags().
		String("subscription", "", "Sets the subscription ID used by a scenarios azure-cli commands. Will rely on the default subscription if not set.")
	inspectCommand.PersistentFlags().
		String("working-directory", ".", "Sets the working directory for innovation engine to operate out of. Restores the current working directory when finished.")

	// StringArray flags
	inspectCommand.PersistentFlags().
		StringArray("var", []string{}, "Sets an environment variable for the scenario. Format: --var <key>=<value>")
}

var inspectCommand = &cobra.Command{
	Use:   "inspect",
	Short: "Execute a document in inspect mode.",
	Run: func(cmd *cobra.Command, args []string) {
		markdownFile := args[0]
		if markdownFile == "" {
			logging.GlobalLogger.Errorf("Error: No markdown file specified.")
			cmd.Help()
			os.Exit(1)
		}

		environmentVariables, _ := cmd.Flags().GetStringArray("var")
		// features, _ := cmd.Flags().GetStringArray("feature")

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
		// Parse the markdown file and create a scenario
		scenario, err := engine.CreateScenarioFromMarkdown(
			markdownFile,
			[]string{"bash", "azurecli", "azurecli-inspect", "terraform"},
			cliEnvironmentVariables,
		)
		if err != nil {
			logging.GlobalLogger.Errorf("Error creating scenario: %s", err)
			fmt.Printf("Error creating scenario: %s", err)
			os.Exit(1)
		}

		if err != nil {
			logging.GlobalLogger.Errorf("Error creating engine: %s", err)
			fmt.Printf("Error creating engine: %s", err)
			os.Exit(1)
		}

		fmt.Println(ui.ScenarioTitleStyle.Render(scenario.Name))
		for stepNumber, step := range scenario.Steps {
			stepTitle := fmt.Sprintf("  %d. %s\n", stepNumber+1, step.Name)
			fmt.Println(ui.StepTitleStyle.Render(stepTitle))
			for codeBlockNumber, codeBlock := range step.CodeBlocks {
				fmt.Println(
					ui.InteractiveModeCodeBlockDescriptionStyle.Render(
						fmt.Sprintf(
							"    %d.%d %s",
							stepNumber+1,
							codeBlockNumber+1,
							codeBlock.Description,
						),
					),
				)
				fmt.Print(
					ui.IndentMultiLineCommand(
						fmt.Sprintf(
							"      %s",
							ui.InteractiveModeCodeBlockStyle.Render(
								codeBlock.Content,
							),
						),
						6),
				)
				fmt.Println()
			}
		}

	},
}
