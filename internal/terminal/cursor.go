package terminal

import "fmt"

// Hides the cursor from the CLI using ANSI escape codes.
func HideCursor() string {
	hide := "\033[?25l"
	fmt.Print("\033[?25l")
	return hide
}

// Displays the cursor in the CLI using ANSI escape codes.
func ShowCursor() string {
	show := "\033[?25h"
	fmt.Print("\033[?25h")
	return show
}

// Moves the cursor up a specified number of lines.
func MoveCursorPositionUp(lines int) string {
	position := fmt.Sprintf("\033[%dA", lines)
	fmt.Print(position)
	return position

}

// Moves the cursor down a specified number of lines.
func MoveCursorPositionDown(lines int) string {
	position := fmt.Sprintf("\033[%dB\n", lines)
	fmt.Print(position)
	return position
}
