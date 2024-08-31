package interactive

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/patterns"
	"github.com/Azure/InnovationEngine/internal/ui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type InteractiveModeModel struct {
	azureStatus       environments.AzureDeploymentStatus
	codeBlockState    map[int]common.StatefulCodeBlock
	commands          InteractiveModeCommands
	currentCodeBlock  int
	env               map[string]string
	environment       string
	executingCommand  bool
	stepsToBeExecuted int
	recordingInput    bool
	recordedInput     string
	height            int
	help              help.Model
	resourceGroupName string
	subscription      string
	scenarioTitle     string
	width             int
	scenarioCompleted bool
	components        interactiveModeComponents
	ready             bool
	CommandLines      []string
}

// Initialize the intractive mode model
func (model InteractiveModeModel) Init() tea.Cmd {
	environments.ReportAzureStatus(model.azureStatus, model.environment)
	return tea.Batch(common.ClearScreen(), tea.Tick(time.Millisecond*10, func(t time.Time) tea.Msg {
		return tea.KeyMsg{Type: tea.KeyCtrlL} // This is to force a repaint
	}))
}

// Updates the intractive mode model
func (model InteractiveModeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
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
			model.components.stepViewport.Width = message.Width
			model.components.outputViewport.Width = message.Width
			model.components.azureCLIViewport.Width = message.Width
			model.components.updateViewportHeight(message.Height)
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
		codeBlockState.Status = common.STATUS_SUCCESS

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
				model.azureStatus.AddResourceURI(az.BuildResourceGroupId(model.subscription, model.resourceGroupName))
			}
		}
		model.CommandLines = append(model.CommandLines, codeBlockState.StdOut)

		// Increment the codeblock and update the viewport content.
		model.currentCodeBlock++

		if model.currentCodeBlock < len(model.codeBlockState) {
			nextCommand := model.codeBlockState[model.currentCodeBlock].CodeBlock.Content
			nextLanguage := model.codeBlockState[model.currentCodeBlock].CodeBlock.Language

			model.CommandLines = append(model.CommandLines, ui.CommandPrompt(nextLanguage)+nextCommand)
		}

		// Only increment the step for azure if the step name has changed.
		nextCodeBlockState := model.codeBlockState[model.currentCodeBlock]

		if codeBlockState.StepName != nextCodeBlockState.StepName {
			logging.GlobalLogger.Debugf("Step name has changed, incrementing step for Azure")
			model.azureStatus.CurrentStep++
		} else {
			logging.GlobalLogger.Debugf("Step name has not changed, not incrementing step for Azure")
		}

		model.stepsToBeExecuted--

		// If the scenario has been completed, we need to update the azure
		// status and quit the program.
		if model.currentCodeBlock == len(model.codeBlockState) {
			model.scenarioCompleted = true
			model.azureStatus.Status = "Succeeded"
			environments.AttachResourceURIsToAzureStatus(
				&model.azureStatus,
				model.resourceGroupName,
				model.environment,
			)
			model.azureStatus.SetOutput(strings.Join(model.CommandLines, "\n"))
			commands = append(
				commands,
				tea.Sequence(
					common.UpdateAzureStatus(model.azureStatus, model.environment),
					tea.Quit,
				),
			)
		} else {
			commands = append(
				commands,
				tea.Sequence(
					common.UpdateAzureStatus(model.azureStatus, model.environment),
					// Send a key event to trigger
					func() tea.Msg {
						if model.stepsToBeExecuted <= 0 {
							return nil
						}
						return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
					},
				),
			)
		}

	case common.FailedCommandMessage:
		// Handle failed command executions

		// Update the state of the codeblock which finished executing.
		step := model.currentCodeBlock
		codeBlockState := model.codeBlockState[step]
		codeBlockState.StdOut = message.StdOut
		codeBlockState.StdErr = message.StdErr
		codeBlockState.Status = common.STATUS_FAILURE

		model.codeBlockState[step] = codeBlockState
		model.CommandLines = append(model.CommandLines, codeBlockState.StdErr)

		// Report the error
		model.executingCommand = false
		model.azureStatus.SetError(message.Error)
		environments.AttachResourceURIsToAzureStatus(
			&model.azureStatus,
			model.resourceGroupName,
			model.environment,
		)

		model.azureStatus.SetOutput(strings.Join(model.CommandLines, "\n"))
		commands = append(
			commands,
			tea.Sequence(
				common.UpdateAzureStatus(model.azureStatus, model.environment),
				tea.Quit,
			),
		)

	case common.AzureStatusUpdatedMessage:
		// After the status has been updated, we force a window resize to
		// render over the status update. For some reason, clearing the screen
		// manually seems to cause the text produced by View() to not render
		// properly if we don't trigger a window size event.
		commands = append(commands,
			tea.Sequence(
				tea.ClearScreen,
				func() tea.Msg {
					return tea.WindowSizeMsg{
						Width:  model.width,
						Height: model.height,
					}
				},
			),
		)
	}

	// Update viewport content
	block := model.codeBlockState[model.currentCodeBlock]

	renderedStepSection := fmt.Sprintf(
		"%s\n\n%s",
		block.CodeBlock.Description,
		block.CodeBlock.Content,
	)

	// TODO(vmarcella): We shoulkd figure out a way to not have to recreate
	// the renderer every time we update the view.
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(model.width-4),
	)

	if err == nil {
		var glamourizedSection string
		glamourizedSection, err = renderer.Render(
			fmt.Sprintf(
				"%s\n```%s\n%s```",
				block.CodeBlock.Description,
				block.CodeBlock.Language,
				block.CodeBlock.Content,
			),
		)
		if err != nil {
			logging.GlobalLogger.Errorf(
				"Error rendering codeblock: %s, using non rendered codeblock instead",
				err,
			)
		} else {
			renderedStepSection = glamourizedSection
		}
	} else {
		logging.GlobalLogger.Errorf(
			"Error creating glamour renderer: %s, using non rendered codeblock instead",
			err,
		)
	}

	model.components.stepViewport.SetContent(
		renderedStepSection,
	)

	if block.WasExecuted() {
		if block.Succeeded() {
			model.components.outputViewport.SetContent(block.StdOut)
		} else {
			model.components.outputViewport.SetContent(block.StdErr)
		}
	} else {
		model.components.outputViewport.SetContent("")
	}

	model.components.azureCLIViewport.SetContent(strings.Join(model.CommandLines, "\n"))

	// Update all the viewports and append resulting commands.
	updatedComponents, componentCommands := updateComponents(
		model.components,
		model.currentCodeBlock,
		message,
	)
	commands = append(commands, componentCommands...)
	model.components = updatedComponents

	return model, tea.Batch(commands...)
}

// Shows the commands that the user can use to interact with the interactive
// mode model.
func (model InteractiveModeModel) helpView() string {
	if model.environment == "azure" {
		return ""
	}
	keyBindingGroups := [][]key.Binding{
		// Command related bindings
		{
			model.commands.execute,
			model.commands.executeAll,
			model.commands.executeMany,
			model.commands.previous,
			model.commands.next,
		},
		// Scenario related bindings
		{
			model.commands.quit,
		},
	}

	return model.help.FullHelpView(keyBindingGroups)
}

// Renders the interactive mode model.
func (model InteractiveModeModel) View() string {
	// When running in the portal, we only want to show the Azure CLI viewport
	// which mimics a command line interface during execution.
	if model.environment == "azure" {
		return model.components.azureCLIViewport.View()
	}

	scenarioTitle := ui.ScenarioTitleStyle.Width(model.width).
		Align(lipgloss.Center).
		Render(model.scenarioTitle)

	border := lipgloss.NewStyle().
		Width(model.components.stepViewport.Width - 2).
		Border(lipgloss.NormalBorder())

	stepTitle := ui.StepTitleStyle.Render(
		fmt.Sprintf(
			"Step %d - %s",
			model.currentCodeBlock+1,
			model.codeBlockState[model.currentCodeBlock].StepName,
		),
	)
	stepView := border.Render(model.components.stepViewport.View())
	stepSection := fmt.Sprintf("%s\n%s\n\n", stepTitle, stepView)

	outputTitle := ui.StepTitleStyle.Render("Output")
	outputView := border.Render(model.components.outputViewport.View())
	outputSection := fmt.Sprintf("%s\n%s\n\n", outputTitle, outputView)

	paginator := lipgloss.NewStyle().
		Width(model.width).
		Align(lipgloss.Center).
		Render(model.components.paginator.View())

	var executing string

	if model.executingCommand {
		executing = "Executing command..."
	} else {
		executing = ""
	}

	// TODO(vmarcella): Format this to be more readable.
	return ((scenarioTitle + "\n") +
		(paginator + "\n\n") +
		(stepSection) +
		(outputSection) +
		(model.helpView())) +
		("\n" + executing)
}

// Create a new interactive mode model.
func NewInteractiveModeModel(
	title string,
	subscription string,
	environment string,
	steps []common.Step,
	env map[string]string,
) (InteractiveModeModel, error) {
	// TODO: In the future we should just set the current step for the azure status
	// to one as the default.
	azureStatus := environments.NewAzureDeploymentStatus()
	azureStatus.CurrentStep = 1
	totalCodeBlocks := 0
	codeBlockState := make(map[int]common.StatefulCodeBlock)

	err := az.SetSubscription(subscription)
	if err != nil {
		logging.GlobalLogger.Errorf("Invalid Config: Failed to set subscription: %s", err)
		azureStatus.SetError(err)
		environments.ReportAzureStatus(azureStatus, environment)
		return InteractiveModeModel{}, err
	}

	for stepNumber, step := range steps {
		azureCodeBlocks := []environments.AzureCodeBlock{}
		for blockNumber, block := range step.CodeBlocks {
			azureCodeBlocks = append(azureCodeBlocks, environments.AzureCodeBlock{
				Command:     block.Content,
				Description: block.Description,
			})

			codeBlockState[totalCodeBlocks] = common.StatefulCodeBlock{
				StepName:        step.Name,
				CodeBlock:       block,
				StepNumber:      stepNumber,
				CodeBlockNumber: blockNumber,
				StdOut:          "",
				StdErr:          "",
				Error:           nil,
				Status:          common.STATUS_PENDING,
			}

			totalCodeBlocks += 1
		}
		azureStatus.AddStep(fmt.Sprintf("%d. %s", stepNumber+1, step.Name), azureCodeBlocks)
	}

	language := codeBlockState[0].CodeBlock.Language
	commandLines := []string{
		ui.CommandPrompt(language) + codeBlockState[0].CodeBlock.Content,
	}

	return InteractiveModeModel{
		scenarioTitle:     title,
		commands:          NewInteractiveModeCommands(),
		stepsToBeExecuted: 0,
		env:               env,
		subscription:      subscription,
		resourceGroupName: "",
		azureStatus:       azureStatus,
		codeBlockState:    codeBlockState,
		executingCommand:  false,
		currentCodeBlock:  0,
		help:              help.New(),
		environment:       environment,
		scenarioCompleted: false,
		ready:             false,
		CommandLines:      commandLines,
	}, nil
}
