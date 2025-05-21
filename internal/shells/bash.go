package shells

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/Azure/InnovationEngine/internal/lib"
)

// streamWriter implements io.Writer to capture and forward command output in real-time
type streamWriter struct {
	callback OutputCallback
	isStderr bool
}

func (w *streamWriter) Write(p []byte) (n int, err error) {
	if w.callback != nil {
		w.callback(string(p), w.isStderr)
	}
	return len(p), nil
}

func appendToBashHistory(command string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Lock the file to prevent other processes from writing to it concurrently
	// and then  unlock after we're done writing to it.
	if err := unix.Flock(int(file.Fd()), unix.LOCK_EX); err != nil {
		return fmt.Errorf("failed to lock file: %w", err)
	}
	defer unix.Flock(int(file.Fd()), unix.LOCK_UN)

	// Append the command and a newline to the file
	_, err = file.WriteString(command + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

type CommandOutput struct {
	StdOut string
	StdErr string
}

type OutputCallback func(string, bool)

type BashCommandConfiguration struct {
	EnvironmentVariables map[string]string
	InheritEnvironment   bool
	InteractiveCommand   bool
	WriteToHistory       bool
	StreamOutput         bool
	OutputCallback       OutputCallback
}

var ExecuteBashCommand = executeBashCommandImpl

// Executes a bash command and returns the output or error.
func executeBashCommandImpl(
	command string,
	config BashCommandConfiguration,
) (CommandOutput, error) {
	commandWithStateSaved := []string{
		"set -e",
		command,
		"IE_LAST_COMMAND_EXIT_CODE=\"$?\"",
		"env > " + lib.DefaultEnvironmentStateFile,
		"exit $IE_LAST_COMMAND_EXIT_CODE",
	}

	commandToExecute := exec.Command("bash", "-c", strings.Join(commandWithStateSaved, "\n"))

	var stdoutBuffer, stderrBuffer bytes.Buffer

	// If the command requires interaction, we provide the user with the ability
	// to interact with the command. However, we cannot capture the buffer this
	// way.
	if config.InteractiveCommand {
		commandToExecute.Stdout = os.Stdout
		commandToExecute.Stderr = os.Stderr
		commandToExecute.Stdin = os.Stdin
	} else if config.StreamOutput && config.OutputCallback != nil {
		// Create multi-writers to capture output both in buffer and stream it via callback
		stdoutWriter := io.MultiWriter(&stdoutBuffer, &streamWriter{
			callback: config.OutputCallback,
			isStderr: false,
		})
		stderrWriter := io.MultiWriter(&stderrBuffer, &streamWriter{
			callback: config.OutputCallback,
			isStderr: true,
		})
		commandToExecute.Stdout = stdoutWriter
		commandToExecute.Stderr = stderrWriter
	} else {
		commandToExecute.Stdout = &stdoutBuffer
		commandToExecute.Stderr = &stderrBuffer
	}

	if config.InheritEnvironment {
		commandToExecute.Env = os.Environ()
	}

	// Sharing environment variable state between isolated shell executions is a
	// bit tough, but how we handle it is by storing the environment variables
	// after a command is executed within a file and then loading that file
	// before executing the next command. This allows us to share state between
	// isolated command calls.
	envFromPreviousStep, err := lib.LoadEnvironmentStateFile(lib.DefaultEnvironmentStateFile)
	if err == nil {
		merged := lib.MergeMaps(config.EnvironmentVariables, envFromPreviousStep)
		for k, v := range merged {
			commandToExecute.Env = append(commandToExecute.Env, fmt.Sprintf("%s=%s", k, v))
		}
	} else {
		for k, v := range config.EnvironmentVariables {
			commandToExecute.Env = append(commandToExecute.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	if config.WriteToHistory {

		homeDir, err := lib.GetHomeDirectory()
		if err != nil {
			return CommandOutput{}, fmt.Errorf("failed to get home directory: %w", err)
		}

		err = appendToBashHistory(command, homeDir+"/.bash_history")
		if err != nil {
			return CommandOutput{}, fmt.Errorf("failed to write command to history: %w", err)
		}
	}

	err = commandToExecute.Run()

	// TODO(vmarcella): Find a better way to handle this.
	if config.InteractiveCommand {
		return CommandOutput{}, err
	}

	standardOutput, standardError := stdoutBuffer.String(), stderrBuffer.String()

	if err != nil {
		return CommandOutput{
				StdOut: standardOutput,
				StdErr: standardError,
			}, fmt.Errorf(
				"command exited with '%w' and the message '%s'",
				err,
				standardError,
			)
	}

	return CommandOutput{
		StdOut: standardOutput,
		StdErr: standardError,
	}, nil
}
