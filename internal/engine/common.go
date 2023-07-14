package engine

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles used for rendering output to the terminal.
var (
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF"))
	checkStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32"))
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	errorMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733"))
	verboseStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Align(lipgloss.Left)
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CB6FF")).Align(lipgloss.Center).Bold(true)
)

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
