package interactive

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/patterns"
	"github.com/Azure/InnovationEngine/internal/ui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// All interactive mode inputs.
type InteractiveModeCommands struct {
	execute     key.Binding
	executeAll  key.Binding
	executeMany key.Binding
	next        key.Binding
	pause       key.Binding
	previous    key.Binding
	quit        key.Binding
}

type interactiveModeComponents struct {
	paginator        paginator.Model
	stepViewport     viewport.Model
	outputViewport   viewport.Model
	azureCLIViewport viewport.Model
}

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
	markdownSource    string
	CommandLines      []string
}

// Initialize the intractive mode model
func (model InteractiveModeModel) Init() tea.Cmd {
	environments.ReportAzureStatus(model.azureStatus, model.environment)
	return tea.Batch(common.ClearScreen(), tea.Tick(time.Millisecond*10, func(t time.Time) tea.Msg {
		return tea.KeyMsg{Type: tea.KeyCtrlL} // This is to force a repaint
	}))
}

// Initializes the viewports for the interactive mode model.
func initializeComponents(model InteractiveModeModel, width, height int) interactiveModeComponents {
	// paginator setup
	p := paginator.New()
	p.TotalPages = len(model.codeBlockState)
	p.Type = paginator.Dots
	// Dots
	p.ActiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).
		Render("•")
	p.InactiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).
		Render("•")

	p.KeyMap.PrevPage = model.commands.previous
	p.KeyMap.NextPage = model.commands.next

	stepViewport := viewport.New(width, 4)
	outputViewport := viewport.New(width, 2)
	azureCLIViewport := viewport.New(width, height)

	components := interactiveModeComponents{
		paginator:        p,
		stepViewport:     stepViewport,
		outputViewport:   outputViewport,
		azureCLIViewport: azureCLIViewport,
	}

	components.updateViewportHeight(height)
	return components
}

// Handle user input for interactive mode.
func handleUserInput(
	model InteractiveModeModel,
	message tea.KeyMsg,
) (InteractiveModeModel, []tea.Cmd) {
	var commands []tea.Cmd

	// If we're recording input for a multi-char command,
	if model.recordingInput {
		isNumber := lib.IsNumber(message.String())

		// If the input is a number, append it to the recorded input.
		if message.Type == tea.KeyRunes && isNumber {
			model.recordedInput += message.String()
			return model, commands
		}

		// If the input is not a number, we'll stop recording input and reset
		// the commands remaining to the recorded input.
		if message.Type == tea.KeyEnter || !isNumber {
			commandsRemaining, _ := strconv.Atoi(model.recordedInput)

			if commandsRemaining > len(model.codeBlockState)-model.currentCodeBlock {
				commandsRemaining = len(model.codeBlockState) - model.currentCodeBlock
			}

			logging.GlobalLogger.Debugf("Will execute the next %d steps", commandsRemaining)
			model.stepsToBeExecuted = commandsRemaining
			commands = append(commands, func() tea.Msg {
				return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			})

			model.recordingInput = false
			model.recordedInput = ""
			logging.GlobalLogger.Debugf(
				"Recording input stopped and previously recorded input cleared.",
			)
			return model, commands
		}
	}

	switch {
	case key.Matches(message, model.commands.execute):
		if model.executingCommand {
			logging.GlobalLogger.Info("Command is already executing, ignoring execute command")
			break
		}

		// Prevent the user from executing a command if the previous command has
		// not been executed successfully or executed at all.
		previousCodeBlock := model.currentCodeBlock - 1
		if previousCodeBlock >= 0 {
			previousCodeBlockState := model.codeBlockState[previousCodeBlock]
			if !previousCodeBlockState.Success {
				logging.GlobalLogger.Info(
					"Previous command has not been executed successfully, ignoring execute command",
				)
				break
			}
		}

		// Prevent the user from executing a command if the current command has
		// already been executed successfully.
		codeBlockState := model.codeBlockState[model.currentCodeBlock]
		if codeBlockState.Success {
			logging.GlobalLogger.Info(
				"Command has already been executed successfully, ignoring execute command",
			)
			break
		}

		codeBlock := codeBlockState.CodeBlock

		model.executingCommand = true

		// If we're on the last step and the command is an SSH command, we need
		// to report the status before executing the command. This is needed for
		// one click deployments and does not affect the normal execution flow.
		if model.currentCodeBlock == len(model.codeBlockState)-1 &&
			patterns.SshCommand.MatchString(codeBlock.Content) {
			model.azureStatus.Status = "Succeeded"
			environments.AttachResourceURIsToAzureStatus(
				&model.azureStatus,
				model.resourceGroupName,
				model.environment,
			)

			environmentVariables, err := lib.LoadEnvironmentStateFile(
				lib.DefaultEnvironmentStateFile,
			)
			if err != nil {
				logging.GlobalLogger.Errorf("Failed to load environment state file: %s", err)
				model.azureStatus.SetError(err)
			}

			model.azureStatus.ConfigureMarkdownForDownload(
				model.markdownSource,
				environmentVariables,
				model.environment,
			)

			commands = append(commands, tea.Sequence(
				common.UpdateAzureStatus(model.azureStatus, model.environment),
				func() tea.Msg {
					return common.ExecuteCodeBlockSync(codeBlock, lib.CopyMap(model.env))
				}))

		} else {
			commands = append(commands, common.ExecuteCodeBlockAsync(
				codeBlock,
				lib.CopyMap(model.env),
			))
		}

	case key.Matches(message, model.commands.previous):
		if model.executingCommand {
			logging.GlobalLogger.Info("Command is already executing, ignoring execute command")
			break
		}
		if model.currentCodeBlock > 0 {
			model.currentCodeBlock--
		}
	case key.Matches(message, model.commands.next):
		if model.executingCommand {
			logging.GlobalLogger.Info("Command is already executing, ignoring execute command")
			break
		}
		if model.currentCodeBlock < len(model.codeBlockState)-1 {
			model.currentCodeBlock++
		}

	case key.Matches(message, model.commands.quit):
		commands = append(commands, tea.Quit)

	case key.Matches(message, model.commands.executeAll):
		model.stepsToBeExecuted = len(model.codeBlockState) - model.currentCodeBlock
		commands = append(
			commands,
			func() tea.Msg {
				return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			},
		)
	case key.Matches(message, model.commands.executeMany):
		model.recordingInput = true
	case key.Matches(message, model.commands.pause):
		if !model.executingCommand {
			logging.GlobalLogger.Info("No command is currently executing, ignoring pause command")
		}
		model.stepsToBeExecuted = 0
	}

	return model, commands
}

func (components *interactiveModeComponents) updateViewportHeight(terminalHeight int) {
	stepViewportPercent := 0.4
	outputViewportPercent := 0.2
	stepViewportHeight := int(float64(terminalHeight) * stepViewportPercent)
	outputViewportHeight := int(float64(terminalHeight) * outputViewportPercent)

	if stepViewportHeight < 4 {
		stepViewportHeight = 4
	}

	if outputViewportHeight < 2 {
		outputViewportHeight = 2
	}

	components.stepViewport.Height = stepViewportHeight
	components.outputViewport.Height = outputViewportHeight
	components.azureCLIViewport.Height = terminalHeight - 1
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

			environmentVariables, err := lib.LoadEnvironmentStateFile(lib.DefaultEnvironmentStateFile)
			if err != nil {
				logging.GlobalLogger.Errorf("Failed to load environment state file: %s", err)
				model.azureStatus.SetError(err)
			}

			model.azureStatus.ConfigureMarkdownForDownload(
				model.markdownSource,
				environmentVariables,
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
		codeBlockState.Success = false

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

	if block.Success {
		model.components.outputViewport.SetContent(block.StdOut)
	} else {
		model.components.outputViewport.SetContent(block.StdErr)
	}

	model.components.azureCLIViewport.SetContent(strings.Join(model.CommandLines, "\n"))

	// Update all the viewports and append resulting commands.
	var command tea.Cmd

	model.components.paginator.Page = model.currentCodeBlock

	model.components.stepViewport, command = model.components.stepViewport.Update(message)
	commands = append(commands, command)

	model.components.outputViewport, command = model.components.outputViewport.Update(message)
	commands = append(commands, command)

	model.components.azureCLIViewport, command = model.components.azureCLIViewport.Update(message)
	commands = append(commands, command)

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
	markdownSource string,
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
				Success:         false,
			}

			totalCodeBlocks += 1
		}
		azureStatus.AddStep(fmt.Sprintf("%d. %s", stepNumber+1, step.Name), azureCodeBlocks)
	}

	language := codeBlockState[0].CodeBlock.Language
	commandLines := []string{
		ui.CommandPrompt(language) + codeBlockState[0].CodeBlock.Content,
	}

	// Configure extra keybinds used for executing the many/all commands.
	executeAllKeybind := key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Execute all remaining commands."),
	)

	executeManyKeybind := key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m<number><enter>", "Execute the next <number> commands."),
	)
	pauseKeybind := key.NewBinding(
		key.WithKeys("p", "Pause execution of commands."),
	)

	return InteractiveModeModel{
		scenarioTitle: title,
		commands: InteractiveModeCommands{
			execute: key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "Execute the current command."),
			),
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "Quit the scenario."),
			),
			previous: key.NewBinding(
				key.WithKeys("left"),
				key.WithHelp("←", "Go to the previous command."),
			),
			next: key.NewBinding(
				key.WithKeys("right"),
				key.WithHelp("→", "Go to the next command."),
			),
			// Only enabled when in the azure environment.
			executeAll:  executeAllKeybind,
			executeMany: executeManyKeybind,
			pause:       pauseKeybind,
		},
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
		markdownSource:    markdownSource,
		CommandLines:      commandLines,
	}, nil
}
