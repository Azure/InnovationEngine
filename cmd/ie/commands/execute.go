package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/render"
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

		// Check if the markdown file exists.
		if !utils.FileExists(markdownFile) {
			fmt.Printf("Markdown file '%s' does not exist.\n", markdownFile)
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
			fmt.Printf("INI file '%s' does not exist, skipping...", markdownINI)
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

		codeBlocks := parsers.ExtractCodeBlocksFromAst(markdown, source, []string{"bash", "azurecli", "azurecli-init", "azurecli-interactive", "terraform", "terraform-interactive"})

		render.ExecuteAndRenderCodeBlocks(codeBlocks, environmentVariables)
	},
}
