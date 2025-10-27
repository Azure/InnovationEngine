package testutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Azure/InnovationEngine/internal/ui"
)

const (
	spinnerFrames  = `-\|/`
	spinnerRefresh = 100 * time.Millisecond
)

// StreamCommand executes a command and streams its output in real-time
// while showing a spinner with elapsed time between outputs
func StreamCommand(command string) error {
	fmt.Println("Executing command with real-time output streaming:")
	fmt.Println("$ " + command)

	// Create the command
	cmd := exec.Command("bash", "-c", command)

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	// Create channels to signal when stdout/stderr reading is done
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	// Track when the last output was printed
	lastOutputTime := time.Now()
	commandStartTime := time.Now()

	// Function to read from a pipe and print to console
	readFromPipe := func(pipe io.ReadCloser, isDone chan<- struct{}, prefix string) {
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Printf("\r%s%s\n", prefix, text)
			lastOutputTime = time.Now()
		}
		close(isDone)
	}

	// Read stdout and stderr concurrently
	go readFromPipe(stdout, stdoutDone, "")
	go readFromPipe(stderr, stderrDone, "[stderr] ")

	// Spinner goroutine
	spinnerDone := make(chan struct{})
	go func() {
		frame := 0
		for {
			select {
			case <-spinnerDone:
				return
			default:
				// Only show spinner if it's been a while since last output
				if time.Since(lastOutputTime) > 200*time.Millisecond {
					elapsedTime := time.Since(commandStartTime)
					minutes := int(elapsedTime.Minutes())
					seconds := int(elapsedTime.Seconds()) % 60
					
					// Clear the current line and show spinner
					fmt.Printf("\r  %s [%02d:%02d elapsed]", ui.SpinnerStyle.Render(string(spinnerFrames[frame])), minutes, seconds)
					
					frame = (frame + 1) % len(spinnerFrames)
				}
				time.Sleep(spinnerRefresh)
			}
		}
	}()

	// Wait for stdout and stderr to finish
	<-stdoutDone
	<-stderrDone

	// Stop the spinner
	close(spinnerDone)
	fmt.Print("\r                                 \r") // Clear the spinner line

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("command exited with code %d", exitErr.ExitCode())
		}
		return fmt.Errorf("error waiting for command: %v", err)
	}

	return nil
}

// RunStreamingTest runs a command that demonstrates the spinner with elapsed time
// Returns a string with the captured output for test verification
func RunStreamingTest() string {
	// Create a temporary file to capture output
	outputCaptured := strings.Builder{}
	
	// Get the path to the stream_output.sh script
	scriptPath := os.Getenv("INNOVATION_ENGINE_ROOT")
	if scriptPath == "" {
		// If environment variable not set, use current directory
		var err error
		scriptPath, err = os.Getwd()
		if err != nil {
			outputCaptured.WriteString(fmt.Sprintf("Error getting working directory: %v\n", err))
			return outputCaptured.String()
		}
	}
	
	fullScriptPath := fmt.Sprintf("%s/internal/testutil/stream_output.sh", scriptPath)
	outputCaptured.WriteString(fmt.Sprintf("Looking for script at: %s\n", fullScriptPath))
	
	// Verify script exists
	if _, err := os.Stat(fullScriptPath); os.IsNotExist(err) {
		outputCaptured.WriteString(fmt.Sprintf("Error: Script not found at %s\n", fullScriptPath))
		return outputCaptured.String()
	}
	
	outputCaptured.WriteString("Streaming test started\n")
	
	// Run the command
	err := StreamCommand(fullScriptPath)
	if err != nil {
		outputCaptured.WriteString(fmt.Sprintf("Error: %v\n", err))
	} else {
		outputCaptured.WriteString("Streaming test completed successfully\n")
	}
	
	return outputCaptured.String()
}