package ocd

import (
	"encoding/json"

	"github.com/Azure/InnovationEngine/internal/logging"
)

// / The status of a one-click deployment.
type OneClickDeploymentStatus struct {
	Steps        []string `json:"steps"`
	CurrentStep  int      `json:"currentStep"`
	Status       string   `json:"status"`
	ResourceURIs []string `json:"resourceURIs"`
	Error        string   `json:"error"`
}

func NewOneClickDeploymentStatus() OneClickDeploymentStatus {
	return OneClickDeploymentStatus{
		Steps:        []string{},
		CurrentStep:  0,
		Status:       "Executing",
		ResourceURIs: []string{},
		Error:        "",
	}
}

// Get the status as a JSON string.
func (status *OneClickDeploymentStatus) AsJsonString() (string, error) {
	json, err := json.Marshal(status)
	if err != nil {
		logging.GlobalLogger.Error("Failed to marshal status", err)
		return "", err
	}

	return string(json), nil
}

func (status *OneClickDeploymentStatus) AddStep(step string) {
	status.Steps = append(status.Steps, step)
}

func (status *OneClickDeploymentStatus) AddResourceURI(uri string) {
	status.ResourceURIs = append(status.ResourceURIs, uri)
}

func (status *OneClickDeploymentStatus) SetError(err error) {
	status.Status = "Failed"
	status.Error = err.Error()
}
