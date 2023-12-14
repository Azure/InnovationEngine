package engine

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/terminal"
	"github.com/Azure/InnovationEngine/internal/ui"
)

func (e *Engine) InteractWithSteps(steps []Step, env map[string]string) error {
	stepsToExecute := filterDeletionCommands(steps, true)
	for stepNumber, step := range stepsToExecute {
		stepTitle := fmt.Sprintf("  %d. %s\n", stepNumber+1, step.Name)
		fmt.Println(ui.StepTitleStyle.Render(stepTitle))
		terminal.MoveCursorPositionUp(1)
		terminal.HideCursor()
	}

  terminal.ShowCursor()
	return nil
}
