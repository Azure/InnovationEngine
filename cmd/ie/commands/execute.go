package commands

import (
	"github.com/Azure/InnovationEngine/internal/engine"
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
}

var executeCommand = &cobra.Command{
	Use:   "execute [markdown file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Execute the commands for an Azure deployment scenario.",
	Run: func(cmd *cobra.Command, args []string) {
		markdownFile := args[0]
		if markdownFile == "" {
			cmd.Help()
			return
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		do_not_delete, _ := cmd.Flags().GetBool("do-not-delete")
		subscription, _ := cmd.Flags().GetString("subscription")
		correlation_id, _ := cmd.Flags().GetString("correlation-id")

		innovationEngine := engine.NewEngine(engine.EngineConfiguration{
			Verbose:       verbose,
			DoNotDelete:   do_not_delete,
			Subscription:  subscription,
			CorrelationId: correlation_id,
		})

		scenario, err := engine.CreateScenarioFromMarkdown(markdownFile, []string{"bash", "azurecli", "azurecli-interactive", "terraform"})
		if err != nil {
			panic(err)
		}

		innovationEngine.ExecuteScenario(scenario)
	},
}
