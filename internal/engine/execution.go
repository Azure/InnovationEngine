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
	"github.com/charmbracelet/lipgloss"
)

const (
	// TODO - Make this configurable for terminals that support it.
	// spinnerFrames  = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
	spinnerFrames  = `-\|/`
	spinnerLength  = 1
	spinnerRefresh = 100 * time.Millisecond
)

var azGroupDelete = regexp.MustCompile(`az group delete`)

// If a scenario has an `az group delete` command and the `--preserve-resources`
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

// Executes the steps from a scenario and renders the output to the terminal.
func (e *Engine) ExecuteAndRenderSteps(steps []Step, env map[string]string) {

	// Enable resource tracking for all of the deployments performed by the
	// innovation engine.
	if e.Configuration.ResourceTracking {
		tracking_id := "6edbe7b9-4e03-4ab0-8213-230ba21aeaba"
		env["AZURE_HTTP_USER_AGENT"] = fmt.Sprintf("pid-%s", tracking_id)
		if e.Configuration.Verbose {
			logging.GlobalLogger.Info("Resource tracking enabled. Tracking ID: " + env["AZURE_HTTP_USER_AGENT"])
			fmt.Println("Resource tracking enabled. Tracking ID: " + env["AZURE_HTTP_USER_AGENT"])
		}
	}

	stepsToExecute := filterDeletionCommands(steps, e.Configuration.DoNotDelete)
	for stepNumber, step := range stepsToExecute {
		fmt.Printf("%d. %s\n", stepNumber+1, step.Name)
		for _, block := range step.CodeBlocks {
			// Render the codeblock.
			indentedBlock := indentMultiLineCommand(block.Content, 4)
			fmt.Print("    " + indentedBlock)

			// Grab the number of lines it contains & set the cursor to the
			// beginning of the block.
			lines := strings.Count(block.Content, "\n")
			fmt.Printf("\033[%dA", lines)

			// Render the spinner and hide the cursor.
			fmt.Print(spinnerStyle.Render("  "+string(spinnerFrames[0])) + " ")
			fmt.Print("\033[?25l")

			// execute the command as a goroutine to allow for the spinner to be
			// rendered while the command is executing.
			done := make(chan error)
			var commandOutput shells.CommandOutput
			go func(block parsers.CodeBlock) {
				output, err := shells.ExecuteBashCommand(block.Content, utils.CopyMap(env), true)
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
					// Show the cursor, check the result of the command, and display the
					// final status.
					fmt.Print("\033[?25h")
					if err == nil {
						fmt.Printf("\r  %s \n", checkStyle.Render("✔"))
						fmt.Printf("\033[%dB\n", lines)
						if e.Configuration.Verbose {
							fmt.Printf("    %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Render(commandOutput.StdOut))
						}
					} else {
						fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
						fmt.Printf("\033[%dB", lines)
						fmt.Printf("    %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render(err.Error()))
					}

					break loop
				default:
					frame = (frame + 1) % len(spinnerFrames)
					fmt.Printf("\r  %s", spinnerStyle.Render(string(spinnerFrames[frame])))
					time.Sleep(spinnerRefresh)
				}
			}
		}
	}
	shells.ResetStoredEnvironmentVariables()
}
