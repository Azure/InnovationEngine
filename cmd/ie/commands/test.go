package commands

import (
	"github.com/spf13/cobra"
)

var testCommand = &cobra.Command{
	Use:   "test",
	Short: "Test document commands against it's expected outputs.",
}

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(testCommand)
}
