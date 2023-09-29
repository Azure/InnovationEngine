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

		innovationEngine, err := engine.NewEngine(engine.EngineConfiguration{
			Verbose:       verbose,
			DoNotDelete:   false,
			Subscription:  subscription,
			CorrelationId: "",
		})

		if err != nil {
			logging.GlobalLogger.Errorf("Error creating engine %s", err)
			fmt.Printf("Error creating engine %s", err)
			os.Exit(1)
		}

		scenario, err := engine.CreateScenarioFromMarkdown(markdownFile, []string{"bash", "azurecli", "azurecli-interactive", "terraform"}, nil)
		if err != nil {
			logging.GlobalLogger.Errorf("Error creating scenario %s", err)
			fmt.Printf("Error creating engine %s", err)
			os.Exit(1)
		}

		innovationEngine.TestScenario(scenario)

	},
}
