package parsers

import (
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
