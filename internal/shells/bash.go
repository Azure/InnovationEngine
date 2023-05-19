package shells

import (
	"bufio"
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
			parts := strings.Split(line, "=")
			env[parts[0]] = parts[1]
		}
	}
	return env, nil
}

// Resets the stored environment variables file.
func ResetStoredEnvironmentVariables() error {
	return os.Remove(environmentStateFile)
}

func mergeMaps(a, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}

	return a
}

// Executes a bash command and returns the output or error.
func ExecuteBashCommand(command string, env map[string]string, inherit_environment_variables bool) (string, error) {
	var commandWithState = []string{
		command,
		"env > /tmp/env.txt",
	}
	commandToExecute := exec.Command("bash", "-c", strings.Join(commandWithState, "\n"))

	if inherit_environment_variables {
		commandToExecute.Env = os.Environ()
	}

	envFromPreviousStep, err := loadEnvFile("/tmp/env.txt")
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

	stdOutAndErr, err := commandToExecute.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command exited with '%w' and the message '%s'", err, stdOutAndErr)
	}

	return string(stdOutAndErr), nil
}
