package commands

import (
	"github.com/spf13/cobra"
)

var executeCommand = &cobra.Command{
	Use:   "execute",
	Short: "Execute a document.",
}

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(executeCommand)
}
