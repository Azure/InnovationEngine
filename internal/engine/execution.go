package engine

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/logging"
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

var azGroupDelete = regexp.MustCompile(`az group delete`)
var azCommand = regexp.MustCompile(`az\s+([a-z]+)\s+([a-z]+)`)
var sshCommand = regexp.MustCompile(`(^|\s)\bssh\b\s`)

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

// Check for errors from the Azure CLI. The Azure CLI doesn't return a non-zero
// exit code when an error occurs, so we have to check the output for errors.
func checkForAzCLIError(command string, output shells.CommandOutput) bool {
	if !azCommand.MatchString(command) {
		return false
	}

	if output.StdOut == "" && output.StdErr != "" {
		return true
	}

	return false
}

// Executes the steps from a scenario and renders the output to the terminal.
func (e *Engine) ExecuteAndRenderSteps(steps []Step, env map[string]string) {

	// If the correlation ID is set, we need to set the AZURE_HTTP_USER_AGENT
	// environment variable so that the Azure CLI will send the correlation ID
	// with Azure Resource Manager requests.
	if e.Configuration.CorrelationId != "" {
		env["AZURE_HTTP_USER_AGENT"] = fmt.Sprintf("innovation-engine-%s", e.Configuration.CorrelationId)
		if e.Configuration.Verbose {
			logging.GlobalLogger.Info("Resource tracking enabled. Tracking ID: " + env["AZURE_HTTP_USER_AGENT"])
		}
	}

	stepsToExecute := filterDeletionCommands(steps, e.Configuration.DoNotDelete)
	for stepNumber, step := range stepsToExecute {
		fmt.Printf("%d. %s\n", stepNumber+1, step.Name)
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
			forward_input_output := false
			if sshCommand.MatchString(block.Content) {
				forward_input_output = true
			}

			logging.GlobalLogger.WithField("forward_intput_output", forward_input_output).Info("Executing command: " + block.Content)

			var commandErr error
			var frame int = 0

			// If forwarding input/output, don't render the spinner.
			if !forward_input_output {
				// Grab the number of lines it contains & set the cursor to the
				// beginning of the block.
				lines := strings.Count(block.Content, "\n")
				fmt.Printf("\033[%dA", lines)

				// Render the spinner and hide the cursor.
				fmt.Print(spinnerStyle.Render("  "+string(spinnerFrames[0])) + " ")
				fmt.Print("\033[?25l")

				go func(block parsers.CodeBlock) {
					output, err := shells.ExecuteBashCommand(block.Content, utils.CopyMap(env), true, forward_input_output)
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
						fmt.Print("\033[?25h")

						if checkForAzCLIError(block.Content, commandOutput) {
							commandErr = fmt.Errorf(commandOutput.StdErr)
						}

						if commandErr == nil {
							fmt.Printf("\r  %s \n", checkStyle.Render("✔"))

							fmt.Printf("\033[%dB\n", lines)
							if e.Configuration.Verbose {
								fmt.Printf("  %s\n", verboseStyle.Render(commandOutput.StdOut))
							}
						} else {
							fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
							fmt.Printf("\033[%dB", lines)
							fmt.Printf("  %s\n", errorMessageStyle.Render(commandErr.Error()))
						}

						break loop
					default:
						frame = (frame + 1) % len(spinnerFrames)
						fmt.Printf("\r  %s", spinnerStyle.Render(string(spinnerFrames[frame])))
						time.Sleep(spinnerRefresh)
					}
				}
			} else {
				func(block parsers.CodeBlock) {
					shells.ExecuteBashCommand(block.Content, utils.CopyMap(env), true, forward_input_output)
				}(block)
			}
		}
	}
	shells.ResetStoredEnvironmentVariables()
}
