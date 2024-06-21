package commands

import (
	"fmt"
	"os"

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
		workingDirectory, _ := cmd.Flags().GetString("working-directory")
		environment, _ := cmd.Flags().GetString("environment")

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
			nil,
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
