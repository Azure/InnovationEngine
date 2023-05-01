package render

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

type GroupedCodeBlock struct {
	Header string
	Blocks []parsers.CodeBlock
}

func groupCodeBlocksByHeader(blocks []parsers.CodeBlock) []GroupedCodeBlock {
	var groupedBlocks []GroupedCodeBlock
	var headerIndex = make(map[string]int)

	for _, block := range blocks {
		if index, ok := headerIndex[block.Header]; ok {
			groupedBlocks[index].Blocks = append(groupedBlocks[index].Blocks, block)
		} else {
			headerIndex[block.Header] = len(groupedBlocks)
			groupedBlocks = append(groupedBlocks, GroupedCodeBlock{
				Header: block.Header,
				Blocks: []parsers.CodeBlock{block},
			})
		}
	}

	return groupedBlocks
}

func ExecuteAndRenderCodeBlocks(codeblocks []parsers.CodeBlock, env map[string]string) {
	groupedBlocks := groupCodeBlocksByHeader(codeblocks)

	for _, groupedBlock := range groupedBlocks {
		for _, block := range groupedBlock.Blocks {
			fmt.Print(spinnerStyle.Render(string(spinnerFrames[0])) + " ")

			done := make(chan bool)

			go func(block parsers.CodeBlock) {
				_, err := shells.ExecuteBashCommand(block.Content, env, true)
				if err != nil {
					fmt.Println("Error executing command: ", err)
				}
				done <- true
			}(block)

			frame := 0
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
