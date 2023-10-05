package commands

import (
	"errors"
	"fmt"

	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(testCommand)
	testCommand.PersistentFlags().
		Bool("verbose", false, "Enable verbose logging & standard output.")
	testCommand.PersistentFlags().
		String("subscription", "", "Sets the subscription ID used by a scenarios azure-cli commands. Will rely on the default subscription if not set.")
}

var testCommand = &cobra.Command{
	Use:   "test [markdown file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Test document commands against it's expected outputs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		markdownFile := args[0]
		if markdownFile == "" {
			return errors.New("no markdown file specified")
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
			return fmt.Errorf("creating engine: %w", err)
		}

		scenario, err := engine.CreateScenarioFromMarkdown(
			markdownFile,
			[]string{"bash", "azurecli", "azurecli-interactive", "terraform"},
			nil,
		)
		if err != nil {
			return fmt.Errorf("creating scenario: %w", err)
		}

		if err := innovationEngine.TestScenario(scenario); err != nil {
			return fmt.Errorf("testing scenario: %w", err)
		}

		return nil
	},
}
