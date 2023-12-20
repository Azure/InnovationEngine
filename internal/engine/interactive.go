package engine

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/ui"
)

// Interact with each individual step from a scenario and let the user
// interact with the codeblocks.
func (e *Engine) InteractWithSteps(steps []Step, env map[string]string) error {
	stepsToExecute := filterDeletionCommands(steps, e.Configuration.DoNotDelete)

	for stepNumber, step := range stepsToExecute {
		stepTitle := fmt.Sprintf("  %d. %s\n", stepNumber+1, step.Name)
		fmt.Println(ui.StepTitleStyle.Render(stepTitle))
		for codeBlockNumber, codeBlock := range step.CodeBlocks {
			fmt.Println(
				ui.InteractiveModeCodeblockDescriptionStyle.Render(
					fmt.Sprintf(
						"    %d.%d %s",
						stepNumber+1,
						codeBlockNumber+1,
						codeBlock.Description,
					),
				),
			)
			fmt.Print(
				ui.IndentMultiLineCommand(
					fmt.Sprintf(
						"      %s",
						ui.InteractiveModeCodeblockStyle.Render(
							codeBlock.Content,
						),
					),
					6),
			)
			fmt.Println()
		}
	}

	return nil
}
