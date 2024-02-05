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
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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

type interactiveModeViewPorts struct {
	step   viewport.Model
	output viewport.Model
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
	viewports         interactiveModeViewPorts
	ready             bool
}

// Initialize the intractive mode model
func (model InteractiveModeModel) Init() tea.Cmd {
	environments.ReportAzureStatus(model.azureStatus, model.environment)
	return tea.Batch(clearScreen(), tea.Tick(time.Millisecond*10, func(t time.Time) tea.Msg {
		return tea.KeyMsg{Type: tea.KeyCtrlL} // This is to force a repaint
	}))
}

// Initializes the viewports for the interactive mode model.
func initializeViewports(model InteractiveModeModel, width, height int) interactiveModeViewPorts {
	currentBlock := model.codeBlockState[model.currentCodeBlock]

	stepViewport := viewport.New(width, 8)
	stepViewport.SetContent(currentBlock.CodeBlock.Description)

	// Initialize the output view ports

	outputViewport := viewport.New(width, 4)
	outputViewport.SetContent(currentBlock.StdOut)

	return interactiveModeViewPorts{
		step:   stepViewport,
		output: outputViewport,
	}
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

// Updates the intractive mode model
func (model InteractiveModeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch message := message.(type) {

	case tea.WindowSizeMsg:
		model.width = message.Width
		model.height = message.Height
		logging.GlobalLogger.Debugf("Window size changed to: %d x %d", message.Width, message.Height)
		if !model.ready {
			model.viewports = initializeViewports(model, message.Width, message.Height)
			model.ready = true
		} else {
			model.viewports.step.Width = message.Width
			model.viewports.output.Width = message.Width
		}
		commands = append(commands, clearScreen())

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

		// Increment the codeblock and update the viewport content.
		model.currentCodeBlock++

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
		commands = append(commands, func() tea.Msg {
			return tea.WindowSizeMsg{
				Width:  model.width,
				Height: model.height,
			}
		})
	}

	// Update viewport content
	block := model.codeBlockState[model.currentCodeBlock]
	model.viewports.step.SetContent(block.CodeBlock.Description + "\n\n" + block.CodeBlock.Content)

	if block.Success {
		model.viewports.output.SetContent(block.StdOut)
	} else {
		model.viewports.output.SetContent(block.StdErr)
	}

	// Update all the viewports and append resulting commands.
	var command tea.Cmd
	model.viewports.step, command = model.viewports.step.Update(message)
	commands = append(commands, command)
	model.viewports.output, command = model.viewports.output.Update(message)
	commands = append(commands, command)

	return model, tea.Batch(commands...)
}

// Shows the commands that the user can use to interact with the interactive
// mode model.
func (model InteractiveModeModel) helpView() string {
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
	stepName := model.codeBlockState[model.currentCodeBlock].StepName
	stepView := fmt.Sprintf("%s\n%s\n%s",
		viewportHeaderView(
			fmt.Sprintf("Step %d: %s", model.currentCodeBlock+1, stepName),
			model.width,
		),
		model.viewports.step.View(),
		viewportFooterView(
			fmt.Sprintf("%3.f%%", model.viewports.step.ScrollPercent()*100),
			model.width,
		),
	)

	outputView := fmt.Sprintf("%s\n%s\n%s",
		viewportHeaderView(
			"Output",
			model.width,
		),
		model.viewports.output.View(),
		viewportFooterView(
			fmt.Sprintf("%3.f%%", model.viewports.output.ScrollPercent()*100),
			model.width,
		),
	)

	return stepView + "\n" + outputView + "\n" + model.helpView()
}

// Renders the header for a viewport.
func viewportHeaderView(header string, viewportWidth int) string {
	title := ui.InteractiveModeStepTitleStyle.Render(header)
	line := strings.Repeat("-", lib.Max(0, viewportWidth-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

// Renders the footer for a viewport.
func viewportFooterView(footer string, viewportWidth int) string {
	footer = ui.InteractiveModeStepFooterStyle.Render(footer)
	line := strings.Repeat("-", lib.Max(0, viewportWidth-lipgloss.Width(footer)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, footer)
}

// TODO: Ideally we won't need a global program variable. We should
// refactor this in the future such that each tea program is localized to the
// function that creates it and ExecuteCodeBlockSync doesn't mutate the global
// program variable.
var program *tea.Program = nil

// Create a new interactive mode model.
func NewInteractiveModeModel(
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

	return InteractiveModeModel{
		scenarioTitle: "Test",
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
	}, nil
}
