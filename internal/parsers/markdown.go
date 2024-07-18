package parsers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var markdownParser = goldmark.New(
	goldmark.WithExtensions(extension.GFM, meta.New(meta.WithStoresInDocument())),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithBlockParsers(),
	),
	goldmark.WithRendererOptions(
		html.WithXHTML(),
	),
)

// Parses a markdown file into an AST representing the markdown document.
func ParseMarkdownIntoAst(source []byte) ast.Node {
	document := markdownParser.Parser().Parse(text.NewReader(source))
	return document
}

func ExtractYamlMetadataFromAst(node ast.Node) map[string]interface{} {
	return node.OwnerDocument().Meta()
}

// The representation of an expected output block in a markdown file. This is
// for scenarios that have expected output that should be validated against the
// actual output.
type ExpectedOutputBlock struct {
	Language           string
	Content            string
	ExpectedSimilarity float64
	ExpectedRegex      *regexp.Regexp
}

// The representation of a code block in a markdown file.
type CodeBlock struct {
	Language       string
	Content        string
	Header         string
	Description    string
	ExpectedOutput ExpectedOutputBlock
}

// Assumes the title of the scenario is the first h1 header in the
// markdown file.
func ExtractScenarioTitleFromAst(node ast.Node, source []byte) (string, error) {
	header := ""
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch n := node.(type) {
			case *ast.Heading:
				if n.Level == 1 {
					header = string(extractTextFromMarkdown(&n.BaseBlock, source))
					return ast.WalkStop, nil
				}
			}
		}
		return ast.WalkContinue, nil
	})

	if header == "" {
		return "", fmt.Errorf("no h1 header found to use as the scenario title")
	}

	return header, nil
}

var expectedSimilarityRegex = regexp.MustCompile(
	`<!--\s*expected_similarity=\s*(\d+\.?\d*)|"(.*)"\s*-->`,
)

// Extracts the code blocks from a provided markdown AST that match the
// languagesToExtract.
func ExtractCodeBlocksFromAst(
	node ast.Node,
	source []byte,
	languagesToExtract []string,
) []CodeBlock {
	var lastHeader string
	var commands []CodeBlock
	var nextBlockIsExpectedOutput bool
	var lastExpectedSimilarityScore float64
	var lastExpectedRegex *regexp.Regexp
	var lastNode ast.Node

	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch n := node.(type) {
			// Set the last header when we encounter a heading.
			case *ast.Heading:
				lastHeader = string(extractTextFromMarkdown(&n.BaseBlock, source))
				lastNode = node
			case *ast.Paragraph:
				lastNode = node
			// Extract the code block if it matches the language.
			case *ast.HTMLBlock:
				content := extractTextFromMarkdown(&n.BaseBlock, source)
				matches := expectedSimilarityRegex.FindStringSubmatch(content)

				if len(matches) < 3 {
					break
				}

				match := matches[1]
				if match != "" {
					score, err := strconv.ParseFloat(match, 64)
					logging.GlobalLogger.Debugf("Simalrity score of %f found", score)
					if err != nil {
						return ast.WalkStop, err
					}
					lastExpectedSimilarityScore = score
				} else {
					match = matches[2]
					logging.GlobalLogger.Debugf("Regex %q found", match)

					if match == "" {
						return ast.WalkStop, errors.New("No regex found")
					}

					re, err := regexp.Compile(match)
					if err != nil {
						return ast.WalkStop, fmt.Errorf("Cannot compile the following regex: %q", match)
					}

					lastExpectedRegex = re
				}

				nextBlockIsExpectedOutput = true
			case *ast.FencedCodeBlock:
				language := string(n.Language((source)))
				content := extractTextFromMarkdown(&n.BaseBlock, source)
				description := ""

				if lastNode != nil {
					switch n := lastNode.(type) {
					case *ast.Paragraph:
						description = string(extractTextFromMarkdown(&n.BaseBlock, source))
					default:
						logging.GlobalLogger.Warnf("The node before the codeblock `%s` is not a paragraph, it is a %s", content, n.Kind())
					}
				} else {
					logging.GlobalLogger.Warnf("There are no markdown elements before the last codeblock `%s`", content)
				}

				lastNode = node
				for _, desiredLanguage := range languagesToExtract {
					if language == desiredLanguage {
						command := CodeBlock{
							Language:    language,
							Content:     content,
							Header:      lastHeader,
							Description: description,
						}
						commands = append(commands, command)
						break
					} else if nextBlockIsExpectedOutput {
						// Map the expected output to the last command. If there
						// are no commands, then we ignore the expected output.
						if len(commands) > 0 {
							expectedOutputBlock := ExpectedOutputBlock{
								Language:           language,
								Content:            extractTextFromMarkdown(&n.BaseBlock, source),
								ExpectedSimilarity: lastExpectedSimilarityScore,
								ExpectedRegex:      lastExpectedRegex,
							}
							commands[len(commands)-1].ExpectedOutput = expectedOutputBlock

							// Reset the expected output state.
							nextBlockIsExpectedOutput = false
							lastExpectedSimilarityScore = 0
							lastExpectedRegex = nil
						}
						break
					}
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return commands
}

// This regex matches HTML comments within markdown blocks that contain
// variables to use within the scenario.
var variableCommentBlockRegex = regexp.MustCompile("(?s)<!--.*?```variables(.*?)```.*?")

// Extracts the variables from a provided markdown AST.
func ExtractScenarioVariablesFromAst(node ast.Node, source []byte) map[string]string {
	scenarioVariables := make(map[string]string)

	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == ast.KindHTMLBlock {
			htmlNode := node.(*ast.HTMLBlock)
			blockContent := extractTextFromMarkdown(&htmlNode.BaseBlock, source)
			logging.GlobalLogger.Debugf("Found HTML block with the content: %s\n", blockContent)
			match := variableCommentBlockRegex.FindStringSubmatch(blockContent)

			// Extract the variables from the comment block.
			if len(match) > 1 {
				variables := convertScenarioVariablesToMap(match[1])
				for key, value := range variables {
					scenarioVariables[key] = value
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return scenarioVariables
}

// Converts a string of shell variable exports into a map of key/value pairs.
// I.E. `export FOO=bar\nexport BAZ=qux` becomes `{"FOO": "bar", "BAZ": "qux"}`
func convertScenarioVariablesToMap(variableBlock string) map[string]string {
	variableMap := make(map[string]string)

	// Only process statements that begin with export.
	for _, variable := range strings.Split(variableBlock, "\n") {
		if strings.HasPrefix(variable, "export") {
			parts := strings.SplitN(variable, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], "export ")
				value := parts[1]
				logging.GlobalLogger.Debugf("Found variable: %s=%s\n", key, value)
				variableMap[key] = value
			}
		}
	}

	return variableMap
}

// Extract the text from a code blocks base block and return it as a string.
func extractTextFromMarkdown(baseBlock *ast.BaseBlock, source []byte) string {
	lines := baseBlock.Lines()
	var command strings.Builder

	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		command.WriteString(string(line.Value(source)))
	}

	return command.String()
}
