package common

import (
	"fmt"
	"strings"

	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/shells"
	tea "github.com/charmbracelet/bubbletea"
)

// Emitted when a command has been executed successfully.
type SuccessfulCommandMessage struct {
	StdOut          string
	StdErr          string
	SimilarityScore float64
}

// Emitted when a command has failed to execute.
type FailedCommandMessage struct {
	StdOut          string
	StdErr          string
	Error           error
	SimilarityScore float64
}

// Emitted when command output is streaming
type StreamingOutputMessage struct {
	Output   string
	IsStderr bool
}

type ExitMessage struct {
	EncounteredFailure bool
}

func SendStreamingOutput(output string, isStderr bool) tea.Cmd {
	return func() tea.Msg {
		return StreamingOutputMessage{
			Output:   output,
			IsStderr: isStderr,
		}
	}
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

		var accumulatedOutput strings.Builder
		
		output, err := shells.ExecuteBashCommand(codeBlock.Content, shells.BashCommandConfiguration{
			EnvironmentVariables: env,
			InheritEnvironment:   true,
			InteractiveCommand:   false,
			WriteToHistory:       true,
			StreamOutput:         true,
			OutputCallback: func(output string, isStderr bool) {
				// Accumulate the output
				accumulatedOutput.WriteString(output)
				
				// Log the streaming output
				if isStderr {
					logging.GlobalLogger.Infof("Streaming stderr: %s", output)
				} else {
					logging.GlobalLogger.Infof("Streaming stdout: %s", output)
				}
				
				// Print the output directly to show streaming works
				// This is a simple approach for testing the streaming functionality
				fmt.Print(output)
			},
		})
		
		// Update output with accumulated content if needed
		if output.StdOut == "" && accumulatedOutput.Len() > 0 {
			output.StdOut = accumulatedOutput.String()
		}
		
		if err != nil {
			logging.GlobalLogger.Errorf("Error executing command:\n %s", err.Error())
			return FailedCommandMessage{
				StdOut:          output.StdOut,
				StdErr:          output.StdErr,
				Error:           err,
				SimilarityScore: 0,
			}
		}

		// Check command output against the expected output.
		actualOutput := output.StdOut
		expectedOutput := codeBlock.ExpectedOutput.Content
		expectedSimilarity := codeBlock.ExpectedOutput.ExpectedSimilarity
		expectedRegex := codeBlock.ExpectedOutput.ExpectedRegex
		expectedOutputLanguage := codeBlock.ExpectedOutput.Language

		score, outputComparisonError := CompareCommandOutputs(
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
				StdOut:          output.StdOut,
				StdErr:          output.StdErr,
				Error:           outputComparisonError,
				SimilarityScore: score,
			}

		}

		logging.GlobalLogger.Infof("Command output to stdout:\n %s", output.StdOut)
		return SuccessfulCommandMessage{
			StdOut:          output.StdOut,
			StdErr:          output.StdErr,
			SimilarityScore: score,
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
