package lib

import "fmt"

// Makes a copy of a map
func CopyMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Merge two maps together.
func MergeMaps(a, b map[string]string) map[string]string {
	merged := CopyMap(a)
	for k, v := range b {
		merged[k] = v
	}

	return merged
}

// Compares two maps by key and returns the difference between the two.
func DiffMapsByKey(a, b map[string]string) map[string]string {
	diff := make(map[string]string)

	aLength := len(a)
	bLength := len(b)

	if aLength > bLength {
		for k, v := range a {
			if b[k] == "" && v != "" {
				fmt.Print("diff: ", k, v)
				diff[k] = v
			}
		}
	} else {
		for k, v := range b {
			if a[k] == "" && v != "" {
				fmt.Print("diff: ", k, v)
				diff[k] = v
			}
		}
	}

	return diff
}

func DiffMapsByValue(a, b map[string]string) map[string]string {
	diff := make(map[string]string)

	for k, v := range a {
		if b[k] != v {
			diff[k] = v
		}
	}

	return diff
}
