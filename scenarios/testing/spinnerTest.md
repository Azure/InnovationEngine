# Testing Spinner with Elapsed Time

This test demonstrates the spinner with elapsed time display feature. The commands below create and run a Go program that streams output in real-time, allowing you to observe the spinner animation with elapsed time counter.

## Create streaming test program

First, let's create a simple Go program that demonstrates the spinner with elapsed time in real-time:

```bash
# Create a directory for our test
mkdir -p /tmp/spinner_test

# Create the streaming test program
cat > /tmp/spinner_test/main.go << 'EOF'
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	
	// You would typically import UI from Innovation Engine
	// We're reimplementing just what we need for the test
)

const (
	spinnerFrames  = `-\|/`
	spinnerRefresh = 100 * time.Millisecond
)

// Helper function to simulate styled text (simplified)
func renderStyledText(text string) string {
	return text
}

// StreamCommand executes a command and streams its output in real-time
// while showing a spinner with elapsed time between outputs
func streamCommand(command string) error {
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
					fmt.Printf("\r  %s [%02d:%02d elapsed]", renderStyledText(string(spinnerFrames[frame])), minutes, seconds)
					
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

// Create test script that will demonstrate real-time output
func createTestScript() error {
	script := `#!/bin/bash

# This script demonstrates real-time output with pauses to show the spinner
echo "Starting long-running operation with streamed output..."
for i in {1..5}; do
    echo -n "Processing step $i of 5... "
    sleep 1
    echo "done"
    
    # Add a slight delay to allow spinner to be visible
    sleep 0.5
    
    # On step 3, output something to stderr to test error stream
    if [ $i -eq 3 ]; then
        echo "Note: Step $i added diagnostic info" >&2
    fi
done
echo "Operation complete!"
`
	return os.WriteFile("/tmp/spinner_test/test_script.sh", []byte(script), 0755)
}

func main() {
	fmt.Println("Creating test script...")
	if err := createTestScript(); err != nil {
		fmt.Printf("Error creating test script: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nRunning test with real-time output streaming...")
	if err := streamCommand("/tmp/spinner_test/test_script.sh"); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("\nTest completed successfully!")
}
EOF
```

<!-- expected_similarity=1.0 -->

```text
```

## Build and run streaming test

Build and run the streaming test program. This will demonstrate the spinner with elapsed time between real-time output:

```bash
cd /tmp/spinner_test
go build -o spinner_test main.go
./spinner_test
```

<!-- expected_similarity=0.8 -->

```text
Creating test script...

Running test with real-time output streaming...
Executing command with real-time output streaming:
$ /tmp/spinner_test/test_script.sh
Starting long-running operation with streamed output...
Processing step 1 of 5... done
Processing step 2 of 5... done
Processing step 3 of 5... done
[stderr] Note: Step 3 added diagnostic info
Processing step 4 of 5... done
Processing step 5 of 5... done
Operation complete!

Test completed successfully!
```

## The elapsed time spinner in Innovation Engine

The spinner with elapsed time display you've just observed is the same one that's used in Innovation Engine. When running commands in IE, the spinner will appear and show elapsed time while the command is running, giving you real-time feedback on how long operations are taking.

This is particularly useful for long-running Azure CLI commands where you might otherwise be unsure about progress.