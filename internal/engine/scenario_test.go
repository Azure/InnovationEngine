package engine

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
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
