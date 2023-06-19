package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/logging"
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
						actualOutput := commandOutput.StdOut
						expectedOutput := block.ExpectedOutput.Content
						expectedSimilarity := block.ExpectedOutput.ExpectedSimilarity

						if block.ExpectedOutput.Language == "json" {
							logging.GlobalLogger.Debugf("Comparing JSON strings:\nExpected: %s\nActual%s", expectedOutput, actualOutput)
							meetsThreshold, err := utils.CompareJsonStrings(actualOutput, expectedOutput, expectedSimilarity)
							if err != nil {
								fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
								fmt.Printf("\033[%dB", lines)
								fmt.Printf("  %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render(err.Error()))
								break loop
							}

							if !meetsThreshold {
								fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
								fmt.Printf("\033[%dB", lines)
								fmt.Printf("  %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render("Expected output does not match actual output."))
								fmt.Printf("	%s\n", utils.GetDifferenceBetweenStrings(expectedOutput, actualOutput))
								break loop
							}

							if e.Configuration.Verbose {
								score, _ := utils.ComputeJsonStringSimilarity(actualOutput, expectedOutput)

								actual, _ := utils.OrderJsonFields(actualOutput)
								expected, _ := utils.OrderJsonFields(expectedOutput)

								logging.GlobalLogger.WithField("actual", actual).WithField("expected", expected).Debugf("JaroWinkler score: %f Expected Similarity: %f", score, expectedSimilarity)
							}
						} else {
							score := smetrics.JaroWinkler(block.ExpectedOutput.Content, commandOutput.StdOut, 0.7, 4)
							if block.ExpectedOutput.ExpectedSimilarity > score {
								fmt.Printf("\r  %s \n", errorStyle.Render("✗"))
								fmt.Printf("\033[%dB", lines)
								fmt.Printf("    %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733")).Render("Expected output does not match actual output."))
								fmt.Printf("	%s\n", utils.GetDifferenceBetweenStrings(block.ExpectedOutput.Content, commandOutput.StdOut))
							}
						}

						fmt.Printf("\r  %s \n", checkStyle.Render("✔"))
						fmt.Printf("\033[%dB\n", lines)
						if e.Configuration.Verbose {
							fmt.Printf("  %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Render(commandOutput.StdOut))
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
