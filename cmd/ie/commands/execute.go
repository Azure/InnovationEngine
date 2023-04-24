package commands

import (
	"fmt"
	"os"

	"github.com/Azure/InnovationEngine/internal/parsers"
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
		source, err := os.ReadFile(markdownFile)

		if err != nil {
			panic(err)
		}

		markdown := parsers.ParseMarkdownIntoAst(source)
		commands := parsers.ExtractCodeBlocksFromAst(markdown, source, []string{"bash", "azurecli", "azurecli-interactive", "terraform", "terraform-interactive"})

		for _, command := range commands {
			fmt.Println(command)
		}
	},
}
