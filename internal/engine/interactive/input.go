package interactive

import (
	"strconv"

	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/patterns"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

// NewInteractiveModeCommands creates a new set of interactive mode commands.
func NewInteractiveModeCommands() InteractiveModeCommands {
	return InteractiveModeCommands{
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
		executeAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Execute all remaining commands."),
		),
		executeMany: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m<number><enter>", "Execute the next <number> commands."),
		),
		pause: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Pause execution of commands."),
		),
	}
}

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
