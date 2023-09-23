package shells

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/lib/fs"
)

// Location where the environment state from commands is captured and sent to
// for being able to share state across commands.
var environmentStateFile = "/tmp/env.txt"

func loadEnvFile(path string) (map[string]string, error) {
	if !fs.FileExists(path) {
		return nil, fmt.Errorf("env file '%s' does not exist", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open env file '%s': %w", path, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	env := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2) // Split at the first "=" only
			env[parts[0]] = parts[1]
		}
	}
	return env, nil
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

// Resets the stored environment variables file.
func ResetStoredEnvironmentVariables() error {
	return os.Remove(environmentStateFile)
}

type CommandOutput struct {
	StdOut string
	StdErr string
}

type BashCommandConfiguration struct {
	EnvironmentVariables map[string]string
	InheritEnvironment   bool
	InteractiveCommand   bool
	WriteToHistory       bool
}

// Executes a bash command and returns the output or error.
func ExecuteBashCommand(command string, config BashCommandConfiguration) (CommandOutput, error) {
	var commandWithStateSaved = []string{
		command,
		"IE_LAST_COMMAND_EXIT_CODE=\"$?\"",
		"env > /tmp/env.txt",
		"exit $IE_LAST_COMMAND_EXIT_CODE",
	}

	commandToExecute := exec.Command("bash", "-c", strings.Join(commandWithStateSaved, "\n"))

	var stdoutBuffer, stderrBuffer bytes.Buffer

	if config.InteractiveCommand {
		commandToExecute.Stdout = os.Stdout
		commandToExecute.Stderr = os.Stderr
		commandToExecute.Stdin = os.Stdin
	} else {
		// Capture std out and std err as separate buffers.
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
	envFromPreviousStep, err := loadEnvFile(environmentStateFile)
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
		return CommandOutput{}, fmt.Errorf("command exited with '%w' and the message '%s'", err, standardError)
	}

	return CommandOutput{
		StdOut: standardOutput,
		StdErr: standardError,
	}, nil
}
