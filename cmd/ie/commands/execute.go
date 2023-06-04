package commands

import (
	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(executeCommand)

	executeCommand.PersistentFlags().Bool("verbose", false, "Enable verbose logging & standard output.")
	executeCommand.PersistentFlags().Bool("tracking", false, "Enable tracking for Azure resources created by the Azure CLI commands executed.")
	executeCommand.PersistentFlags().Bool("do-not-delete", false, "Do not delete the Azure resources created by the Azure CLI commands executed.")
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
		tracking, _ := cmd.Flags().GetBool("tracking")
		do_not_delete, _ := cmd.Flags().GetBool("do-not-delete")

		innovationEngine := engine.NewEngine(engine.EngineConfiguration{
			Verbose:          verbose,
			ResourceTracking: tracking,
			DoNotDelete:      do_not_delete,
		})
		scenario, err := engine.CreateScenarioFromMarkdown(markdownFile, []string{"bash", "azurecli", "azurecli-interactive", "terraform"})
		if err != nil {
			panic(err)
		}

		innovationEngine.ExecuteScenario(scenario)
	},
}
