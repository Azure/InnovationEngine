package commands

import (
	"fmt"
	"os"

	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/spf13/cobra"
)

// Register the command with our command runner.
func init() {
	rootCommand.AddCommand(clearEnvCommand)

	// Bool flags
	clearEnvCommand.PersistentFlags().
		Bool("all", false, "Clear both environment variables and working directory state.")
	clearEnvCommand.PersistentFlags().
		Bool("working-dir", false, "Also clear the working directory state.")
	clearEnvCommand.PersistentFlags().
		Bool("force", false, "Force clear without confirmation prompt.")
}

var clearEnvCommand = &cobra.Command{
	Use:   "clear-env",
	Short: "Clear the stored environment variables and optionally working directory state.",
	Long: `Clear the stored environment variables and optionally working directory state.
	
This command removes the environment state file that stores variables between
Innovation Engine command executions. By default, it only clears environment 
variables, but you can also clear working directory state using the flags.

Examples:
  ie clear-env                    # Clear only environment variables
  ie clear-env --working-dir      # Clear env vars and working directory
  ie clear-env --all              # Clear both env vars and working directory
  ie clear-env --force            # Clear without confirmation`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		clearAll, _ := cmd.Flags().GetBool("all")
		clearWorkingDir, _ := cmd.Flags().GetBool("working-dir")

		// Determine what to clear
		shouldClearEnv := true // Always clear env vars
		shouldClearWD := clearAll || clearWorkingDir

		// Show confirmation unless --force is used
		if !force {
			fmt.Print("This will clear the stored environment state")
			if shouldClearWD {
				fmt.Print(" and working directory state")
			}
			fmt.Print(". Continue? (y/N): ")
			
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" && response != "yes" {
				fmt.Println("Operation cancelled.")
				return
			}
		}

		// Clear environment variables
		if shouldClearEnv {
			if err := lib.DeleteEnvironmentStateFile(lib.DefaultEnvironmentStateFile); err != nil {
				// Don't error if file doesn't exist
				if !os.IsNotExist(err) {
					logging.GlobalLogger.Errorf("Error clearing environment variables: %s", err)
					fmt.Printf("Error clearing environment variables: %s\n", err)
					os.Exit(1)
				} else {
					fmt.Println("Environment variables state file was already clear.")
				}
			} else {
				fmt.Println("Environment variables cleared successfully.")
			}
		}

		// Clear working directory state if requested
		if shouldClearWD {
			if err := lib.DeleteWorkingDirectoryStateFile(lib.DefaultWorkingDirectoryStateFile); err != nil {
				// Don't error if file doesn't exist
				if !os.IsNotExist(err) {
					logging.GlobalLogger.Errorf("Error clearing working directory state: %s", err)
					fmt.Printf("Error clearing working directory state: %s\n", err)
					os.Exit(1)
				} else {
					fmt.Println("Working directory state file was already clear.")
				}
			} else {
				fmt.Println("Working directory state cleared successfully.")
			}
		}

		logging.GlobalLogger.Info("Environment state cleared successfully")
	},
}