package commands

import (
	"fmt"
	"os"

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
			panic(err)
		}
		logging.Init(logging.LevelFromString(logLevel))
	},
}

// Entrypoint into the Innovation Engine CLI.
func ExecuteCLI() {
	rootCommand.PersistentFlags().String("log-level", string(logging.Debug), "Configure the log level")

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
