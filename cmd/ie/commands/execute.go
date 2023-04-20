package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/spf13/cobra"
)

var markdownFile string

// / Register the command with our command runner.
func init() {
	rootCommand.AddCommand(executeCommand)
	executeCommand.Flags().StringVar(&markdownFile, "markdown", "", "The markdown file to execute.")
}

var executeCommand = &cobra.Command{
	Use:   "execute",
	Short: "Execute a document.",
	Run: func(cmd *cobra.Command, args []string) {
		if markdownFile == "" {
			cmd.Help()
			return
		}
		source, err := ioutil.ReadFile(markdownFile)

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
