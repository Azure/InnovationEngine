package shells

import (
	"fmt"
	"os"
	"os/exec"
)

// Executes a bash command and returns the output or error.
func ExecuteBashCommand(command string, env map[string]string, inherit_environment_variables bool) (string, error) {
	commandToExecute := exec.Command("bash", "-c", command)

	if inherit_environment_variables {
		commandToExecute.Env = os.Environ()
	}

	for k, v := range env {
		commandToExecute.Env = append(commandToExecute.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdOutAndErr, err := commandToExecute.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command exited with '%w' and the message '%s'", err, stdOutAndErr)
	}
	return string(stdOutAndErr), nil
}
