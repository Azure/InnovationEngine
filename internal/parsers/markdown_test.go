package parsers

import (
	"testing"
)

func TestParsingTitle(t *testing.T) {
	markdown := []byte(`# Hello World`)
	document := ParseMarkdownIntoAst(markdown)
	title, err := ExtractScenarioTitleFromAst(document, markdown)

	if err != nil {
		t.Errorf("Error parsing title: %s", err)
	}

	if title != "Hello World" {
		t.Errorf("Title is wrong: %s", title)
	}
}
