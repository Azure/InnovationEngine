package lib

import (
	"testing"
)

// Simple test to ensure that Max() returns the correct value.
func TestMax(t *testing.T) {
	got := Max(1, 2)
	if got != 2 {
		t.Errorf("Max(1, 2) = %d; want 2", got)
	}
}
