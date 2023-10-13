package environments

const (
	EnvironmentsLocal = "local"
	EnvironmentsCI    = "ci"
	EnvironmentsOCD   = "ocd"
	EnvironmentsAzure = "azure"
)

// Check if the environment is valid.
func IsValidEnvironment(environment string) bool {
	switch environment {
	case EnvironmentsLocal, EnvironmentsCI, EnvironmentsOCD, EnvironmentsAzure:
		return true
	default:
		return false
	}
}

func IsAzureEnvironment(environment string) bool {
	return environment == EnvironmentsAzure || environment == EnvironmentsOCD
}
