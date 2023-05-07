package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/charmbracelet/lipgloss"
)

const (
	// TODO - Make this configurable for terminals that support it.
	// spinnerFrames  = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
	spinnerFrames  = `-\|/`
	spinnerLength  = 1
	spinnerRefresh = 100 * time.Millisecond
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF"))
	checkStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Align(lipgloss.Center).Bold(true)
)

func indentSubsequentLines(content string, index string) string {
	lines := strings.Split(content, "\n")
	for i := 1; i < len(lines); i++ {
		if strings.HasSuffix(strings.TrimSpace(lines[i-1]), "\\") {
			lines[i] = index + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}

func ExecuteAndRenderSteps(steps []Step, env map[string]string) {
	for stepNumber, step := range steps {
		fmt.Printf("%d. %s\n", stepNumber+1, step.Name)
		for _, block := range step.CodeBlocks {
			// Render the codeblock.
			indentedBlock := indentSubsequentLines(block.Content, "    ")
			fmt.Print("    " + indentedBlock)

			// Grab the number of lines it contains & set the cursor to the
			// beginning of the block.
			lines := strings.Count(block.Content, "\n")
			fmt.Printf("\033[%dA", lines)

			// Render the spinner and hide the curosr
			fmt.Print(spinnerStyle.Render("  "+string(spinnerFrames[0])) + " ")

			fmt.Print("\033[?25l")
			done := make(chan error)

			// Create a goroutine to execute the command and pass in the codeblock to
			go func(block parsers.CodeBlock) {
				_, err := shells.ExecuteBashCommand(block.Content, env, true)
				done <- err
			}(block)

			frame := 0
			var err error

		loop:
			// While the command is executing, render the spinner.
			for {
				select {
				case err = <-done:
					fmt.Print("\033[?25h")
					// Show the cursor & clear the spinner.
					if err == nil {
						fmt.Printf("\r  %s \n", checkStyle.Render("✔"))
						fmt.Printf("\033[%dB", lines)
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
