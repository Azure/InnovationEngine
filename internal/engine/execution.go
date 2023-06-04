package engine

import (
	"fmt"
	"regexp"
	"strings"
	"time"

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

// Styles used for rendering output to the terminal.
var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF"))
	checkStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Align(lipgloss.Center).Bold(true)
)

var azGroupDelete = regexp.MustCompile(`az group delete`)

// Indents a multi-line command to be nested under the first line of the
// command.
func indentMultiLineCommand(content string, indentation int) string {
	lines := strings.Split(content, "\n")
	for i := 1; i < len(lines); i++ {
		if strings.HasSuffix(strings.TrimSpace(lines[i-1]), "\\") {
			lines[i] = strings.Repeat(" ", indentation) + lines[i]
		} else if strings.TrimSpace(lines[i]) != "" {
			lines[i] = strings.Repeat(" ", indentation) + lines[i]
		}

	}
	return strings.Join(lines, "\n")
}

// Executes the steps from a scenario and renders the output to the terminal.
func (e *Engine) ExecuteAndRenderSteps(steps []Step, env map[string]string) {

	// Enable resource tracking for all of the deployments performed by the
	// innovation engine.
	if e.Configuration.ResourceTracking {
		tracking_id := "6edbe7b9-4e03-4ab0-8213-230ba21aeaba"
		env["AZURE_HTTP_USER_AGENT"] = fmt.Sprintf("pid-%s", tracking_id)
		if e.Configuration.Verbose {
			fmt.Println("Resource tracking enabled. Tracking ID: " + env["AZURE_HTTP_USER_AGENT"])
		}
	}

	// If a scenario has an `az group delete` command and the `--do-not-delete`
	// flag is set, we remove it from the steps.

	stepsToExecute := []Step{}
	if e.Configuration.DoNotDelete {
		for _, step := range steps {
			newBlocks := []parsers.CodeBlock{}
			for _, block := range step.CodeBlocks {
				if azGroupDelete.MatchString(block.Content) {
					if e.Configuration.Verbose {
						fmt.Printf("Found az group delete command within the step: %s\n", step.Name)
					}
				} else {
					newBlocks = append(newBlocks, block)
				}
			}
			if len(newBlocks) > 0 {
				stepsToExecute = append(stepsToExecute, Step{
					Name:       step.Name,
					CodeBlocks: newBlocks,
				})
			}
		}
	} else {
		stepsToExecute = steps
	}

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
			var commandOutput string
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
							fmt.Printf("    %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Render(commandOutput))
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
}
