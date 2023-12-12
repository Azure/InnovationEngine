package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(executeCommand)

	// Bool flags
	executeCommand.PersistentFlags().
		Bool("verbose", false, "Enable verbose logging & standard output.")
	executeCommand.PersistentFlags().
		Bool("do-not-delete", false, "Do not delete the Azure resources created by the Azure CLI commands executed.")

	// String flags
	executeCommand.PersistentFlags().
		String("correlation-id", "", "Adds a correlation ID to the user agent used by a scenarios azure-cli commands.")
	executeCommand.PersistentFlags().
		String("subscription", "", "Sets the subscription ID used by a scenarios azure-cli commands. Will rely on the default subscription if not set.")
	executeCommand.PersistentFlags().
		String("working-directory", ".", "Sets the working directory for innovation engine to operate out of. Restores the current working directory when finished.")

	// StringArray flags
	executeCommand.PersistentFlags().
		StringArray("var", []string{}, "Sets an environment variable for the scenario. Format: --var <key>=<value>")
}

var executeCommand = &cobra.Command{
	Use:   "execute [markdown file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Execute the commands for an Azure deployment scenario.",
	RunE: func(cmd *cobra.Command, args []string) error {
		markdownFile := args[0]
		if markdownFile == "" {
			return errors.New("no markdown file specified")
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		doNotDelete, _ := cmd.Flags().GetBool("do-not-delete")

		subscription, _ := cmd.Flags().GetString("subscription")
		correlationId, _ := cmd.Flags().GetString("correlation-id")
		environment, _ := cmd.Flags().GetString("environment")
		workingDirectory, _ := cmd.Flags().GetString("working-directory")

		environmentVariables, _ := cmd.Flags().GetStringArray("var")
		features, _ := cmd.Flags().GetStringArray("feature")

		// Known features
		renderValues := false

		// Parse the environment variables from the command line into a map
		cliEnvironmentVariables := make(map[string]string)
		for _, environmentVariable := range environmentVariables {
			keyValuePair := strings.SplitN(environmentVariable, "=", 2)
			if len(keyValuePair) != 2 {
				return fmt.Errorf("invalid environment variable format: %s", environmentVariable)
			}

			cliEnvironmentVariables[keyValuePair[0]] = keyValuePair[1]
		}

		for _, feature := range features {
			switch feature {
			case "render-values":
				renderValues = true
			default:
				return fmt.Errorf("invalid feature: %s", feature)
			}
		}

		// Parse the markdown file and create a scenario
		scenario, err := engine.CreateScenarioFromMarkdown(
			markdownFile,
			[]string{"bash", "azurecli", "azurecli-interactive", "terraform"},
			cliEnvironmentVariables,
		)
		if err != nil {
			return fmt.Errorf("creating scenario: %w", err)
		}

		innovationEngine, err := engine.NewEngine(engine.EngineConfiguration{
			Verbose:          verbose,
			DoNotDelete:      doNotDelete,
			Subscription:     subscription,
			CorrelationId:    correlationId,
			Environment:      environment,
			WorkingDirectory: workingDirectory,
			RenderValues:     renderValues,
		})
		if err != nil {
			return fmt.Errorf("creating engine: %w", err)
		}

		// Execute the scenario
		if err = innovationEngine.ExecuteScenario(scenario); err != nil {
			return fmt.Errorf("executing scenario: %w", err)
		}

		return nil
	},
}
