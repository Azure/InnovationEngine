package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/engine/common"
	"github.com/Azure/InnovationEngine/internal/engine/environments"
	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/patterns"
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/Azure/InnovationEngine/internal/terminal"
	"github.com/Azure/InnovationEngine/internal/ui"
)

const (
	// TODO - Make this configurable for terminals that support it.
	// spinnerFrames  = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
	spinnerFrames  = `-\|/`
	spinnerRefresh = 100 * time.Millisecond
)

// If a scenario has an `az group delete` command and the `--do-not-delete`
// flag is set, we remove it from the steps.
func filterDeletionCommands(steps []common.Step, preserveResources bool) []common.Step {
	filteredSteps := []common.Step{}
	if preserveResources {
		for _, step := range steps {
			newBlocks := []parsers.CodeBlock{}
			for _, block := range step.CodeBlocks {
				if patterns.AzGroupDelete.MatchString(block.Content) {
					continue
				} else {
					newBlocks = append(newBlocks, block)
				}
			}
			if len(newBlocks) > -1 {
				filteredSteps = append(filteredSteps, common.Step{
					Name:       step.Name,
					CodeBlocks: newBlocks,
				})
			}
		}
	} else {
		filteredSteps = steps
	}
	return filteredSteps
}

func renderCommand(blockContent string) (shells.CommandOutput, error) {
	escapedCommand := blockContent
	if !patterns.MultilineQuotedStringCommand.MatchString(blockContent) {
		escapedCommand = strings.ReplaceAll(blockContent, "\\\n", "\\\\\n")
	}
	renderedCommand, err := shells.ExecuteBashCommand(
		"echo -e \""+escapedCommand+"\"",
		shells.BashCommandConfiguration{
			EnvironmentVariables: map[string]string{},
			InteractiveCommand:   false,
			WriteToHistory:       false,
			InheritEnvironment:   true,
		},
	)
	return renderedCommand, err
}

// Executes the steps from a scenario and renders the output to the terminal.
func (e *Engine) ExecuteAndRenderSteps(steps []common.Step, env map[string]string) error {
	var resourceGroupName string = ""
	azureStatus := environments.NewAzureDeploymentStatus()

	err := az.SetSubscription(e.Configuration.Subscription)
	if err != nil {
		logging.GlobalLogger.Errorf("Invalid Config: Failed to set subscription: %s", err)
		azureStatus.SetError(err)
		environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)
		return err
	}

	stepsToExecute := filterDeletionCommands(steps, e.Configuration.DoNotDelete)

	for stepNumber, step := range stepsToExecute {

		azureCodeBlocks := []environments.AzureCodeBlock{}
		for _, block := range step.CodeBlocks {
			azureCodeBlocks = append(azureCodeBlocks, environments.AzureCodeBlock{
				Command:     block.Content,
				Description: block.Description,
			})
		}

		azureStatus.AddStep(fmt.Sprintf("%d. %s", stepNumber+1, step.Name), azureCodeBlocks)
	}

	environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)

	for stepNumber, step := range stepsToExecute {
		stepTitle := fmt.Sprintf("%d. %s\n", stepNumber+1, step.Name)
		fmt.Println(ui.StepTitleStyle.Render(stepTitle))
		azureStatus.CurrentStep = stepNumber + 1

		for _, block := range step.CodeBlocks {
			var finalCommandOutput string
			if e.Configuration.RenderValues {
				// Render the codeblock.
				renderedCommand, err := renderCommand(block.Content)
				if err != nil {
					logging.GlobalLogger.Errorf("Failed to render command: %s", err.Error())
					azureStatus.SetError(err)
					environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)
					return err
				}
				finalCommandOutput = ui.IndentMultiLineCommand(renderedCommand.StdOut, 4)
			} else {
				finalCommandOutput = ui.IndentMultiLineCommand(block.Content, 4)
			}

			fmt.Print("    " + finalCommandOutput)

			// execute the command as a goroutine to allow for the spinner to be
			// rendered while the command is executing.
			done := make(chan error)
			var commandOutput shells.CommandOutput

			// If the command is an SSH command, we need to forward the input and
			// output
			interactiveCommand := false
			if patterns.SshCommand.MatchString(block.Content) {
				interactiveCommand = true
			}

			logging.GlobalLogger.WithField("isInteractive", interactiveCommand).
				Infof("Executing command: %s", block.Content)

			var commandErr error
			var frame int = 0

			// If forwarding input/output, don't render the spinner.
			if !interactiveCommand {
				// Grab the number of lines it contains & set the cursor to the
				// beginning of the block.

				lines := strings.Count(finalCommandOutput, "\n")
				terminal.MoveCursorPositionUp(lines)

				// Render the spinner and hide the cursor.
				fmt.Print(ui.SpinnerStyle.Render("  "+string(spinnerFrames[0])) + " ")
				terminal.HideCursor()

				go func(block parsers.CodeBlock) {
					output, err := shells.ExecuteBashCommand(
						block.Content,
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
				}(block)
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
							expectedOutput := block.ExpectedOutput.Content
							expectedSimilarity := block.ExpectedOutput.ExpectedSimilarity
							expectedRegex := block.ExpectedOutput.ExpectedRegex
							expectedOutputLanguage := block.ExpectedOutput.Language

							_, outputComparisonError := common.CompareCommandOutputs(actualOutput, expectedOutput, expectedSimilarity, expectedRegex, expectedOutputLanguage)

							if outputComparisonError != nil {
								logging.GlobalLogger.Errorf("Error comparing command outputs: %s", outputComparisonError.Error())
								fmt.Printf("\r  %s \n", ui.ErrorStyle.Render("✗"))
								terminal.MoveCursorPositionDown(lines)
								fmt.Printf("  %s\n", ui.ErrorMessageStyle.Render(outputComparisonError.Error()))
								fmt.Printf("	%s\n", lib.GetDifferenceBetweenStrings(block.ExpectedOutput.Content, commandOutput.StdOut))

								azureStatus.SetError(outputComparisonError)
								environments.AttachResourceURIsToAzureStatus(
									&azureStatus,
									resourceGroupName,
									e.Configuration.Environment,
								)
								environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)

								return outputComparisonError
							}

							fmt.Printf("\r  %s \n", ui.CheckStyle.Render("✔"))
							terminal.MoveCursorPositionDown(lines)

							fmt.Printf("%s\n", ui.RemoveHorizontalAlign(ui.VerboseStyle.Render(commandOutput.StdOut)))

							// Extract the resource group name from the command output if
							// it's not already set.
							if resourceGroupName == "" && patterns.AzCommand.MatchString(block.Content) {
								logging.GlobalLogger.Info("Attempting to extract resource group name from command output")
								tmpResourceGroup := az.FindResourceGroupName(commandOutput.StdOut)
								if tmpResourceGroup != "" {
									logging.GlobalLogger.WithField("resourceGroup", tmpResourceGroup).Info("Found resource group")
									resourceGroupName = tmpResourceGroup
									azureStatus.AddResourceURI(az.BuildResourceGroupId(e.Configuration.Subscription, resourceGroupName))
								}
							}

							if stepNumber != len(stepsToExecute)-1 {
								environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)
							}

						} else {
							terminal.ShowCursor()
							fmt.Printf("\r  %s \n", ui.ErrorStyle.Render("✗"))
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
			} else {
				lines := strings.Count(block.Content, "\n")

				// If we're on the last step and the command is an SSH command, we need
				// to report the status before executing the command. This is needed for
				// one click deployments and does not affect the normal execution flow.
				if stepNumber == len(stepsToExecute)-1 && patterns.SshCommand.MatchString(block.Content) {
					azureStatus.Status = "Succeeded"
					environments.AttachResourceURIsToAzureStatus(&azureStatus, resourceGroupName, e.Configuration.Environment)
					environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)
				}

				output, commandExecutionError := shells.ExecuteBashCommand(
					block.Content,
					shells.BashCommandConfiguration{
						EnvironmentVariables: lib.CopyMap(env),
						InheritEnvironment:   true,
						InteractiveCommand:   true,
						WriteToHistory:       false,
					},
				)

				terminal.ShowCursor()

				if commandExecutionError == nil {
					fmt.Printf("\r  %s \n", ui.CheckStyle.Render("✔"))
					terminal.MoveCursorPositionDown(lines)

					fmt.Printf("  %s\n", ui.VerboseStyle.Render(output.StdOut))

					if stepNumber != len(stepsToExecute)-1 {
						environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)
					}
				} else {
					fmt.Printf("\r  %s \n", ui.ErrorStyle.Render("✗"))
					terminal.MoveCursorPositionDown(lines)
					fmt.Printf("  %s\n", ui.ErrorMessageStyle.Render(commandExecutionError.Error()))

					azureStatus.SetError(commandExecutionError)
					environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)
					return commandExecutionError
				}
			}
		}
	}

	// Report the final status of the deployment (Only applies to one click deployments).
	azureStatus.Status = "Succeeded"
	environments.AttachResourceURIsToAzureStatus(
		&azureStatus,
		resourceGroupName,
		e.Configuration.Environment,
	)
	environments.ReportAzureStatus(azureStatus, e.Configuration.Environment)

	switch e.Configuration.Environment {
	case environments.EnvironmentsAzure, environments.EnvironmentsOCD:
		logging.GlobalLogger.Info(
			"Cleaning environment variable file located at /tmp/env-vars",
		)
		err := lib.CleanEnvironmentStateFile(lib.DefaultEnvironmentStateFile)
		if err != nil {
			logging.GlobalLogger.Errorf("Error cleaning environment variables: %s", err.Error())
			return err
		}

		logging.GlobalLogger.Info(
			"Cleaning working directory file located at /tmp/working-dir",
		)
		err = lib.DeleteWorkingDirectoryStateFile(lib.DefaultWorkingDirectoryStateFile)
		if err != nil {
			logging.GlobalLogger.Errorf("Error cleaning working directory: %s", err.Error())
			return err
		}

	default:
		lib.DeleteEnvironmentStateFile(lib.DefaultEnvironmentStateFile)
		lib.DeleteWorkingDirectoryStateFile(lib.DefaultWorkingDirectoryStateFile)
	}

	return nil
}
