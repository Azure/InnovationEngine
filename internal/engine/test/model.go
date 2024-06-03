package test

import (
	"strings"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/patterns"
	"github.com/Azure/InnovationEngine/internal/ui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Commands accessible to the user for test mode.
type TestModeCommands struct {
	quit key.Binding
}

// The state required for testing scenarios.
type TestModeModel struct {
	codeBlockState    map[int]common.StatefulCodeBlock
	commands          TestModeCommands
	currentCodeBlock  int
	env               map[string]string
	environment       string
	executingCommand  bool
	height            int
	help              help.Model
	resourceGroupName string
	scenarioTitle     string
	width             int
	scenarioCompleted bool
	components        testModeComponents
	ready             bool
	commandLines      []string
}

// Init the test mode model.
func (model TestModeModel) Init() tea.Cmd {
	return common.ExecuteCodeBlockAsync(
		model.codeBlockState[model.currentCodeBlock].CodeBlock,
		model.env,
	)
}

// Update the test mode model.
func (model TestModeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch message := message.(type) {

	case tea.WindowSizeMsg:
		model.width = message.Width
		model.height = message.Height
		logging.GlobalLogger.Debugf("Window size changed to: %d x %d", message.Width, message.Height)
		if !model.ready {
			model.components = initializeComponents(model, message.Width, message.Height)
			model.ready = true
		} else {
			model.components.updateViewportSizing(message.Width, message.Height)
		}

	case tea.KeyMsg:
		model, commands = handleUserInput(model, message)

	case common.SuccessfulCommandMessage:
		// Handle successful command executions
		model.executingCommand = false
		step := model.currentCodeBlock

		// Update the state of the codeblock which finished executing.
		codeBlockState := model.codeBlockState[step]
		codeBlockState.StdOut = message.StdOut
		codeBlockState.StdErr = message.StdErr
		codeBlockState.Success = true
		model.codeBlockState[step] = codeBlockState

		logging.GlobalLogger.Infof("Finished executing:\n %s", codeBlockState.CodeBlock.Content)

		// Extract the resource group name from the command output if
		// it's not already set.
		if model.resourceGroupName == "" && patterns.AzCommand.MatchString(codeBlockState.CodeBlock.Content) {
			logging.GlobalLogger.Debugf("Attempting to extract resource group name from command output")
			tmpResourceGroup := az.FindResourceGroupName(codeBlockState.StdOut)
			if tmpResourceGroup != "" {
				logging.GlobalLogger.Infof("Found resource group named: %s", tmpResourceGroup)
				model.resourceGroupName = tmpResourceGroup
			}
		}
		model.commandLines = append(model.commandLines, codeBlockState.StdOut)

		// Increment the codeblock and update the viewport content.
		model.currentCodeBlock++

		if model.currentCodeBlock < len(model.codeBlockState) {
			nextCommand := model.codeBlockState[model.currentCodeBlock].CodeBlock.Content
			nextLanguage := model.codeBlockState[model.currentCodeBlock].CodeBlock.Language

			model.commandLines = append(model.commandLines, ui.CommandPrompt(nextLanguage)+nextCommand)
		}

		// Only increment the step for azure if the step name has changed.
		nextCodeBlockState := model.codeBlockState[model.currentCodeBlock]

		if codeBlockState.StepName != nextCodeBlockState.StepName {
			logging.GlobalLogger.Debugf("Step name has changed, incrementing step for Azure")
		} else {
			logging.GlobalLogger.Debugf("Step name has not changed, not incrementing step for Azure")
		}

		// If the scenario has been completed, we need to update the azure
		// status and quit the program. else,
		if model.currentCodeBlock == len(model.codeBlockState) {
			model.scenarioCompleted = true
			commands = append(
				commands,
				tea.Quit,
			)
		} else {
			// If the scenario has not been completed, we need to execute the next command
			commands = append(
				commands,
				common.ExecuteCodeBlockAsync(nextCodeBlockState.CodeBlock, model.env),
			)
		}

	case common.FailedCommandMessage:
		// Handle failed command executions

		// Update the state of the codeblock which finished executing.
		step := model.currentCodeBlock
		codeBlockState := model.codeBlockState[step]
		codeBlockState.StdOut = message.StdOut
		codeBlockState.StdErr = message.StdErr
		codeBlockState.Success = false

		model.codeBlockState[step] = codeBlockState
		model.commandLines = append(model.commandLines, codeBlockState.StdErr)

		// Report the error
		model.executingCommand = false
		commands = append(commands, tea.Quit)
	}

	model.components.commandViewport.SetContent(strings.Join(model.commandLines, "\n"))

	// Update all the viewports and append resulting commands.
	var command tea.Cmd

	model.components.commandViewport, command = model.components.commandViewport.Update(message)
	commands = append(commands, command)

	return model, tea.Batch(commands...)
}

// View the test mode model.
func (model TestModeModel) View() string {
	return model.components.commandViewport.View()
}

// Create a new test mode model.
func NewTestModeModel(
	title string,
	subscription string,
	environment string,
	steps []common.Step,
	env map[string]string,
) (TestModeModel, error) {
	totalCodeBlocks := 0
	codeBlockState := make(map[int]common.StatefulCodeBlock)

	err := az.SetSubscription(subscription)
	if err != nil {
		logging.GlobalLogger.Errorf("Invalid Config: Failed to set subscription: %s", err)
		return TestModeModel{}, err
	}

	// TODO(vmarcella): The codeblock state building should be reused across
	// Interactive mode and test mode in the future.
	for stepNumber, step := range steps {
		for blockNumber, block := range step.CodeBlocks {

			codeBlockState[totalCodeBlocks] = common.StatefulCodeBlock{
				StepName:        step.Name,
				CodeBlock:       block,
				StepNumber:      stepNumber,
				CodeBlockNumber: blockNumber,
				StdOut:          "",
				StdErr:          "",
				Error:           nil,
				Success:         false,
			}

			totalCodeBlocks += 1
		}
	}

	language := codeBlockState[0].CodeBlock.Language
	commandLines := []string{
		ui.CommandPrompt(language) + codeBlockState[0].CodeBlock.Content,
	}

	return TestModeModel{
		scenarioTitle: title,
		commands: TestModeCommands{
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "Quit the scenario."),
			),
		},
		env:               env,
		resourceGroupName: "",
		codeBlockState:    codeBlockState,
		executingCommand:  false,
		currentCodeBlock:  0,
		help:              help.New(),
		environment:       environment,
		scenarioCompleted: false,
		ready:             false,
		commandLines:      commandLines,
	}, nil
}
