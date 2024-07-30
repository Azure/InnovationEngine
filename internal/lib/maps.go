package lib

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

// Returns the difference between two maps.
func DiffMaps(a, b map[string]string) map[string]string {
	diff := make(map[string]string)
	for k, v := range a {
		if b[k] != v {
			diff[k] = v
		}
	}

	return diff
}
