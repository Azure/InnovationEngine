package lib

import (
	"testing"
)

func TestMapUtilities(t *testing.T) {

	t.Run("Copying maps", func(t *testing.T) {
		original := make(map[string]string)
		original["key"] = "value"

		copy := CopyMap(original)

		if len(copy) != 1 {
			t.Errorf("Copy length is wrong: %d", len(copy))
		}

		if copy["key"] != "value" {
			t.Errorf("Copy is wrong: %s", copy["key"])
		}
	})

	t.Run("Merging maps", func(t *testing.T) {
		original := make(map[string]string)
		original["key"] = "value"

		merge := make(map[string]string)
		merge["key2"] = "value2"

		merged := MergeMaps(original, merge)

		if len(merged) != 2 {
			t.Errorf("Merged length is wrong: %d", len(merged))
		}

		if merged["key"] != "value" {
			t.Errorf("Merged is wrong: %s", merged["key"])
		}

		if merged["key2"] != "value2" {
			t.Errorf("Merged is wrong: %s", merged["key2"])
		}
	})

}
