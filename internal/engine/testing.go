package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/parsers"
	"github.com/Azure/InnovationEngine/internal/shells"
	"github.com/Azure/InnovationEngine/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/xrash/smetrics"
)

func (e *Engine) TestSteps(steps []Step, env map[string]string) {
	for stepNumber, step := range steps {
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
						if block.ExpectedOutput.Language == "json" {
							actualOutput, err := utils.OrderJsonFields(commandOutput)
							if err != nil {
								fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
								fmt.Printf("\033[%dB", lines)
								fmt.Printf("    %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render(err.Error()))
								break loop
							}

							expectedOutput, err := utils.OrderJsonFields(block.ExpectedOutput.Content)
							if err != nil {
								fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
								fmt.Printf("\033[%dB", lines)
								fmt.Printf("    %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render(err.Error()))
								break loop
							}

							score := smetrics.JaroWinkler(expectedOutput, actualOutput, 0.7, 4)

							if block.ExpectedOutput.ExpectedSimilarity > score {
								fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
								fmt.Printf("\033[%dB", lines)
								fmt.Printf("    %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render("Expected output does not match actual output."))
								fmt.Printf("	%s\n", utils.GetDifferenceBetweenStrings(expectedOutput, actualOutput))
							}
						} else {
							score := smetrics.JaroWinkler(block.ExpectedOutput.Content, commandOutput, 0.7, 4)
							if block.ExpectedOutput.ExpectedSimilarity > score {
								fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
								fmt.Printf("\033[%dB", lines)
								fmt.Printf("    %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render("Expected output does not match actual output."))
								fmt.Printf("	%s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render(utils.GetDifferenceBetweenStrings(block.ExpectedOutput.Content, commandOutput)))
							}
						}

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
	shells.ResetStoredEnvironmentVariables()
}
