package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles used for rendering output to the terminal.
var (
	ScenarioTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6CB6FF")).
				Align(lipgloss.Center).
				Bold(true).
				Underline(true)
	StepTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#518BAD")).
			Align(lipgloss.Left).
			Bold(true)
	SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#518BAD"))
	VerboseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#437684")).
			Align(lipgloss.Left)
	CheckStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#32CD32"))
	ErrorStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	ErrorMessageStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733"))
	OcdStatusUpdateStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
)

var (
	InteractiveModeCodeBlockDescriptionStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("#ffffff"))
	InteractiveModeCodeBlockStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#fff"))

	InteractiveModeStepTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}().Foreground(lipgloss.Color("#518BAD")).Bold(true)

	InteractiveModeStepFooterStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return InteractiveModeStepTitleStyle.Copy().BorderStyle(b)
	}().Foreground(lipgloss.Color("#fff"))
)

// Indents a multi-line command to be nested under the first line of the
// command.
func IndentMultiLineCommand(content string, indentation int) string {
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

func RemoveHorizontalAlign(s string) string {
	return strings.Join(
		mapSliceString(
			strings.Split(s, "\n"),
			func(s string) string { return strings.TrimRight(s, " ") },
		),
		"\n",
	)
}

func mapSliceString(slice []string, apply func(string) string) []string {
	var result []string
	for _, s := range slice {
		result = append(result, apply(s))
	}
	return result
}
