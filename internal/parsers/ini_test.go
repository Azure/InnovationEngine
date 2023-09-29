package parsers

import (
	"os"
	"testing"
)

func TestParsingINIFiles(t *testing.T) {

	t.Run("INI with valid contents", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test")

		if err != nil {
			t.Errorf("Error creating temp file: %s", err)
		}

		defer os.Remove(tempFile.Name())

		contents := []byte(`[section]
      key=value`)

		if _, err := tempFile.Write(contents); err != nil {
			t.Errorf("Error writing to temp file: %s", err)
		}

		data, err := ParseINIFile(tempFile.Name())

		if err != nil {
			t.Errorf("Error parsing INI file: %s", err)
		}

		if len(data) != 1 {
			t.Errorf("Data length is wrong: %d", len(data))
		}

		if data["key"] != "value" {
			t.Errorf("Data is wrong: %s", data["key"])
		}
	})

}
