package test

import (
	"testing"

	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/stretchr/testify/assert"
)

// This suite of tests is responsible for ensuring that the model around test mode
// is well defined and behaves as expected.
func TestTestModeModel(t *testing.T) {
	t.Run("Initializing a test model with an invalid subscription fails.", func(t *testing.T) {
		// Test the initialization of the test mode model.
		_, err := NewTestModeModel("test", "invalid", "test", nil, nil)
		assert.Error(t, err)
	})

	t.Run("Creating a valid test model works.", func(t *testing.T) {
		// Test the initialization of the test mode model.
		model, err := NewTestModeModel("test", "", "test", nil, nil)
		assert.NoError(t, err)

		assert.Equal(t, "test", model.scenarioTitle)

		assert.Equal(t, "", model.components.commandViewport.View())
		assert.Equal(t, map[string]string{}, model.env)
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
		assert.NoError(t, err)

		assert.Equal(t, 0, model.currentCodeBlock)
		assert.Equal(t, 1, len(model.codeBlockState))

		state := model.codeBlockState[0]

		assert.Equal(t, "step1", state.StepName)
		assert.Equal(t, "bash", state.CodeBlock.Language)
		assert.Equal(t, "header1", state.CodeBlock.Header)
		assert.Equal(t, "echo 'hello world'", state.CodeBlock.Content)
		assert.Equal(t, false, state.Success)
	})

	t.Run(
		"Initializing the test model invokes the first command to start running tests.",
		func(t *testing.T) {
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
			assert.NoError(t, err)

			m, _ := model.Update(model.Init()())

			if model, ok := m.(TestModeModel); ok {
				assert.Equal(t, 1, model.currentCodeBlock)
			}
		},
	)
}
