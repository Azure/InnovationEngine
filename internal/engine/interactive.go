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
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/Azure/InnovationEngine/internal/ui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InteractiveModeCommands struct {
	execute key.Binding
	skip    key.Binding
	quit    key.Binding
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
	step    viewport.Model
	command viewport.Model
	output  viewport.Model
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

type SuccessfulCommandMessage struct {
	StdOut string
	StdErr string
}

type FailedCommandMessage struct {
	StdOut string
	StdErr string
	Error  error
}

// Executes a bash command and returns a tea message with the output. This function
// will be executed asycnhronously.
func ExecuteCodeBlockAsync(codeBlock parsers.CodeBlock, env map[string]string) tea.Cmd {
	return func() tea.Msg {
		logging.GlobalLogger.Info("Executing command asynchronously: ", codeBlock.Content)
		output, err := shells.ExecuteBashCommand(codeBlock.Content, shells.BashCommandConfiguration{
			EnvironmentVariables: env,
			InheritEnvironment:   true,
			InteractiveCommand:   false,
			WriteToHistory:       true,
		})

		if err != nil {
			logging.GlobalLogger.Errorf("Error executing command: %s", err.Error())
			return FailedCommandMessage{
				StdOut: output.StdOut,
				StdErr: output.StdErr,
				Error:  err,
			}
		}

		// Check command output against the expected output.
		actualOutput := output.StdOut
		expectedOutput := codeBlock.ExpectedOutput.Content
		expectedSimilarity := codeBlock.ExpectedOutput.ExpectedSimilarity
		expectedRegex := codeBlock.ExpectedOutput.ExpectedRegex
		expectedOutputLanguage := codeBlock.ExpectedOutput.Language

		outputComparisonError := compareCommandOutputs(
			actualOutput,
			expectedOutput,
			expectedSimilarity,
			expectedRegex,
			expectedOutputLanguage,
		)

		if outputComparisonError != nil {
			logging.GlobalLogger.Errorf(
				"Error comparing command outputs: %s",
				outputComparisonError.Error(),
			)

			return FailedCommandMessage{
				StdOut: output.StdOut,
				StdErr: output.StdErr,
				Error:  outputComparisonError,
			}

		}

		logging.GlobalLogger.Infof("Command output to stdout:\n %s", output.StdOut)
		return SuccessfulCommandMessage{
			StdOut: output.StdOut,
			StdErr: output.StdErr,
		}
	}
}

// Executes a bash command syncrhonously. This function will block until the command
// finishes executing.
func ExecuteCodeBlockSync(codeBlock parsers.CodeBlock, env map[string]string) tea.Msg {
	logging.GlobalLogger.Info("Executing command synchronously: ", codeBlock.Content)
	output, err := shells.ExecuteBashCommand(codeBlock.Content, shells.BashCommandConfiguration{
		EnvironmentVariables: env,
		InheritEnvironment:   true,
		InteractiveCommand:   true,
		WriteToHistory:       true,
	})

	if err != nil {
		return FailedCommandMessage{
			StdOut: output.StdOut,
			StdErr: output.StdErr,
			Error:  err,
		}
	}

	logging.GlobalLogger.Infof("Command output to stdout:\n %s", output.StdOut)
	return SuccessfulCommandMessage{
		StdOut: output.StdOut,
		StdErr: output.StdErr,
	}
}

// clearScreen returns a command that clears the terminal screen and positions the cursor at the top-left corner
func clearScreen() tea.Cmd {
	return func() tea.Msg {
		fmt.Print(
			"\033[H\033[2J",
		) // ANSI escape codes for clearing the screen and repositioning the cursor
		return nil
	}
}

// Updates the azure status with the current state of the interactive mode
// model.
func updateAzureStatus(model InteractiveModeModel) tea.Cmd {
	return func() tea.Msg {
		logging.GlobalLogger.Infof(
			"Attempting to update the azure status: %+v",
			model.azureStatus,
		)
		environments.ReportAzureStatus(model.azureStatus, model.environment)
		return AzureStatusUpdatedMessage{}
	}
}

// Empty struct used to indicate that the azure status has been updated so
// that we can respond to it within the Update() function.
type AzureStatusUpdatedMessage struct{}

// Initializes the viewports for the interactive mode model.
func initializeViewports(model InteractiveModeModel, width, height int) interactiveModeViewPorts {

	currentBlock := model.codeBlockState[model.currentCodeBlock]

	stepViewport := viewport.New(width, 8)
	stepViewport.SetContent(currentBlock.CodeBlock.Description)

	commandViewport := viewport.New(width, 6)
	commandViewport.SetContent(currentBlock.CodeBlock.Content)

	// Initialize the output view ports

	outputViewport := viewport.New(width, 4)
	outputViewport.SetContent(currentBlock.StdOut)

	return interactiveModeViewPorts{
		step:    stepViewport,
		command: commandViewport,
		output:  outputViewport,
	}
}

// Updates the intractive mode model
func (model InteractiveModeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	var commands []tea.Cmd

	switch message := message.(type) {

	case tea.WindowSizeMsg:
		model.width = message.Width
		model.height = message.Height
		if !model.ready {
			model.viewports = initializeViewports(model, message.Width, message.Height)
			model.ready = true
		} else {
			model.viewports.step.Width = message.Width
			model.viewports.command.Width = message.Width
			model.viewports.output.Width = message.Width
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(message, model.commands.execute):
			if model.executingCommand {
				logging.GlobalLogger.Info("Command is already executing, ignoring execute command")
				break
			}
			model.executingCommand = true
			codeBlock := model.codeBlockState[model.currentCodeBlock].CodeBlock

			// If we're on the last step and the command is an SSH command, we need
			// to report the status before executing the command. This is needed for
			// one click deployments and does not affect the normal execution flow.

			if model.currentCodeBlock == len(model.codeBlockState)-1 && patterns.SshCommand.MatchString(codeBlock.Content) {
				model.azureStatus.Status = "Succeeded"
				environments.AttachResourceURIsToAzureStatus(&model.azureStatus, model.resourceGroupName, model.environment)

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

		case key.Matches(message, model.commands.skip):
			model.currentCodeBlock++
		case key.Matches(message, model.commands.quit):
			commands = append(commands, tea.Quit)
		}

	case SuccessfulCommandMessage:
		// Handle successful command executions
		model.executingCommand = false
		step := model.currentCodeBlock

		// Update the state of the codeblock which finished executing.
		codeBlockState := model.codeBlockState[step]
		codeBlockState.StdOut = message.StdOut
		codeBlockState.StdErr = message.StdErr
		codeBlockState.Success = true
		model.viewports.output.SetContent(codeBlockState.StdOut)
		logging.GlobalLogger.Infof("Finished executing: %s", codeBlockState.CodeBlock.Content)

		// Extract the resource group name from the command output if
		// it's not already set.
		if model.resourceGroupName == "" && patterns.AzCommand.MatchString(codeBlockState.CodeBlock.Content) {
			logging.GlobalLogger.Info("Attempting to extract resource group name from command output")
			tmpResourceGroup := az.FindResourceGroupName(codeBlockState.StdOut)
			if tmpResourceGroup != "" {
				logging.GlobalLogger.WithField("resourceGroup", tmpResourceGroup).Info("Found resource group")
				model.resourceGroupName = tmpResourceGroup
			}
		}

		// Increment the codeblock and update the viewport content.
		model.currentCodeBlock++

		commands = append(commands, updateAzureStatus(model))

	case FailedCommandMessage:
		// Handle failed command executions
		model.executingCommand = false
		model.azureStatus.SetError(message.Error)
		environments.AttachResourceURIsToAzureStatus(
			&model.azureStatus,
			model.resourceGroupName,
			model.environment,
		)
		model.viewports.output.SetContent(message.StdErr)
		commands = append(commands, updateAzureStatus(model))

	case AzureStatusUpdatedMessage:
		// After the status has been updated, we force a window resize to
		// render over the status update. For some reason, clearing the screen
		// manually seems to cause the text produced by View() to not render
		// properly if we don't trigger a window size event.
		model.azureStatus.CurrentStep++
		commands = append(commands, func() tea.Msg {
			return tea.WindowSizeMsg{
				Width:  model.width,
				Height: model.height,
			}
		})
	}

	// Update viewport content
	block := model.codeBlockState[model.currentCodeBlock]
	model.viewports.step.SetContent(block.CodeBlock.Description)
	model.viewports.command.SetContent(block.CodeBlock.Content)
	model.viewports.output.SetContent(block.StdOut)

	// Update all the viewports and append resulting commands.
	var command tea.Cmd
	model.viewports.step, command = model.viewports.step.Update(message)
	commands = append(commands, command)
	model.viewports.command, command = model.viewports.command.Update(message)
	commands = append(commands, command)
	model.viewports.output, command = model.viewports.output.Update(message)
	commands = append(commands, command)

	return model, tea.Batch(commands...)
}

// Shows the commands that the user can use to interact with the interactive
// mode model.
func (model InteractiveModeModel) helpView() string {
	return "\n" + model.help.ShortHelpView([]key.Binding{
		model.commands.execute,
		model.commands.skip,
		model.commands.quit,
	})
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

	commandView := fmt.Sprintf("%s\n%s\n%s",
		viewportHeaderView(
			"Command",
			model.width,
		),
		model.viewports.command.View(),
		viewportFooterView(
			" ",
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

	return stepView + "\n" + commandView + "\n" + outputView + "\n" + model.helpView()
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

func NewInteractiveModeModel(
	engine *Engine,
	steps []Step,
	env map[string]string,
) (InteractiveModeModel, error) {
	// TODO: In the future we should just set the current step for the azure status
	// to one as the default.
	var azureStatus = environments.NewAzureDeploymentStatus()
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
			skip: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "Skip the current command."),
			),
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "Quit the scenario."),
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

// Interact with each individual step from a scenario and let the user
// interact with the codecodeBlocks.
func (e *Engine) InteractWithSteps(steps []Step, env map[string]string) error {

	stepsToExecute := filterDeletionCommands(steps, e.Configuration.DoNotDelete)

	model, err := NewInteractiveModeModel(e, stepsToExecute, env)

	if err != nil {
		return err
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		logging.GlobalLogger.Fatalf("Error initializing interactive mode: %s", err)
		return err
	}

	return nil

}
