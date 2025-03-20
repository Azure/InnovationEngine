package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eg",
	Short: "A CLI application to discover and work with Executable Docs using existing documentation as examples.",
	Long: `EG (meaning "for example") is a command line tool that assists in finding, customizing and executing Executable Docs.\n
Eg uses Copilot to interact with existing documentation in order to create custom executable docs.\n 
It then uses Innovation Engine (IE) to execute these docs.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define any flags and configuration settings here
}
