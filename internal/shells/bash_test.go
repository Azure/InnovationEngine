package shells

import (
	"testing"
)

func TestEnvironmentVariableValidationAndFiltering(t *testing.T) {
	// Test key validation
	t.Run("Key Validation", func(t *testing.T) {
		validCases := []struct {
			key      string
			expected bool
		}{
			{"ValidKey", true},
			{"VALID_VARIABLE", true},
			{"_AnotherValidKey", true},
			{"123Key", false},                   // Starts with a digit
			{"key-with-hyphen", false},          // Contains a hyphen
			{"key.with.dot", false},             // Contains a period
			{"Fabric_NET-0-[Delegated]", false}, // From cloud shell environment.
		}

		for _, tc := range validCases {
			t.Run(tc.key, func(t *testing.T) {
				result := environmentVariableName.MatchString(tc.key)
				if result != tc.expected {
					t.Errorf(
						"Expected isValidKey(%s) to be %v, got %v",
						tc.key,
						tc.expected,
						result,
					)
				}
			})
		}
	})

	// Test key filtering
	t.Run("Key Filtering", func(t *testing.T) {
		envMap := map[string]string{
			"ValidKey":                 "value1",
			"_AnotherValidKey":         "value2",
			"123Key":                   "value3",
			"key-with-hyphen":          "value4",
			"key.with.dot":             "value5",
			"Fabric_NET-0-[Delegated]": "false", // From cloud shell environment.
		}

		validEnvMap := filterInvalidKeys(envMap)

		expectedValidEnvMap := map[string]string{
			"ValidKey":         "value1",
			"_AnotherValidKey": "value2",
		}

		if len(validEnvMap) != len(expectedValidEnvMap) {
			t.Errorf(
				"Expected validEnvMap to have %d keys, got %d",
				len(expectedValidEnvMap),
				len(validEnvMap),
			)
		}

		for key, value := range validEnvMap {
			if expectedValue, ok := expectedValidEnvMap[key]; !ok || value != expectedValue {
				t.Errorf("Expected validEnvMap[%s] to be %s, got %s", key, expectedValue, value)
			}
		}
	})
}

func TestBashCommandExecution(t *testing.T) {
	// Test command execution
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

	t.Run("Command with multiple commands", func(t *testing.T) {
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
