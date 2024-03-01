package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/parsers"
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

type InteractiveModeCommands struct {
	execute  key.Binding
	quit     key.Binding
	previous key.Binding
	next     key.Binding
}

// State for the codeblock in interactive mode. Used to keep track of the
// state of each codeblock.
type CodeBlockState struct {
	CodeBlock       parsers.CodeBlock
	CodeBlockNumber int
	Error           error
	StdErr          string
	StdOut          string
	StepName        string
	StepNumber      int
	Success         bool
}

type interactiveModeComponents struct {
	paginator        paginator.Model
	stepViewport     viewport.Model
	outputViewport   viewport.Model
	azureCLIViewport viewport.Model
}

type InteractiveModeModel struct {
	azureStatus       environments.AzureDeploymentStatus
	codeBlockState    map[int]CodeBlockState
	commands          InteractiveModeCommands
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
	components        interactiveModeComponents
	ready             bool
	commandLines      []string
}

// Initialize the intractive mode model
func (model InteractiveModeModel) Init() tea.Cmd {
	environments.ReportAzureStatus(model.azureStatus, model.environment)
	return tea.Batch(clearScreen(), tea.Tick(time.Millisecond*10, func(t time.Time) tea.Msg {
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

			commands = append(commands, tea.Sequence(
				updateAzureStatus(model),
				func() tea.Msg {
					return ExecuteCodeBlockSync(codeBlock, lib.CopyMap(model.env))
				}))

		} else {
			commands = append(commands, ExecuteCodeBlockAsync(
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

	case SuccessfulCommandMessage:
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
		nextCommand := model.codeBlockState[model.currentCodeBlock].CodeBlock.Content
		nextLanguage := model.codeBlockState[model.currentCodeBlock].CodeBlock.Language
		model.commandLines = append(model.commandLines, ui.CommandPrompt(nextLanguage)+nextCommand)

		// Only increment the step for azure if the step name has changed.
		nextCodeBlockState := model.codeBlockState[model.currentCodeBlock]

		if codeBlockState.StepName != nextCodeBlockState.StepName {
			logging.GlobalLogger.Debugf("Step name has changed, incrementing step for Azure")
			model.azureStatus.CurrentStep++
		} else {
			logging.GlobalLogger.Debugf("Step name has not changed, not incrementing step for Azure")
		}

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
			commands = append(
				commands,
				tea.Sequence(
					updateAzureStatus(model),
					tea.Quit,
				),
			)
		} else {
			commands = append(commands, updateAzureStatus(model))
		}

	case FailedCommandMessage:
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
		model.azureStatus.SetError(message.Error)
		environments.AttachResourceURIsToAzureStatus(
			&model.azureStatus,
			model.resourceGroupName,
			model.environment,
		)
		commands = append(commands, tea.Sequence(updateAzureStatus(model), tea.Quit))

	case AzureStatusUpdatedMessage:
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

	model.components.azureCLIViewport.SetContent(strings.Join(model.commandLines, "\n"))

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
	if model.environment == "azure" {
		return model.components.azureCLIViewport.View()
	}

	scenarioTitle := ui.ScenarioTitleStyle.Width(model.width).
		Align(lipgloss.Center).
		Render(model.scenarioTitle)
	var stepTitle string
	var stepView string
	var stepSection string
	stepTitle = ui.StepTitleStyle.Render(
		fmt.Sprintf(
			"Step %d - %s",
			model.currentCodeBlock+1,
			model.codeBlockState[model.currentCodeBlock].StepName,
		),
	)

	border := lipgloss.NewStyle().
		Width(model.components.stepViewport.Width - 2).
		Border(lipgloss.NormalBorder())

	stepView = border.Render(model.components.stepViewport.View())

	stepSection = fmt.Sprintf("%s\n%s\n\n", stepTitle, stepView)

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

// TODO: Ideally we won't need a global program variable. We should
// refactor this in the future such that each tea program is localized to the
// function that creates it and ExecuteCodeBlockSync doesn't mutate the global
// program variable.
var program *tea.Program = nil

// Create a new interactive mode model.
func NewInteractiveModeModel(
	title string,
	engine *Engine,
	steps []Step,
	env map[string]string,
) (InteractiveModeModel, error) {
	// TODO: In the future we should just set the current step for the azure status
	// to one as the default.
	azureStatus := environments.NewAzureDeploymentStatus()
	azureStatus.CurrentStep = 1
	totalCodeBlocks := 0
	codeBlockState := make(map[int]CodeBlockState)

	err := az.SetSubscription(engine.Configuration.Subscription)
	if err != nil {
		logging.GlobalLogger.Errorf("Invalid Config: Failed to set subscription: %s", err)
		azureStatus.SetError(err)
		environments.ReportAzureStatus(azureStatus, engine.Configuration.Environment)
		return InteractiveModeModel{}, err
	}

	for stepNumber, step := range steps {
		azureCodeBlocks := []environments.AzureCodeBlock{}
		for blockNumber, block := range step.CodeBlocks {
			azureCodeBlocks = append(azureCodeBlocks, environments.AzureCodeBlock{
				Command:     block.Content,
				Description: block.Description,
			})

			codeBlockState[totalCodeBlocks] = CodeBlockState{
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
		},
		env:               env,
		resourceGroupName: "",
		azureStatus:       azureStatus,
		codeBlockState:    codeBlockState,
		executingCommand:  false,
		currentCodeBlock:  0,
		help:              help.New(),
		environment:       engine.Configuration.Environment,
		scenarioCompleted: false,
		ready:             false,
		commandLines:      commandLines,
	}, nil
}
