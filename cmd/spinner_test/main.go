package main

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/Azure/InnovationEngine/internal/testutil"
)

func main() {
	// Find the root directory and set it as environment variable
	execPath, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(execPath)
		// Go up two directories (from cmd/spinner_test to root)
		rootDir := filepath.Dir(filepath.Dir(dir))
		os.Setenv("INNOVATION_ENGINE_ROOT", rootDir)
	}

	fmt.Println("Running spinner test with real-time streaming...")
	output := testutil.RunStreamingTest()
	fmt.Println("\nTest result:", output)
}