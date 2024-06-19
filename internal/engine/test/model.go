package test

import (
	"fmt"
	"strings"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/patterns"
	"github.com/Azure/InnovationEngine/internal/shells"
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
	codeBlockState       map[int]common.StatefulCodeBlock
	commands             TestModeCommands
	currentCodeBlock     int
	environmentVariables map[string]string
	environment          string
	help                 help.Model
	resourceGroupName    string
	scenarioTitle        string
	scenarioCompleted    bool
	components           testModeComponents
	ready                bool
	CommandLines         []string
}

// Init the test mode model by executing the first code block.
func (model TestModeModel) Init() tea.Cmd {
	return common.ExecuteCodeBlockAsync(
		model.codeBlockState[model.currentCodeBlock].CodeBlock,
		model.environmentVariables,
	)
}

// Update the test mode model.
func (model TestModeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	viewportContentUpdated := false

	switch message := message.(type) {

	case tea.WindowSizeMsg:
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
		model.CommandLines = append(model.CommandLines, codeBlockState.StdOut)
		viewportContentUpdated = true

		// Increment the codeblock and update the viewport content.
		model.currentCodeBlock++

		if model.currentCodeBlock < len(model.codeBlockState) {
			nextCommand := model.codeBlockState[model.currentCodeBlock].CodeBlock.Content
			nextLanguage := model.codeBlockState[model.currentCodeBlock].CodeBlock.Language

			model.CommandLines = append(model.CommandLines, ui.CommandPrompt(nextLanguage)+nextCommand)
		}

		// Only increment the step for azure if the step name has changed.
		nextCodeBlockState := model.codeBlockState[model.currentCodeBlock]

		// If the scenario has been completed, we need to update the azure
		// status and quit the program. else,
		if model.currentCodeBlock == len(model.codeBlockState) {
			logging.GlobalLogger.Infof("The last codeblock was executed. Requesting to exit test mode...")
			commands = append(
				commands,
				common.Exit(false),
			)

		} else {
			// If the scenario has not been completed, we need to execute the next command
			commands = append(
				commands,
				common.ExecuteCodeBlockAsync(nextCodeBlockState.CodeBlock, model.environmentVariables),
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
		model.CommandLines = append(model.CommandLines, codeBlockState.StdErr+message.Error.Error())
		viewportContentUpdated = true

		commands = append(commands, common.Exit(true))

	case common.ExitMessage:
		// TODO: Generate test report

		// Delete any found resource groups.
		if model.resourceGroupName != "" {
			logging.GlobalLogger.Infof("Attempting to delete the deployed resource group with the name: %s", model.resourceGroupName)
			command := fmt.Sprintf("az group delete --name %s --yes --no-wait", model.resourceGroupName)
			_, err := shells.ExecuteBashCommand(
				command,
				shells.BashCommandConfiguration{
					EnvironmentVariables: lib.CopyMap(model.environmentVariables),
					InheritEnvironment:   true,
					InteractiveCommand:   false,
					WriteToHistory:       true,
				},
			)
			if err != nil {
				model.CommandLines = append(model.CommandLines, ui.ErrorStyle.Render("Error deleting resource group: %s\n", err.Error()))
				logging.GlobalLogger.Errorf("Error deleting resource group: %s", err.Error())
			} else {
				model.CommandLines = append(model.CommandLines, "Resource group deleted successfully.")
			}

		}

		// If the model didn't encounter a failure, then the scenario was scenario
		// was completed successfully.
		model.scenarioCompleted = !message.EncounteredFailure

		commands = append(commands, tea.Quit)

	}

	model.components.commandViewport.SetContent(strings.Join(model.CommandLines, "\n"))

	if viewportContentUpdated {
		model.components.commandViewport.GotoBottom()
	}

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

	// If the environment variables are not set, set it to an empty map.
	if len(env) == 0 || env == nil {
		env = make(map[string]string)
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
		environmentVariables: env,
		resourceGroupName:    "",
		codeBlockState:       codeBlockState,
		currentCodeBlock:     0,
		help:                 help.New(),
		environment:          environment,
		scenarioCompleted:    false,
		ready:                false,
		CommandLines:         commandLines,
	}, nil
}
