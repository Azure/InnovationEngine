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

// Compares two maps by key and returns the difference between the two.
// This comparison doesn't take into account the value of the key, only the
// presence of the key.
func DiffMapsByKey(a, b map[string]string) map[string]string {
	diff := make(map[string]string)

	if len(a) > len(b) {
		for k, v := range a {
			if b[k] == "" && v != "" {
				diff[k] = v
			}
		}
	} else {
		for k, v := range b {
			if a[k] == "" && v != "" {
				diff[k] = v
			}
		}
	}

	return diff
}

// Compares two maps by key and returns the difference between the two based
// on the value of the key.
func DiffMapsByValue(a, b map[string]string) map[string]string {
	diff := make(map[string]string)

	if len(a) > len(b) {
		for k, v := range a {
			if b[k] != v {
				diff[k] = v
			}
		}
	} else {
		for k, v := range b {
			if a[k] != v {
				diff[k] = v
			}
		}
	}

	return diff
}
