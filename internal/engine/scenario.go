package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/lib/fs"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/parsers"
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

// Groups the codeblocks into steps based on the header of the codeblock.
// This organizes the codeblocks into steps that can be executed linearly.
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
func CreateScenarioFromMarkdown(
	path string,
	languagesToExecute []string,
	environmentVariableOverrides map[string]string,
) (*Scenario, error) {
	if !fs.FileExists(path) {
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
	if !fs.FileExists(markdownINI) {
		logging.GlobalLogger.Infof("INI file '%s' does not exist, skipping...", markdownINI)
	} else {
		logging.GlobalLogger.Infof("INI file '%s' exists, loading...", markdownINI)
		environmentVariables, err = parsers.ParseINIFile(markdownINI)

		if err != nil {
			return nil, err
		}

		for key, value := range environmentVariables {
			logging.GlobalLogger.Debugf("Setting %s=%s\n", key, value)
		}
	}

	// Convert the markdonw into an AST and extract the scenario variables.
	markdown := parsers.ParseMarkdownIntoAst(source)
	scenarioVariables := parsers.ExtractScenarioVariablesFromAst(markdown, source)
	for key, value := range scenarioVariables {
		environmentVariables[key] = value
	}

	// Extract the code blocks from the markdown file.
	codeBlocks := parsers.ExtractCodeBlocksFromAst(markdown, source, languagesToExecute)
	logging.GlobalLogger.WithField("CodeBlocks", codeBlocks).
		Debugf("Found %d code blocks", len(codeBlocks))

	varsToExport := lib.CopyMap(environmentVariableOverrides)
	for key, value := range environmentVariableOverrides {
		logging.GlobalLogger.Debugf("Attempting to override %s with %s", key, value)
		exportRegex := regexp.MustCompile(fmt.Sprintf(`export %s=["']?([a-z-A-Z0-9_]+)["']?`, key))

		for index, codeBlock := range codeBlocks {
			matches := exportRegex.FindAllStringSubmatch(codeBlock.Content, -1)

			if len(matches) != 0 {
				logging.GlobalLogger.Debugf(
					"Found %d matches for %s, deleting from varsToExport",
					len(matches),
					key,
				)
				delete(varsToExport, key)
			} else {
				logging.GlobalLogger.Debugf("Found no matches for %s inside of %s", key, codeBlock.Content)
			}

			for _, match := range matches {
				oldLine := match[0]
				oldValue := match[1]

				// Replace the old export with the new export statement
				newLine := strings.Replace(oldLine, oldValue, value, 1)
				logging.GlobalLogger.Debugf("Replacing '%s' with '%s'", oldLine, newLine)

				// Update the code block with the new export statement
				codeBlocks[index].Content = strings.Replace(codeBlock.Content, oldLine, newLine, 1)
			}

		}
	}

	// If there are some variables left after going through each of the codeblocks,
	// do not update the scenario
	// steps.
	if len(varsToExport) != 0 {
		logging.GlobalLogger.Debugf(
			"Found %d variables to add to the scenario as a step.",
			len(varsToExport),
		)
		exportCodeBlock := parsers.CodeBlock{
			Language:       "bash",
			Content:        "",
			Header:         "Exporting variables defined via the CLI and not in the markdown file.",
			ExpectedOutput: parsers.ExpectedOutputBlock{},
		}
		for key, value := range varsToExport {
			exportCodeBlock.Content += fmt.Sprintf("export %s=\"%s\"\n", key, value)
		}

		codeBlocks = append([]parsers.CodeBlock{exportCodeBlock}, codeBlocks...)
	}

	// Group the code blocks into steps.
	steps := groupCodeBlocksIntoSteps(codeBlocks)
	title, err := parsers.ExtractScenarioTitleFromAst(markdown, source)
	if err != nil {
		return nil, err
	}

	logging.GlobalLogger.Infof("Successfully built out the scenario: %s", title)

	return &Scenario{
		Name:        title,
		Environment: environmentVariables,
		Steps:       steps,
		MarkdownAst: markdown,
	}, nil
}

func (s *Scenario) OverwriteEnvironmentVariables(environmentVariables map[string]string) {
	for key, value := range environmentVariables {
		s.Environment[key] = value
	}
}

// Convert a scenario into a shell script
func (s *Scenario) ToShellScript() string {
	var script strings.Builder

	for key, value := range s.Environment {
		script.WriteString(fmt.Sprintf("export %s=\"%s\"\n", key, value))
	}

	for _, step := range s.Steps {
		script.WriteString(fmt.Sprintf("# %s\n", step.Name))
		for _, block := range step.CodeBlocks {
			script.WriteString(fmt.Sprintf("%s\n", block.Content))
		}
	}

	return script.String()
}
