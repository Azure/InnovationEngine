package environments

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/ui"
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

// Print out the one click deployment status if in the correct environment.
func ReportAzureStatus(status OneClickDeploymentStatus, environment string) {
	if environment != EnvironmentsOCD {
		return
	}

	statusJson, err := status.AsJsonString()
	if err != nil {
		logging.GlobalLogger.Error("Failed to marshal status", err)
	} else {
		// We add these strings to the output so that the portal can find and parse
		// the JSON status.
		ocdStatus := fmt.Sprintf("ie_us%sie_ue\n", statusJson)
		fmt.Println(ui.OcdStatusUpdateStyle.Render(ocdStatus))
	}
}

// Attach deployed resource URIs to the one click deployment status if we're in
// the correct environment & we have a resource group name.
func AttachResourceURIsToAzureStatus(
	status *OneClickDeploymentStatus,
	resourceGroupName string,
	environment string,
) {

	if environment != EnvironmentsOCD {
		logging.GlobalLogger.Info(
			"Not fetching resource URIs because we're not in the OCD environment.",
		)
		return
	}

	if resourceGroupName == "" {
		logging.GlobalLogger.Warn("No resource group name found.")
		return
	}

	resourceURIs := az.FindAllDeployedResourceURIs(resourceGroupName)

	if len(resourceURIs) > 0 {
		logging.GlobalLogger.WithField("resourceURIs", resourceURIs).
			Info("Found deployed resources.")
		status.ResourceURIs = resourceURIs
	} else {
		logging.GlobalLogger.Warn("No deployed resources found.")
	}
}
