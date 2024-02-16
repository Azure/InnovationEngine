package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerboseStyle(t *testing.T) {
	text := `aaaa
  b`
	styledText := VerboseStyle.Render(text)
	expectedStyledText := `aaaa
  b `
	assert.Equal(t, expectedStyledText, styledText)
	assert.Equal(t, text, RemoveHorizontalAlign(styledText))
}
