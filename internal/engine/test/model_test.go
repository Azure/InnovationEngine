package test

import (
	"testing"

	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/parsers"
)

func TestTestModeModel(t *testing.T) {
	t.Run("Initializing a test model with an invalid subscription fails.", func(t *testing.T) {
		// Test the initialization of the test mode model.
		_, err := NewTestModeModel("test", "invalid", "test", nil, nil)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("Creating a test model works.", func(t *testing.T) {
		// Test the initialization of the test mode model.
		model, err := NewTestModeModel("test", "", "test", nil, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if model.components.commandViewport.View() != "" {
			t.Errorf("Expected view to be empty, got %q", model.components.commandViewport.View())
		}

		if model.scenarioTitle != "test" {
			t.Errorf("Expected scenario title to be %q, got %q", "test", model.scenarioTitle)
		}

		if model.env != nil {
			t.Errorf("Expected env to be nil, got %v", model.env)
		}
	})

	t.Run("Creating a test model with steps works.", func(t *testing.T) {
		// Test the initialization of the test mode model.

		steps := []common.Step{
			{
				Name: "step1",
				CodeBlocks: []parsers.CodeBlock{
					{
						Header:   "header1",
						Content:  "echo 'hello world'",
						Language: "bash",
					},
				},
			},
		}

		model, err := NewTestModeModel("test", "", "test", steps, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if model.currentCodeBlock != 0 {
			t.Errorf("Expected current code block to be 0, got %d", model.currentCodeBlock)
		}

		if model.codeBlockState[0].CodeBlock.Content != "echo 'hello world'" {
			t.Errorf(
				"Expected code block content to be %q, got %q",
				"echo 'hello world'",
				model.codeBlockState[0].CodeBlock.Content,
			)
		}
	})
}
