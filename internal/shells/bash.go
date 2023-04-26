package shells

import (
	"fmt"
	"os"
	"os/exec"
)

// Executes a bash command and returns the output or error.
func ExecuteBashCommand(command string, env map[string]string, inherit_environment_variables bool) (string, error) {
	command_to_execute := exec.Command("bash", "-c", command)

	if inherit_environment_variables {
		command_to_execute.Env = os.Environ()
	}

	for k, v := range env {
		command_to_execute.Env = append(command_to_execute.Env, fmt.Sprintf("%s=%s", k, v))
	}

	out, err := command_to_execute.Output()
	if err != nil {
		return "", fmt.Errorf("error executing bash command: %w", err)
	}
	return string(out), nil
}
