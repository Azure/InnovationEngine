package engine

import (
	"errors"
	"fmt"
	"time"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/patterns"
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/Azure/InnovationEngine/internal/terminal"
	"github.com/Azure/InnovationEngine/internal/ui"
)

func (e *Engine) TestSteps(steps []Step, env map[string]string) error {
	var resourceGroupName string
	stepsToExecute := filterDeletionCommands(steps, true)

	var testRunnerError error = nil
testRunner:
	for stepNumber, step := range stepsToExecute {
		stepTitle := fmt.Sprintf("  %d. %s\n", stepNumber+1, step.Name)
		fmt.Println(ui.StepTitleStyle.Render(stepTitle))
		terminal.MoveCursorPositionUp(1)
		terminal.HideCursor()

		for _, block := range step.CodeBlocks {
			// execute the command as a goroutine to allow for the spinner to be
			// rendered while the command is executing.
			done := make(chan error)
			var commandOutput shells.CommandOutput
			go func(block parsers.CodeBlock) {
				logging.GlobalLogger.Infof("Executing command: %s", block.Content)
				output, err := shells.ExecuteBashCommand(block.Content, shells.BashCommandConfiguration{EnvironmentVariables: lib.CopyMap(env), InheritEnvironment: true, InteractiveCommand: false, WriteToHistory: true})
				logging.GlobalLogger.Infof("Command stdout: %s", output.StdOut)
				logging.GlobalLogger.Infof("Command stderr: %s", output.StdErr)
				commandOutput = output
				done <- err
			}(block)

			frame := 0
			var err error

		loop:
			// While the command is executing, render the spinner.
			for {
				select {
				case err = <-done:
					terminal.ShowCursor()

					if err == nil {
						actualOutput := commandOutput.StdOut
						expectedOutput := block.ExpectedOutput.Content
						expectedSimilarity := block.ExpectedOutput.ExpectedSimilarity
						expectedRegex := block.ExpectedOutput.ExpectedRegex
						expectedOutputLanguage := block.ExpectedOutput.Language

						err := compareCommandOutputs(actualOutput, expectedOutput, expectedSimilarity, expectedRegex, expectedOutputLanguage)

						if err != nil {
							logging.GlobalLogger.Errorf("Error comparing command outputs: %s", err.Error())
							fmt.Print(ui.ErrorStyle.Render("Error when comparing the command outputs: %s\n", err.Error()))
						}

						// Extract the resource group name from the command output if
						// it's not already set.
						if resourceGroupName == "" && patterns.AzCommand.MatchString(block.Content) {
							tmpResourceGroup := az.FindResourceGroupName(commandOutput.StdOut)
							if tmpResourceGroup != "" {
								logging.GlobalLogger.Infof("Found resource group: %s", tmpResourceGroup)
								resourceGroupName = tmpResourceGroup
							}
						}

						fmt.Printf("\r  %s \n", ui.CheckStyle.Render("✔"))
						terminal.MoveCursorPositionDown(1)
					} else {

						fmt.Printf("\r  %s \n", ui.ErrorStyle.Render("✗"))
						terminal.MoveCursorPositionDown(1)
						fmt.Printf(" %s\n", ui.ErrorStyle.Render("Error executing command: %s\n", err.Error()))

						logging.GlobalLogger.Errorf("Error executing command: %s", err.Error())

						testRunnerError = err
						break testRunner
					}

					break loop
				default:
					frame = (frame + 1) % len(spinnerFrames)
					fmt.Printf("\r  %s", ui.SpinnerStyle.Render(string(spinnerFrames[frame])))
					time.Sleep(spinnerRefresh)
				}
			}
		}
	}

	// If the resource group name was set, delete it.
	if resourceGroupName != "" {
		fmt.Printf("\n")
		fmt.Printf("Deleting resource group: %s\n", resourceGroupName)
		command := fmt.Sprintf("az group delete --name %s --yes", resourceGroupName)
		output, err := shells.ExecuteBashCommand(
			command,
			shells.BashCommandConfiguration{
				EnvironmentVariables: lib.CopyMap(env),
				InheritEnvironment:   true,
				InteractiveCommand:   false,
				WriteToHistory:       true,
			},
		)

		if err != nil {
			fmt.Print(ui.ErrorStyle.Render("Error deleting resource group: %s\n", err.Error()))
			logging.GlobalLogger.Errorf("Error deleting resource group: %s", err.Error())
			testRunnerError = errors.Join(testRunnerError, err)
		}

		fmt.Print(output.StdOut)
	}

	shells.ResetStoredEnvironmentVariables()
	return testRunnerError
}
