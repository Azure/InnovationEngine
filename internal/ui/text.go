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
	InteractiveModeStepTitleStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#518BAD")).
					Bold(true)
	InteractiveModeCodeblockDescriptionStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("#ffffff"))
	InteractiveModeCodeblockStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#fff")).Background(lipgloss.Color("#000000"))
)

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
