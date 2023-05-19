package utils

func CopyMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func MergeMaps(a, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}

	return a
}
