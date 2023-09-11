package commands

import (
	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(testCommand)
	testCommand.PersistentFlags().Bool("verbose", false, "Enable verbose logging & standard output.")
	testCommand.PersistentFlags().String("subscription", "", "Sets the subscription ID used by a scenarios azure-cli commands. Will rely on the default subscription if not set.")
}

var testCommand = &cobra.Command{
	Use:   "test",
	Args:  cobra.MinimumNArgs(1),
	Short: "Test document commands against it's expected outputs.",
	Run: func(cmd *cobra.Command, args []string) {

		markdownFile := args[0]
		if markdownFile == "" {
			cmd.Help()
			return
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		subscription, _ := cmd.Flags().GetString("subscription")

		innovationEngine := engine.NewEngine(engine.EngineConfiguration{
			Verbose:       verbose,
			DoNotDelete:   false,
			Subscription:  subscription,
			CorrelationId: "",
		})

		scenario, err := engine.CreateScenarioFromMarkdown(markdownFile, []string{"bash", "azurecli", "azurecli-interactive", "terraform"}, nil)
		if err != nil {
			panic(err)
		}

		innovationEngine.TestScenario(scenario)

	},
}
