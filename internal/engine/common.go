package engine

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

const (
	setCursor = "\033[?25h"
)

func moveCursorPositionDown(lines int) {
	fmt.Printf("\033[" + string(lines) + "B\n")
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
