package parsers

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var markdownParser = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
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

// Extracts the code blocks from a provided markdown AST that match the
// languagesToExtract.
func ExtractCodeBlocksFromAst(node ast.Node, source []byte, languagesToExtract []string) []string {
	var commands []string
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == ast.KindFencedCodeBlock {
			codeBlock := node.(*ast.FencedCodeBlock)
			for _, language := range languagesToExtract {
				if string(codeBlock.Language(source)) == language {
					commands = append(commands, extractCommandFromCodeBlock(codeBlock, source))
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return commands
}

// Extracts the command text from an already parsed markdown code block.
func extractCommandFromCodeBlock(codeBlock *ast.FencedCodeBlock, source []byte) string {
	lines := codeBlock.Lines()
	var command strings.Builder

	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		command.WriteString(string(line.Value(source)))
	}

	return command.String()
}
