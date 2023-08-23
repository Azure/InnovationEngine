package shells

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Azure/InnovationEngine/internal/utils"
)

// Location where the environment state from commands is captured and sent to
// for being able to share state across commands.
var environmentStateFile = "/tmp/env.txt"

func loadEnvFile(path string) (map[string]string, error) {
	if !utils.FileExists(path) {
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

// Resets the stored environment variables file.
func ResetStoredEnvironmentVariables() error {
	return os.Remove(environmentStateFile)
}

type CommandOutput struct {
	StdOut string
	StdErr string
}

// Executes a bash command and returns the output or error.
func ExecuteBashCommand(command string, env map[string]string, inheritEnvironment bool, forwardInputOutput bool) (CommandOutput, error) {
	var commandWithStateSaved = []string{
		command,
		"IE_LAST_COMMAND_EXIT_CODE=\"$?\"",
		"env > /tmp/env.txt",
		"exit $IE_LAST_COMMAND_EXIT_CODE",
	}

	commandToExecute := exec.Command("bash", "-c", strings.Join(commandWithStateSaved, "\n"))

	var stdoutBuffer, stderrBuffer bytes.Buffer

	if forwardInputOutput {
		commandToExecute.Stdout = os.Stdout
		commandToExecute.Stderr = os.Stderr
		commandToExecute.Stdin = os.Stdin
	} else {
		// Capture std out and std err as separate buffers.
		commandToExecute.Stdout = &stdoutBuffer
		commandToExecute.Stderr = &stderrBuffer
	}

	if inheritEnvironment {
		commandToExecute.Env = os.Environ()
	}

	// Sharing environment variable state between isolated shell executions is a
	// bit tough, but how we handle it is by storing the environment variables
	// after a command is executed within a file and then loading that file
	// before executing the next command. This allows us to share state between
	// isolated command calls.
	envFromPreviousStep, err := loadEnvFile(environmentStateFile)
	if err == nil {
		merged := utils.MergeMaps(env, envFromPreviousStep)
		for k, v := range merged {
			commandToExecute.Env = append(commandToExecute.Env, fmt.Sprintf("%s=%s", k, v))
		}
	} else {
		for k, v := range env {
			commandToExecute.Env = append(commandToExecute.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	err = commandToExecute.Run()
	if forwardInputOutput {
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
