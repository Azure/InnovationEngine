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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logLevel, err := cmd.Flags().GetString("log-level")
		if err != nil {
			return fmt.Errorf("getting log level: %w", err)
		}
		logging.Init(logging.LevelFromString(logLevel))

		// Check environment
		environment, err := cmd.Flags().GetString("environment")
		if err != nil {
			return fmt.Errorf("getting environment: %w", err)
		}

		if !environments.IsValidEnvironment(environment) {
			return fmt.Errorf("validating environment: %w", err)
		}

		return nil
	},
}

// Entrypoint into the Innovation Engine CLI.
func ExecuteCLI() {
	rootCommand.PersistentFlags().
		String("log-level", string(logging.Debug), "Configure the log level")
	rootCommand.PersistentFlags().
		String("environment", environments.EnvironmentsLocal, "The environment that the CLI is running in. (local, ci, ocd)")

	rootCommand.PersistentFlags().
		StringArray("feature", []string{}, "Enables the specified feature. Format: --feature <feature>")

	if err := rootCommand.Execute(); err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		logging.GlobalLogger.Errorf("Error executing command: %s", err)
		os.Exit(1)
	}
}
