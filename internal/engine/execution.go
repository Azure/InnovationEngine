package engine

import (
	"fmt"
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
)

func ExecuteAndRenderSteps(steps []Step, env map[string]string) {
	for _, step := range steps {
		for _, block := range step.CodeBlocks {
			fmt.Print(spinnerStyle.Render(string(spinnerFrames[0])) + " ")

			done := make(chan bool)

			// Create a goroutine to execute the command and pass in the codeblock to
			go func(block parsers.CodeBlock) {
				_, err := shells.ExecuteBashCommand(block.Content, env, true)
				if err != nil {
					fmt.Println("Error executing command: ", err)
				}
				done <- true
			}(block)

			frame := 0

			// While the command is executing, render the spinner.
		loop:
			for {
				select {
				case <-done:
					fmt.Printf("\r%s %s\n", checkStyle.Render("✔"), block.Header)
					break loop
				default:
					frame = (frame + 1) % len(spinnerFrames)
					fmt.Printf("\r%s %s", spinnerStyle.Render(string(spinnerFrames[frame])), block.Header)
					time.Sleep(spinnerRefresh)
				}
			}
		}
	}
}
