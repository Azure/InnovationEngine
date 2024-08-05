package shells

import (
	"testing"
)

func TestBashCommandExecution(t *testing.T) {
	// Ensures that if a command succeeds, the output is returned.
	t.Run("Valid command execution", func(t *testing.T) {
		cmd := "printf hello"
		result, err := ExecuteBashCommand(
			cmd,
			BashCommandConfiguration{
				EnvironmentVariables: nil,
				InheritEnvironment:   true,
				InteractiveCommand:   false,
				WriteToHistory:       false,
			},
		)
		if err != nil {
			t.Errorf("Expected err to be nil, got %v", err)
		}
		if result.StdOut != "hello" {
			t.Errorf("Expected result to be non-empty, got '%s'", result.StdOut)
		}
	})

	// Ensures that if a command fails, an error is returned.
	t.Run("Invalid command execution", func(t *testing.T) {
		cmd := "not_real_command"
		_, err := ExecuteBashCommand(
			cmd,
			BashCommandConfiguration{
				EnvironmentVariables: nil,
				InheritEnvironment:   true,
				InteractiveCommand:   false,
				WriteToHistory:       false,
			},
		)

		if err == nil {
			t.Errorf("Expected an error to occur, but the command succeeded.")
		}
	})

	// Test the execution of commands with multiple subcommands.
	t.Run("Command with multiple subcommands", func(t *testing.T) {
		cmd := "printf hello; printf world"
		result, err := ExecuteBashCommand(
			cmd,
			BashCommandConfiguration{
				EnvironmentVariables: nil,
				InheritEnvironment:   true,
				InteractiveCommand:   false,
				WriteToHistory:       false,
			},
		)
		if err != nil {
			t.Errorf("Expected err to be nil, got %v", err)
		}
		if result.StdOut != "helloworld" {
			t.Errorf("Expected result to be non-empty, got '%s'", result.StdOut)
		}
	})

	// Ensures that if one of the subcommands fail, the other commands do
	// as well.
	t.Run("Command with multiple subcommands exits on first error", func(t *testing.T) {
		cmd := "printf hello; not_real_command; printf world"
		_, err := ExecuteBashCommand(
			cmd,
			BashCommandConfiguration{
				EnvironmentVariables: nil,
				InheritEnvironment:   true,
				InteractiveCommand:   false,
				WriteToHistory:       false,
			},
		)

		if err == nil {
			t.Errorf("Expected an error to occur, but the command succeeded.")
		}
	})

	// Ensures that commands can access environment variables passed into
	// the configuration.
	t.Run("Command with environment variables", func(t *testing.T) {
		cmd := "printf $TEST_ENV_VAR"
		result, err := ExecuteBashCommand(
			cmd,
			BashCommandConfiguration{
				EnvironmentVariables: map[string]string{
					"TEST_ENV_VAR": "hello",
				},
				InheritEnvironment: true,
				InteractiveCommand: false,
				WriteToHistory:     false,
			},
		)
		if err != nil {
			t.Errorf("Expected err to be nil, got %v", err)
		}

		if result.StdOut != "hello" {
			t.Errorf("Expected result to be non-empty, got '%s'", result.StdOut)
		}
	})
}
