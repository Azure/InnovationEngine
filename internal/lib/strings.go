package lib

// Checks if a given string is a number.
func IsNumber(str string) bool {
	for _, r := range str {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
