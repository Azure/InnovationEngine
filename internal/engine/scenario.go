package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/utils"
	"github.com/yuin/goldmark/ast"
)

// Individual steps within a scenario.
type Step struct {
	Name       string
	CodeBlocks []parsers.CodeBlock
}

// Scenarios are the top-level object that represents a scenario to be executed.
type Scenario struct {
	Name        string
	MarkdownAst ast.Node
	Steps       []Step
	Environment map[string]string
}

func groupCodeBlocksIntoSteps(blocks []parsers.CodeBlock) []Step {
	var groupedSteps []Step
	var headerIndex = make(map[string]int)

	for _, block := range blocks {
		if index, ok := headerIndex[block.Header]; ok {
			groupedSteps[index].CodeBlocks = append(groupedSteps[index].CodeBlocks, block)
		} else {
			headerIndex[block.Header] = len(groupedSteps)
			groupedSteps = append(groupedSteps, Step{
				Name:       block.Header,
				CodeBlocks: []parsers.CodeBlock{block},
			})
		}
	}

	return groupedSteps
}

// Creates a scenario object from a given markdown file. languagesToExecute is
// used to filter out code blocks that should not be parsed out of the markdown
// file.
func CreateScenarioFromMarkdown(path string, languagesToExecute []string) (*Scenario, error) {
	if path == "" {
		return nil, nil
	}

	if !utils.FileExists(path) {
		return nil, fmt.Errorf("markdown file '%s' does not exist", path)
	}

	source, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Load environment variables
	markdownINI := strings.TrimSuffix(path, filepath.Ext(path)) + ".ini"
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

	markdown := parsers.ParseMarkdownIntoAst(source)
	scenarioVariables := parsers.ExtractScenarioVariablesFromAst(markdown, source)
	for key, value := range scenarioVariables {
		environmentVariables[key] = value
	}

	codeBlocks := parsers.ExtractCodeBlocksFromAst(markdown, source, languagesToExecute)
	steps := groupCodeBlocksIntoSteps(codeBlocks)
	title, err := parsers.ExtractScenarioTitleFromAst(markdown, source)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found scenario: %s\n", title)

	return &Scenario{
		Name:        title,
		Environment: environmentVariables,
		Steps:       steps,
		MarkdownAst: markdown,
	}, nil
}
