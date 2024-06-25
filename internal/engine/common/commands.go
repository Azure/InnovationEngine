package common

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/shells"
	tea "github.com/charmbracelet/bubbletea"
)

// Emitted when a command has been executed successfully.
type SuccessfulCommandMessage struct {
	StdOut string
	StdErr string
}

// Emitted when a command has failed to execute.
type FailedCommandMessage struct {
	StdOut string
	StdErr string
	Error  error
}

type ExitMessage struct {
	EncounteredFailure bool
}

func Exit(encounteredFailure bool) tea.Cmd {
	return func() tea.Msg {
		return ExitMessage{EncounteredFailure: encounteredFailure}
	}
}

// Executes a bash command and returns a tea message with the output. This function
// will be executed asycnhronously.
func ExecuteCodeBlockAsync(codeBlock parsers.CodeBlock, env map[string]string) tea.Cmd {
	return func() tea.Msg {
		logging.GlobalLogger.Infof(
			"Executing command asynchronously:\n %s", codeBlock.Content)

		output, err := shells.ExecuteBashCommand(codeBlock.Content, shells.BashCommandConfiguration{
			EnvironmentVariables: env,
			InheritEnvironment:   true,
			InteractiveCommand:   false,
			WriteToHistory:       true,
		})
		if err != nil {
			logging.GlobalLogger.Errorf("Error executing command:\n %s", err.Error())
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

		outputComparisonError := CompareCommandOutputs(
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
	Program.ReleaseTerminal()

	output, err := shells.ExecuteBashCommand(
		codeBlock.Content,
		shells.BashCommandConfiguration{
			EnvironmentVariables: env,
			InheritEnvironment:   true,
			InteractiveCommand:   true,
			WriteToHistory:       true,
		},
	)

	Program.RestoreTerminal()

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
func ClearScreen() tea.Cmd {
	return func() tea.Msg {
		fmt.Print(
			"\033[H\033[2J",
		) // ANSI escape codes for clearing the screen and repositioning the cursor
		return nil
	}
}

// Updates the azure status with the current state of the interactive mode
// model.
func UpdateAzureStatus(azureStatus environments.AzureDeploymentStatus, environment string) tea.Cmd {
	return func() tea.Msg {
		logging.GlobalLogger.Tracef(
			"Attempting to update the azure status: %+v",
			azureStatus,
		)
		environments.ReportAzureStatus(azureStatus, environment)
		return AzureStatusUpdatedMessage{}
	}
}

// Empty struct used to indicate that the azure status has been updated so
// that we can respond to it within the Update() function.
type AzureStatusUpdatedMessage struct{}
