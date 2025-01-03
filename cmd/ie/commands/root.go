package commands

import (
	"fmt"
	"os"

	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/spf13/cobra"
)

// The root command for the CLI. Currently initializes the logging for all other
// commands.
var rootCommand = &cobra.Command{
	Use:   "ie",
	Short: "The innovation engine.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel, err := cmd.Flags().GetString("log-level")
		if err != nil {
			fmt.Printf("Error getting log level: %s", err)
			os.Exit(1)
		}
		logging.Init(logging.LevelFromString(logLevel))

		// Check environment
		environment, err := cmd.Flags().GetString("environment")
		if err != nil {
			fmt.Printf("Error getting environment: %s", err)
			logging.GlobalLogger.Errorf("Error getting environment: %s", err)
			os.Exit(1)
		}

		if !environments.IsValidEnvironment(environment) {
			fmt.Printf("Invalid environment: %s", environment)
			logging.GlobalLogger.Errorf("Invalid environment: %s", err)
			os.Exit(1)
		}
	},
}

// Entrypoint into the Innovation Engine CLI.
func ExecuteCLI() {
	rootCommand.PersistentFlags().
		String(
			"log-level",
			string(logging.Debug),
			"Configure the log level for statements written to ie.log. Valid options are 'trace', 'debug', 'info', 'warn', 'error', 'fatal'",
		)
	rootCommand.PersistentFlags().
		String(
			"environment",
			environments.EnvironmentsLocal,
			"The environment that the CLI is running in. Valid options are 'local', 'github-action'. For running ie in your standard terminal, local will work just fine. If using IE inside a github action, use github-action.",
		)

	rootCommand.PersistentFlags().
		StringArray(
			"feature",
			[]string{},
			"Enables the specified feature. Format: --feature <feature>",
		)

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		logging.GlobalLogger.Errorf("Failed to execute ie: %s", err)
		os.Exit(1)
	}
}
