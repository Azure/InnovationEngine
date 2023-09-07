package commands

import (
	"fmt"
	"os"

	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(executeCommand)

	// Bool flags
	executeCommand.PersistentFlags().Bool("verbose", false, "Enable verbose logging & standard output.")
	executeCommand.PersistentFlags().Bool("do-not-delete", false, "Do not delete the Azure resources created by the Azure CLI commands executed.")

	// String flags
	executeCommand.PersistentFlags().String("correlation-id", "", "Adds a correlation ID to the user agent used by a scenarios azure-cli commands.")
	executeCommand.PersistentFlags().String("subscription", "", "Sets the subscription ID used by a scenarios azure-cli commands. Will rely on the default subscription if not set.")
	executeCommand.PersistentFlags().String("working-directory", ".", "Sets the working directory for innovation engine to operate out of. Restores the current working directory when finished.")
}

var executeCommand = &cobra.Command{
	Use:   "execute [markdown file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Execute the commands for an Azure deployment scenario.",
	Run: func(cmd *cobra.Command, args []string) {
		markdownFile := args[0]
		if markdownFile == "" {
			cmd.Help()
			os.Exit(1)
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		doNotDelete, _ := cmd.Flags().GetBool("do-not-delete")
		subscription, _ := cmd.Flags().GetString("subscription")
		correlationId, _ := cmd.Flags().GetString("correlation-id")
		environment, _ := cmd.Flags().GetString("environment")
		workingDirectory, _ := cmd.Flags().GetString("working-directory")

		innovationEngine := engine.NewEngine(engine.EngineConfiguration{
			Verbose:          verbose,
			DoNotDelete:      doNotDelete,
			Subscription:     subscription,
			CorrelationId:    correlationId,
			Environment:      environment,
			WorkingDirectory: workingDirectory,
		})

		// Parse the markdown file and create a scenario
		scenario, err := engine.CreateScenarioFromMarkdown(markdownFile, []string{"bash", "azurecli", "azurecli-interactive", "terraform"})
		if err != nil {
			logging.GlobalLogger.Errorf("Error creating scenario: %s", err)
			fmt.Printf("Error creating scenario: %s", err)
			os.Exit(1)
		}

		// Execute the scenario
		err = innovationEngine.ExecuteScenario(scenario)
		if err != nil {
			logging.GlobalLogger.Errorf("Error executing scenario: %s", err)
			fmt.Printf("Error executing scenario: %s", err)
			os.Exit(1)
		}
	},
}
