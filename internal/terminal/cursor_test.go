package terminal

import "testing"

func TestCursorManipulation(t *testing.T) {

	t.Run("Test moving cursor up", func(t *testing.T) {
		position := MoveCursorPositionUp(1)
		if position != "\033[1A" {
			t.Errorf("Expected cursor to move up 1 line, got %s", position)
		}

		position = MoveCursorPositionUp(2)
		if position != "\033[2A" {
			t.Errorf("Expected cursor to move up 2 lines, got %s", position)
		}
	})

	t.Run("Test moving cursor down", func(t *testing.T) {
		position := MoveCursorPositionDown(1)
		if position != "\033[1B\n" {
			t.Errorf("Expected cursor to move up 1 line, got %s", position)
		}

		position = MoveCursorPositionDown(2)
		if position != "\033[2B\n" {
			t.Errorf("Expected cursor to move up 2 lines, got %s", position)
		}
	})
}
