package parsers

import (
	"fmt"
	"testing"
)

func TestParsingMarkdownHeaders(t *testing.T) {
	t.Run("Markdown with a valid title", func(t *testing.T) {
		markdown := []byte(`# Hello World`)
		document := ParseMarkdownIntoAst(markdown)
		title, err := ExtractScenarioTitleFromAst(document, markdown)

		if err != nil {
			t.Errorf("Error parsing title: %s", err)
		}

		if title != "Hello World" {
			t.Errorf("Title is wrong: %s", title)
		}
	})

	t.Run("Markdown with multiple titles", func(t *testing.T) {
		markdown := []byte("# Hello World \n # Hello again")
		document := ParseMarkdownIntoAst(markdown)
		title, err := ExtractScenarioTitleFromAst(document, markdown)

		if err != nil {
			t.Errorf("Error parsing title: %s", err)
		}

		if title != "Hello World" {
			t.Errorf("Title is wrong: %s", title)
		}
	})

	t.Run("Markdown without a title", func(t *testing.T) {
		markdown := []byte(``)

		document := ParseMarkdownIntoAst(markdown)
		title, err := ExtractScenarioTitleFromAst(document, markdown)

		if err == nil {
			t.Errorf("Error should have been thrown")
		}

		if title != "" {
			t.Errorf("Title should be empty")
		}
	})
}

func TestParsingMarkdownCodeBlocks(t *testing.T) {

	t.Run("Markdown with a valid bash code block", func(t *testing.T) {
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
			t.Errorf(
				"Code block code is wrong. Expected: %s, Got %s",
				"echo Hello\\n",
				codeBlocks[0].Content,
			)
		}
	})

}

func TestParsingMarkdownExpectedSimilarty(t *testing.T) {

	t.Run("Markdown with a expected_similarty tag using float", func(t *testing.T) {
		markdown := []byte(fmt.Sprintf("```bash\n%s\n```\n<!--expected_similarity=0.8-->\n```\nHello\n```\n", "echo Hello"))

		document := ParseMarkdownIntoAst(markdown)
		codeBlocks := ExtractCodeBlocksFromAst(document, markdown, []string{"bash"})

		if len(codeBlocks) != 1 {
			t.Errorf("Code block count is wrong: %d", len(codeBlocks))
		}

		block := codeBlocks[0].ExpectedOutput
		expectedFloat := .8
		if block.ExpectedSimilarity != expectedFloat {
			t.Errorf("ExpectedSimilarity is wrong, got %f, expected %f", block.ExpectedSimilarity, expectedFloat)
		}
	})

}

func TestParsingMarkdownExpectedRegex(t *testing.T) {

	t.Run("Markdown with a expected_similarty tag using regex", func(t *testing.T) {
		markdown := []byte(fmt.Sprintf("```bash\n%s\n```\n<!--expected_similarity=\"Foo \\w+\"-->\n```\nFoo Bar\n```\n", "echo 'Foo Bar'"))

		document := ParseMarkdownIntoAst(markdown)
		codeBlocks := ExtractCodeBlocksFromAst(document, markdown, []string{"bash"})

		if len(codeBlocks) != 1 {
			t.Errorf("Code block count is wrong: %d", len(codeBlocks))
		}

		block := codeBlocks[0].ExpectedOutput
		if block.ExpectedRegex == nil {
			t.Errorf("ExpectedRegex is nil")
		}

		stringRegex := block.ExpectedRegex.String()
		expectedRegex := `Foo \w+`
		if stringRegex != expectedRegex {
			t.Errorf("ExpectedRegex is wrong, got %q, expected %q", stringRegex, expectedRegex)
		}
	})

}
