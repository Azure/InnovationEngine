package parsers

import (
	"fmt"
	"testing"
)

func TestParsingMarkdownTitle(t *testing.T) {
	// Handle when title is present
	markdown := []byte(`# Hello World`)
	document := ParseMarkdownIntoAst(markdown)
	title, err := ExtractScenarioTitleFromAst(document, markdown)

	if err != nil {
		t.Errorf("Error parsing title: %s", err)
	}

	if title != "Hello World" {
		t.Errorf("Title is wrong: %s", title)
	}

	// Handle when title is not present
	markdown = []byte(``)

	document = ParseMarkdownIntoAst(markdown)
	title, err = ExtractScenarioTitleFromAst(document, markdown)

	if err == nil {
		t.Errorf("Error should have been thrown")
	}

	if title != "" {
		t.Errorf("Title should be empty")
	}

}

func TestParsingMarkdownCodeBlocks(t *testing.T) {
	markdown := []byte(fmt.Sprintf("# Hello World\n ```bash\n%s\n```", "echo Hello"))

	document := ParseMarkdownIntoAst(markdown)
	codeBlocks := ExtractCodeBlocksFromAst(document, markdown, []string{"bash"})

	if len(codeBlocks) != 1 {
		t.Errorf("Code block count is wrong: %d", len(codeBlocks))
	}

	if codeBlocks[0].Language != "bash" {
		t.Errorf("Code block language is wrong: %s", codeBlocks[0].Language)
	}

	if codeBlocks[0].Content != "echo Hello\n" {
		t.Errorf("Code block code is wrong. Expected: %s, Got %s", "echo Hello\\n", codeBlocks[0].Content)
	}
}
