package commands

import (
	"github.com/spf13/cobra"
)

var interactiveCommand = &cobra.Command{
	Use:   "interactive",
	Short: "Execute a document in interactive mode.",
}

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(interactiveCommand)
}
