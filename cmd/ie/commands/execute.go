package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/Azure/InnovationEngine/internal/utils"
	"github.com/spf13/cobra"
)

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(executeCommand)
}

var executeCommand = &cobra.Command{
	Use:   "execute [markdown file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Execute the commands for an Azure deployment scenario.",
	Run: func(cmd *cobra.Command, args []string) {
		markdownFile := args[0]
		if markdownFile == "" {
			cmd.Help()
			return
		}

		// Load the markdown file.
		if utils.FileExists(markdownFile) == false {
			fmt.Println("File does not exist.")
			return
		}

		source, err := os.ReadFile(markdownFile)
		if err != nil {
			panic(err)
		}

		// Load environment variables
		markdownINI := strings.TrimSuffix(markdownFile, filepath.Ext(markdownFile)) + ".ini"
		environmentVariables := make(map[string]string)

		// Check if the INI file exists & load it.
		if !utils.FileExists(markdownINI) {
			fmt.Println("INI file does not exist: ", markdownINI)
		} else {
			fmt.Println("INI file exists. Loading: ", markdownINI)
			environmentVariables = parsers.ParseINIFile(markdownINI)

			for key, value := range environmentVariables {
				fmt.Printf("Setting %s=%s\n", key, value)
			}
		}

		fmt.Println(environmentVariables)

		markdown := parsers.ParseMarkdownIntoAst(source)
		scenarioVariables := parsers.ExtractScenarioVariablesFromAst(markdown, source)
		for key, value := range scenarioVariables {
			environmentVariables[key] = value
		}

		commands := parsers.ExtractCodeBlocksFromAst(markdown, source, []string{"bash", "azurecli", "azurecli-init", "azurecli-interactive", "terraform", "terraform-interactive"})

		fmt.Println(scenarioVariables)

		for _, command := range commands {
			fmt.Println(command)
			out, error := shells.ExecuteBashCommand(command, environmentVariables, true)
			if error != nil {
				fmt.Println(error)
			}
			fmt.Println(out)
		}
	},
}
