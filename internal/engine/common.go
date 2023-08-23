package engine

import (
	"fmt"
	"strings"

	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/xrash/smetrics"
)

// Styles used for rendering output to the terminal.
var (
	scenarioTitleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Align(lipgloss.Center).Bold(true).Underline(true)
	stepTitleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#518BAD")).Align(lipgloss.Left).Bold(true)
	spinnerStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#518BAD"))
	verboseStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#437684")).Align(lipgloss.Left)
	checkStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32"))
	errorStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	errorMessageStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733"))
	ocdStatusUpdateStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
)

// Hides the cursor from the CLI using ANSI escape codes.
func hideCursor() {
	fmt.Print("\033[?25l")
}

// Displays the cursor in the CLI using ANSI escape codes.
func showCursor() {
	fmt.Print("\033[?25h")
}

// Moves the cursor up a specified number of lines.
func moveCursorPositionUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

// Moves the cursor down a specified number of lines.
func moveCursorPositionDown(lines int) {
	fmt.Printf("\033[%dB\n", lines)
}

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

// Compares the actual output of a command to the expected output of a command.
func compareCommandOutputs(actualOutput string, expectedOutput string, expectedSimilarity float64, expectedOutputLanguage string) error {
	if expectedOutputLanguage == "json" {
		logging.GlobalLogger.Debugf("Comparing JSON strings:\nExpected: %s\nActual%s", expectedOutput, actualOutput)
		meetsThreshold, err := utils.CompareJsonStrings(actualOutput, expectedOutput, expectedSimilarity)

		if err != nil {
			return err
		}

		if !meetsThreshold {
			return fmt.Errorf(errorMessageStyle.Render("Expected output does not match actual output."))
		}

		score, _ := utils.ComputeJsonStringSimilarity(actualOutput, expectedOutput)

		actual, _ := utils.OrderJsonFields(actualOutput)
		expected, _ := utils.OrderJsonFields(expectedOutput)

		logging.GlobalLogger.WithField("actual", actual).WithField("expected", expected).Debugf("Jaro score: %f Expected Similarity: %f", score, expectedSimilarity)
	} else {
		score := smetrics.JaroWinkler(expectedOutput, actualOutput, 0.7, 4)

		if expectedSimilarity > score {
			return fmt.Errorf(errorMessageStyle.Render("Expected output does not match actual output."))
		}
	}

	return nil
}
