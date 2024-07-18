package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock HTTP server for testing downloading markdown from URL
func mockHTTPServer(content string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, content)
	}))
}

func TestResolveMarkdownSource(t *testing.T) {
	// Test downloading from URL
	t.Run("Download markdown from URL", func(t *testing.T) {
		content := "Test content from URL"
		mockServer := mockHTTPServer(content)
		defer mockServer.Close()

		url := mockServer.URL
		result, err := resolveMarkdownSource(url)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := []byte(content)
		if string(result) != string(expected) {
			t.Errorf("Expected content to be %q, got %q", expected, result)
		}
	})

	// Test reading from local file
	t.Run("Read from a local file", func(t *testing.T) {
		content := "Test content from local file"
		temporaryFile, err := os.CreateTemp("", "example")
		if err != nil {
			t.Fatalf("Error creating temporary file: %v", err)
		}
		defer os.Remove(temporaryFile.Name())

		if _, err := temporaryFile.Write([]byte(content)); err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
		if err := temporaryFile.Close(); err != nil {
			t.Fatalf("Error closing temporary file: %v", err)
		}

		path := temporaryFile.Name()
		result, err := resolveMarkdownSource(path)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := []byte(content)
		if string(result) != string(expected) {
			t.Errorf("Expected content to be %q, got %q", expected, result)
		}
	})

	// Test non-existing file
	t.Run("Non-existing file", func(t *testing.T) {
		nonExistingPath := "non_existing_file.md"
		_, err := resolveMarkdownSource(nonExistingPath)
		if err == nil {
			t.Error("Expected error for non-existing file, but got nil")
		}
		expectedErrorMsg := fmt.Sprintf("markdown file '%s' does not exist", nonExistingPath)
		if !strings.Contains(err.Error(), expectedErrorMsg) {
			t.Errorf("Expected error message to contain %q, got %q", expectedErrorMsg, err.Error())
		}
	})
}

func TestScenarioParsing(t *testing.T) {
	// Test parsing a scenario from markdown
	t.Run("Parse scenario that doesn't have an h1 tag to use for it's title", func(t *testing.T) {
		content := "Test content from local file"
		temporaryFile, err := os.CreateTemp("", "example")
		if err != nil {
			t.Fatalf("Error creating temporary file: %v", err)
		}
		defer os.Remove(temporaryFile.Name())

		if _, err := temporaryFile.Write([]byte(content)); err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
		if err := temporaryFile.Close(); err != nil {
			t.Fatalf("Error closing temporary file: %v", err)
		}

		path := temporaryFile.Name()

		scenario, err := CreateScenarioFromMarkdown(path, []string{"bash"}, nil)

		assert.NoError(t, err)
		fmt.Println(scenario)
		assert.Equal(t, filepath.Base(path), scenario.Name)
	})
}

func TestVariableOverrides(t *testing.T) {
	variableScenarioPath := "../../../scenarios/testing/variables.md"
	// Test overriding environment variables
	t.Run("Override a standard variable declaration", func(t *testing.T) {
		scenario, err := CreateScenarioFromMarkdown(
			variableScenarioPath,
			[]string{"bash"},
			map[string]string{
				"MY_VAR": "my_value",
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, "my_value", scenario.Environment["MY_VAR"])
		assert.Contains(t, scenario.Steps[0].CodeBlocks[0].Content, "export MY_VAR=my_value")
	})

	t.Run(
		"Override a variable that is declared on the same line as another variable, separated by &&",
		func(t *testing.T) {
			scenario, err := CreateScenarioFromMarkdown(
				variableScenarioPath,
				[]string{"bash"},
				map[string]string{
					"NEXT_VAR": "next_value",
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, "next_value", scenario.Environment["NEXT_VAR"])
			assert.Contains(
				t,
				scenario.Steps[1].CodeBlocks[0].Content,
				`export NEXT_VAR=next_value && export OTHER_VAR="Hello, World!"`,
			)
		},
	)

	t.Run(
		"Override a variable that is declared on the same line as another variable, separated by ;",
		func(t *testing.T) {
			scenario, err := CreateScenarioFromMarkdown(
				variableScenarioPath,
				[]string{"bash"},
				map[string]string{
					"THIS_VAR": "this_value",
					"THAT_VAR": "that_value",
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, "this_value", scenario.Environment["THIS_VAR"])
			assert.Equal(t, "that_value", scenario.Environment["THAT_VAR"])
			assert.Contains(
				t,
				scenario.Steps[2].CodeBlocks[0].Content,
				`export THIS_VAR=this_value ; export THAT_VAR=that_value`,
			)
		})

	t.Run("Override a variable that has a subshell command as it's value", func(t *testing.T) {
		scenario, err := CreateScenarioFromMarkdown(
			variableScenarioPath,
			[]string{"bash"},
			map[string]string{
				"SUBSHELL_VARIABLE": "subshell_value",
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, "subshell_value", scenario.Environment["SUBSHELL_VARIABLE"])
		assert.Contains(
			t,
			scenario.Steps[3].CodeBlocks[0].Content,
			`export SUBSHELL_VARIABLE=subshell_value`,
		)
	})

	t.Run("Override a variable that references another variable", func(t *testing.T) {
		scenario, err := CreateScenarioFromMarkdown(
			variableScenarioPath,
			[]string{"bash"},
			map[string]string{
				"VAR2": "var2_value",
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, "var2_value", scenario.Environment["VAR2"])
		assert.Contains(
			t,
			scenario.Steps[4].CodeBlocks[0].Content,
			`export VAR2=var2_value`,
		)
	})
}
