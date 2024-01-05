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
	"github.com/Azure/InnovationEngine/internal/terminal"
	"github.com/Azure/InnovationEngine/internal/ui"
	"github.com/eiannone/keyboard"
)

// Interact with each individual step from a scenario and let the user
// interact with the codecodeBlocks.
func (e *Engine) InteractWithSteps(steps []Step, env map[string]string) error {
	var azureStatus = environments.NewAzureDeploymentStatus()
	var resourceGroupName string = ""
	err := az.SetSubscription(e.Configuration.Subscription)
	if err != nil {
		logging.GlobalLogger.Errorf("Invalid Config: Failed to set subscription: %s", err)
		azureStatus.SetError(err)
		environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)
		return err
	}

	stepsToExecute := filterDeletionCommands(steps, e.Configuration.DoNotDelete)

	// Add steps to the azure status
	for stepNumber, step := range stepsToExecute {
		azureStatus.AddStep(fmt.Sprintf("%d. %s", stepNumber+1, step.Name))
	}
	environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)

	err = keyboard.Open()
	if err != nil {
		logging.GlobalLogger.Fatalf("Error opening keyboard: %s", err)
	}

	defer keyboard.Close()

	for stepNumber, step := range stepsToExecute {
		stepTitle := fmt.Sprintf("Step %d. %s\n", stepNumber+1, step.Name)
		fmt.Println(ui.StepTitleStyle.Render(stepTitle))
		azureStatus.CurrentStep = stepNumber + 1

		for codeBlockNumber, codeBlock := range step.CodeBlocks {
			fmt.Println(
				ui.InteractiveModeCodeBlockDescriptionStyle.Render(
					codeBlock.Description,
				),
			)
			fmt.Println(
				ui.InteractiveModeCodeBlockStyle.Render(
					codeBlock.Content,
				),
			)

			validCommandEntered := false

			for !validCommandEntered {
				fmt.Print(
					"Enter a command to proceed or h to see available commands: ",
				)
				char, _, err := keyboard.GetKey()
				if err != nil {
					logging.GlobalLogger.Fatalf("Error reading keyboard input: %s", err)
				}

				switch char {

				case 'e':
					// execute the command as a goroutine to allow for the spinner to be
					// rendered while the command is executing.
					validCommandEntered = true
					done := make(chan error)
					var commandOutput shells.CommandOutput
					lines := strings.Count(codeBlock.Content, "\n") + 2
					// If the command is an SSH command, we need to forward the input and
					// output
					interactiveCommand := false
					if patterns.SshCommand.MatchString(codeBlock.Content) {
						interactiveCommand = true
					}

					logging.GlobalLogger.WithField("isInteractive", interactiveCommand).
						Infof("Executing command: %s", codeBlock.Content)

					var commandErr error
					var frame int = 0

					// execute the codecodeBlock
					go func(codeBlock parsers.CodeBlock) {
						output, err := shells.ExecuteBashCommand(
							codeBlock.Content,
							shells.BashCommandConfiguration{
								EnvironmentVariables: lib.CopyMap(env),
								InheritEnvironment:   true,
								InteractiveCommand:   false,
								WriteToHistory:       true,
							},
						)
						logging.GlobalLogger.Infof("Command output to stdout:\n %s", output.StdOut)
						logging.GlobalLogger.Infof("Command output to stderr:\n %s", output.StdErr)
						commandOutput = output
						done <- err
					}(codeBlock)
				renderingLoop:
					// While the command is executing, render the spinner.
					for {
						select {
						case commandErr = <-done:
							// Show the cursor, check the result of the command, and display the
							// final status.
							terminal.ShowCursor()

							if commandErr == nil {

								actualOutput := commandOutput.StdOut
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
									logging.GlobalLogger.Errorf("Error comparing command outputs: %s", outputComparisonError.Error())

									fmt.Printf("\r  %s \n", ui.ErrorStyle.Render("Γ£ù"))
									terminal.MoveCursorPositionDown(lines)
									fmt.Printf("  %s\n", ui.ErrorMessageStyle.Render(outputComparisonError.Error()))
									fmt.Printf("	%s\n", lib.GetDifferenceBetweenStrings(codeBlock.ExpectedOutput.Content, commandOutput.StdOut))

									azureStatus.SetError(outputComparisonError)
									environments.AttachResourceURIsToAzureStatus(
										&azureStatus,
										resourceGroupName,
										e.Configuration.Environment,
									)
									environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)

									return outputComparisonError
								}

								fmt.Printf("\r  %s \n", ui.CheckStyle.Render("Γ£ö"))
								terminal.MoveCursorPositionDown(lines)

								fmt.Printf("%s\n", ui.RemoveHorizontalAlign(ui.VerboseStyle.Render(commandOutput.StdOut)))

								// Extract the resource group name from the command output if
								// it's not already set.
								if resourceGroupName == "" && patterns.AzCommand.MatchString(codeBlock.Content) {
									logging.GlobalLogger.Info("Attempting to extract resource group name from command output")
									tmpResourceGroup := az.FindResourceGroupName(commandOutput.StdOut)
									if tmpResourceGroup != "" {
										logging.GlobalLogger.WithField("resourceGroup", tmpResourceGroup).Info("Found resource group")
										resourceGroupName = tmpResourceGroup
									}
								}

								if stepNumber != len(stepsToExecute)-1 {
									environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)
								}

							} else {
								terminal.ShowCursor()
								fmt.Printf("\r  %s \n", ui.ErrorStyle.Render("Γ£ù"))
								terminal.MoveCursorPositionDown(lines)
								fmt.Printf("  %s\n", ui.ErrorMessageStyle.Render(commandErr.Error()))

								logging.GlobalLogger.Errorf("Error executing command: %s", commandErr.Error())

								azureStatus.SetError(commandErr)
								environments.AttachResourceURIsToAzureStatus(
									&azureStatus,
									resourceGroupName,
									e.Configuration.Environment,
								)
								environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)

								return commandErr
							}

							break renderingLoop

						default:
							frame = (frame + 1) % len(spinnerFrames)
							fmt.Printf("\r  %s", ui.SpinnerStyle.Render(string(spinnerFrames[frame])))
							time.Sleep(spinnerRefresh)
						}
					}

				case 's':
					// skip the codeblock.
					validCommandEntered = true
					logging.GlobalLogger.Infof(
						"Skip used on step %d.%d",
						stepNumber,
						codeBlockNumber,
					)

				case 'q':
					// quit the program
					logging.GlobalLogger.Info("Quit command entered, exiting interactive mode")
					return nil

				case 'h':
					fallthrough
				default:
					// If h any other key is entered, show the available commands.
					fmt.Println("Available commands:")
					fmt.Println("  e - execute this step")
					fmt.Println("  h - show this help")
					fmt.Println("  s - skip this step")
					fmt.Println("  q - quit")

				}

			}

		}
	}

	return nil
}
