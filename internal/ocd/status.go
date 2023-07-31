package ocd

// / The status of a one-click deployment.
type OneClickDeploymentStatus struct {
	Steps        []string `json:"steps"`
	CurrentStep  int      `json:"currentStep"`
	Status       string   `json:"status"`
	ResourceURIs []string `json:"resourceURIs"`
	Error        string   `json:"error"`
}
