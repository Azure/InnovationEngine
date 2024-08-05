package lib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Azure/InnovationEngine/internal/lib/fs"
)

// Get environment variables from the current process.
func GetEnvironmentVariables() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}

	return envMap
}

// Location where the environment state from commands are to be captured
// and sent to for being able to share state across commands.
var DefaultEnvironmentStateFile = "/tmp/env-vars"

// Loads a file that contains environment variables
func LoadEnvironmentStateFile(path string) (map[string]string, error) {
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
			value := parts[1]
			if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
				// Remove leading and trailing quotes
				value = value[1 : len(value)-1]
			}
			env[parts[0]] = value
		}
	}
	return env, nil
}

func CleanEnvironmentStateFile(path string) error {
	env, err := LoadEnvironmentStateFile(path)
	if err != nil {
		return err
	}

	env = filterInvalidKeys(env)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for k, v := range env {
		_, err := fmt.Fprintf(writer, "%s=\"%s\"\n", k, v)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

var environmentVariableName = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

func filterInvalidKeys(envMap map[string]string) map[string]string {
	validEnvMap := make(map[string]string)
	for key, value := range envMap {
		if environmentVariableName.MatchString(key) {
			validEnvMap[key] = value
		}
	}
	return validEnvMap
}

// Deletes the stored environment variables file.
func DeleteEnvironmentStateFile(path string) error {
	return os.Remove(path)
}
