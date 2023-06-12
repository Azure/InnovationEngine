package commands

import (
	"github.com/Azure/InnovationEngine/internal/engine"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(testCommand)
	testCommand.PersistentFlags().Bool("verbose", false, "Enable verbose logging & standard output.")
}

var testCommand = &cobra.Command{
	Use:   "test",
	Args:  cobra.MinimumNArgs(1),
	Short: "Test document commands against it's expected outputs.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO(vmarcella): Initialize this via a flag.
		logging.Init("info")

		markdownFile := args[0]
		if markdownFile == "" {
			cmd.Help()
			return
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		innovationEngine := engine.NewEngine(engine.EngineConfiguration{
			Verbose:          verbose,
			ResourceTracking: false,
			DoNotDelete:      false,
		})

		scenario, err := engine.CreateScenarioFromMarkdown(markdownFile, []string{"bash", "azurecli", "azurecli-interactive", "terraform"})
		if err != nil {
			panic(err)
		}

		innovationEngine.TestScenario(scenario)

	},
}
