package test

import (
	"testing"

	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/shells"
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
		assert.Equal(t, map[string]string{}, model.environmentVariables)
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
		assert.Equal(t, true, state.Succeeded())
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

				executedBlock := model.codeBlockState[0]

				// Assert outputs of the executed block.
				assert.Equal(t, "hello world\n", executedBlock.StdOut)
				assert.Equal(t, "", executedBlock.StdErr)
				assert.Equal(t, true, executedBlock.Succeeded())
			}
		},
	)

	t.Run(
		"Test mode doesn't try to delete resource group if none was created.",
		func(t *testing.T) {
			steps := []common.Step{
				{
					Name: "step1",
					CodeBlocks: []parsers.CodeBlock{
						{
							Header:  "header1",
							Content: "echo 'hello world'",
						},
					},
				},
			}

			model, err := NewTestModeModel("test", "", "test", steps, nil)

			assert.NoError(t, err)

			m, _ := model.Update(model.Init()())

			if model, ok := m.(TestModeModel); ok {
				assert.Equal(t, 1, model.currentCodeBlock)

				executedBlock := model.codeBlockState[0]
				model.resourceGroupName = "test"

				// Assert outputs of the executed block.
				assert.Equal(t, "hello world\n", executedBlock.StdOut)
				assert.Equal(t, "", executedBlock.StdErr)
				assert.Equal(t, true, executedBlock.Succeeded())
			} else {
				assert.Fail(t, "Model is not a TestModeModel")
			}

			// Assert that the model doesn't try to delete the resource group when
			// the resource group name is empty.
			m, _ = model.Update(common.Exit(false)())
			counter := 0

			// We create a mock function to replace the shells.ExecuteBashCommand function
			// to make sure that the function is not called.
			original := shells.ExecuteBashCommand
			defer func() { shells.ExecuteBashCommand = original }()

			shells.ExecuteBashCommand = func(
				command string,
				config shells.BashCommandConfiguration,
			) (shells.CommandOutput, error) {
				counter += 1
				return shells.CommandOutput{}, nil
			}

			if model, ok := m.(TestModeModel); ok {
				assert.Equal(t, 0, counter)
				assert.Equal(t, true, model.scenarioCompleted)
			} else {
				assert.Fail(t, "Model is not a TestModeModel")
			}
		},
	)

	t.Run(
		"Test mode tries to delete resource group if one was created.",
		func(t *testing.T) {
			steps := []common.Step{
				{
					Name: "step1",
					CodeBlocks: []parsers.CodeBlock{
						{
							Header:  "header1",
							Content: "echo 'hello world'",
						},
					},
				},
			}

			model, err := NewTestModeModel("test", "", "test", steps, nil)

			assert.NoError(t, err)

			m, _ := model.Update(model.Init()())

			var ok bool

			if model, ok = m.(TestModeModel); ok {
				assert.Equal(t, 1, model.currentCodeBlock)

				executedBlock := model.codeBlockState[0]
				model.resourceGroupName = "test"

				// Assert outputs of the executed block.
				assert.Equal(t, "hello world\n", executedBlock.StdOut)
				assert.Equal(t, "", executedBlock.StdErr)
				assert.Equal(t, true, executedBlock.Succeeded())

			} else {
				assert.Fail(t, "Model is not a TestModeModel")
			}

			// Assert that the model tries to delete the resource group when
			// the resource group name is not empty.
			counter := 0
			recordedCommand := ""

			// We create a mock function to replace the shells.ExecuteBashCommand function
			// to make sure that the function is called.
			original := shells.ExecuteBashCommand
			defer func() { shells.ExecuteBashCommand = original }()

			shells.ExecuteBashCommand = func(
				command string,
				config shells.BashCommandConfiguration,
			) (shells.CommandOutput, error) {
				recordedCommand = command
				counter += 1
				return shells.CommandOutput{}, nil
			}

			m, _ = model.Update(common.Exit(false)())

			if model, ok = m.(TestModeModel); ok {
				assert.Equal(t, 1, counter)
				assert.Equal(t, "az group delete --name test --yes --no-wait", recordedCommand)
				assert.Equal(t, true, model.scenarioCompleted)
			} else {
				assert.Fail(t, "Model is not a TestModeModel")
			}
		},
	)
}
