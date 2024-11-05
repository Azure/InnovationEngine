package commands

import "github.com/spf13/cobra"

var (
	VERSION = "dev"
	COMMIT  = "N/A"
	DATE    = "N/A"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the Innovation Engine",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Printf("Version: %s\n", VERSION)
		cmd.Printf("Commit: %s\n", COMMIT)
		cmd.Printf("Date: %s\n", DATE)

		return nil
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
}
