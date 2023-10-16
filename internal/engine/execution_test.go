package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteBlock(t *testing.T) {
	blocks := []string{
		"echo \"hello \\\nworld\"", // tutorial.md
		"echo hello \\\nworld",
		"echo \"hello world\"",
		"echo hello world",
	}
	for _, blockCommand := range blocks {
		t.Run("render command", func(t *testing.T) {
			_, err := renderCommand(blockCommand)
			assert.Equal(t, nil, err)
		})
	}

}
