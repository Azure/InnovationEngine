package engine

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/ocd"
	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/Azure/InnovationEngine/internal/utils"
)

const (
	// TODO - Make this configurable for terminals that support it.
	// spinnerFrames  = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
	spinnerFrames  = `-\|/`
	spinnerRefresh = 100 * time.Millisecond
)

var (
	sshCommand = regexp.MustCompile(`(^|\s)\bssh\b\s`)

	// Az cli command regex
	azCommand     = regexp.MustCompile(`az\s+([a-z]+)\s+([a-z]+)`)
	azGroupDelete = regexp.MustCompile(`az group delete`)

	// ARM regex
	azResourceURI       = regexp.MustCompile(`\"id\": \"(/subscriptions/[^\"]+)\"`)
	azResourceGroupName = regexp.MustCompile(`resourceGroups/([^\"]+)`)
)

// If a scenario has an `az group delete` command and the `--do-not-delete`
// flag is set, we remove it from the steps.
func filterDeletionCommands(steps []Step, preserveResources bool) []Step {
	filteredSteps := []Step{}
	if preserveResources {
		for _, step := range steps {
			newBlocks := []parsers.CodeBlock{}
			for _, block := range step.CodeBlocks {
				if azGroupDelete.MatchString(block.Content) {
					continue
				} else {
					newBlocks = append(newBlocks, block)
				}
			}
			if len(newBlocks) > -1 {
				filteredSteps = append(filteredSteps, Step{
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

// Print out the one click deployment status if in the correct environment.
func reportOCDStatus(status ocd.OneClickDeploymentStatus, environment string) {
	if environment == EnvironmentsOCD {
		statusJson, err := status.AsJsonString()
		if err != nil {
			logging.GlobalLogger.Error("Failed to marshal status", err)
		} else {
			// We add these strings to the output so that the portal can find and parse
			// the JSON status.
			ocdStatus := fmt.Sprintf("ie_us%sie_ue\n", statusJson)
			fmt.Println(ocdStatusUpdateStyle.Render(ocdStatus))
		}
	}
}

// Find the resource group name from the output of an az command.
func findResourceGroupName(commandOutput string) string {
	matches := azResourceGroupName.FindStringSubmatch(commandOutput)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// Find all the deployed resources in a resource group.
func findAllDeployedResourceURIs(resourceGroup string) []string {
	output, err := shells.ExecuteBashCommand("az resource list -g"+resourceGroup, map[string]string{}, true, false)

	if err != nil {
		logging.GlobalLogger.Error("Failed to list deployments", err)
	}

	matches := azResourceURI.FindAllStringSubmatch(output.StdOut, -1)
	results := []string{}
	for _, match := range matches {
		results = append(results, match[1])
	}
	return results
}

// Executes the steps from a scenario and renders the output to the terminal.
func (e *Engine) ExecuteAndRenderSteps(steps []Step, env map[string]string) {

	// If the correlation ID is set, we need to set the AZURE_HTTP_USER_AGENT
	// environment variable so that the Azure CLI will send the correlation ID
	// with Azure Resource Manager requests.
	if e.Configuration.CorrelationId != "" {
		env["AZURE_HTTP_USER_AGENT"] = fmt.Sprintf("innovation-engine-%s", e.Configuration.CorrelationId)
		logging.GlobalLogger.Info("Resource tracking enabled. Tracking ID: " + env["AZURE_HTTP_USER_AGENT"])
	}

	var resourceGroupName string
	var ocdStatus = ocd.NewOneClickDeploymentStatus()

	stepsToExecute := filterDeletionCommands(steps, e.Configuration.DoNotDelete)

	for stepNumber, step := range stepsToExecute {
		ocdStatus.AddStep(fmt.Sprintf("%d. %s", stepNumber+1, step.Name))
	}

	reportOCDStatus(ocdStatus, e.Configuration.Environment)

	for stepNumber, step := range stepsToExecute {
		stepTitle := fmt.Sprintf("%d. %s\n", stepNumber+1, step.Name)
		fmt.Println(stepTitleStyle.Render(stepTitle))
		ocdStatus.CurrentStep = stepNumber + 1

		for _, block := range step.CodeBlocks {
			// Render the codeblock.
			indentedBlock := indentMultiLineCommand(block.Content, 4)
			fmt.Print("    " + indentedBlock)

			// execute the command as a goroutine to allow for the spinner to be
			// rendered while the command is executing.
			done := make(chan error)
			var commandOutput shells.CommandOutput

			// If the command is an SSH command, we need to forward the input and
			// output
			interactiveCommand := false
			if sshCommand.MatchString(block.Content) {
				interactiveCommand = true
			}

			logging.GlobalLogger.WithField("isInteractive", interactiveCommand).Infof("Executing command: %s", block.Content)

			var commandErr error
			var frame int = 0

			// If forwarding input/output, don't render the spinner.
			if !interactiveCommand {
				// Grab the number of lines it contains & set the cursor to the
				// beginning of the block.
				lines := strings.Count(block.Content, "\n")
				moveCursorPositionUp(lines)

				// Render the spinner and hide the cursor.
				fmt.Print(spinnerStyle.Render("  "+string(spinnerFrames[0])) + " ")
				hideCursor()

				go func(block parsers.CodeBlock) {
					output, err := shells.ExecuteBashCommand(block.Content, utils.CopyMap(env), true, interactiveCommand)
					commandOutput = output
					done <- err
				}(block)
			loop:
				// While the command is executing, render the spinner.
				for {
					select {
					case commandErr = <-done:
						// Show the cursor, check the result of the command, and display the
						// final status.
						showCursor()

						if commandErr == nil {
							fmt.Printf("\r  %s \n", checkStyle.Render("✔"))

							moveCursorPositionDown(lines)
							if e.Configuration.Verbose {
								fmt.Printf("  %s\n", verboseStyle.Render(commandOutput.StdOut))
							}

							// Extract the resource group name from the command output if
							// it's not already set.
							if resourceGroupName == "" && azCommand.MatchString(block.Content) {
								tmpResourceGroup := findResourceGroupName(commandOutput.StdOut)
								if tmpResourceGroup != "" {
									logging.GlobalLogger.WithField("resourceGroup", tmpResourceGroup).Info("Found resource group")
									resourceGroupName = tmpResourceGroup
								}
							}

							if stepNumber != len(stepsToExecute)-1 {
								reportOCDStatus(ocdStatus, e.Configuration.Environment)
							}

						} else {
							showCursor()
							fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
							moveCursorPositionDown(lines)
							fmt.Printf("  %s\n", errorMessageStyle.Render(commandErr.Error()))

							ocdStatus.SetError(commandErr)
							reportOCDStatus(ocdStatus, e.Configuration.Environment)

							os.Exit(1)
						}

						break loop
					default:
						frame = (frame + 1) % len(spinnerFrames)
						fmt.Printf("\r  %s", spinnerStyle.Render(string(spinnerFrames[frame])))
						time.Sleep(spinnerRefresh)
					}
				}
			} else {
				lines := strings.Count(block.Content, "\n")
				output, err := shells.ExecuteBashCommand(block.Content, utils.CopyMap(env), true, interactiveCommand)

				if err == nil {
					showCursor()
					fmt.Printf("\r  %s \n", checkStyle.Render("✔"))
					moveCursorPositionDown(lines)

					if e.Configuration.Verbose {
						fmt.Printf("  %s\n", verboseStyle.Render(output.StdOut))
					}
					if stepNumber != len(stepsToExecute)-1 {
						reportOCDStatus(ocdStatus, e.Configuration.Environment)
					}
				} else {
					showCursor()
					fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
					moveCursorPositionDown(lines)
					fmt.Printf("  %s\n", errorMessageStyle.Render(err.Error()))

					ocdStatus.SetError(err)
					reportOCDStatus(ocdStatus, e.Configuration.Environment)
					os.Exit(1)
				}
			}
		}
	}

	ocdStatus.Status = "Succeeded"

	if resourceGroupName != "" {
		resourceURIs := findAllDeployedResourceURIs(resourceGroupName)

		if len(resourceURIs) > 0 {
			logging.GlobalLogger.WithField("resourceURIs", resourceURIs).Info("Found deployed resources.")
			ocdStatus.ResourceURIs = resourceURIs
		}
	} else {
		logging.GlobalLogger.Warn("No resource group name found. Unable to find deployed resources.")
	}

	reportOCDStatus(ocdStatus, e.Configuration.Environment)

	shells.ResetStoredEnvironmentVariables()
}
