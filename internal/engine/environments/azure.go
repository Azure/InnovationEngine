package environments

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/InnovationEngine/internal/az"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/ui"
)

// codeblock metadata needed for learn mode deployments.
type AzureCodeBlock struct {
	Description string `json:"description"`
	Command     string `json:"command"`
}

// Step metadata needed for learn mode deployments.
type AzureStep struct {
	Name       string           `json:"name"`
	CodeBlocks []AzureCodeBlock `json:"codeblocks"`
}

// The status of a one-click deployment or learn mode deployment.
type AzureDeploymentStatus struct {
	Steps        []AzureStep `json:"steps"`
	CurrentStep  int         `json:"currentStep"`
	Status       string      `json:"status"`
	ResourceURIs []string    `json:"resourceURIs"`
	Error        string      `json:"error"`
	Output       string      `json:"output"`
}

func NewAzureDeploymentStatus() AzureDeploymentStatus {
	return AzureDeploymentStatus{
		Steps:        []AzureStep{},
		CurrentStep:  0,
		Status:       "Executing",
		ResourceURIs: []string{},
		Error:        "",
	}
}

// Get the status as a JSON string.
func (status *AzureDeploymentStatus) AsJsonString() (string, error) {
	json, err := json.Marshal(status)
	if err != nil {
		logging.GlobalLogger.Error("Failed to marshal status", err)
		return "", err
	}

	return string(json), nil
}

func (status *AzureDeploymentStatus) AddStep(step string, codeBlocks []AzureCodeBlock) {
	status.Steps = append(status.Steps, AzureStep{
		Name:       step,
		CodeBlocks: codeBlocks,
	})
}

func (status *AzureDeploymentStatus) AddResourceURI(uri string) {
	status.ResourceURIs = append(status.ResourceURIs, uri)
}

func (status *AzureDeploymentStatus) SetError(err error) {
	status.Status = "Failed"
	status.Error = err.Error()
}

func (status *AzureDeploymentStatus) SetOutput(output string) {
	status.Output = output
}

// Print out the status JSON for azure/cloudshell if in the correct environment.
func ReportAzureStatus(status AzureDeploymentStatus, environment string) {
	if !IsAzureEnvironment(environment) {
		return
	}

	statusJson, err := status.AsJsonString()
	if err != nil {
		logging.GlobalLogger.Error("Failed to marshal status", err)
	} else {
		// We add these strings to the output so that the portal can find and parse
		// the JSON status.
		ocdStatus := fmt.Sprintf("ie_us%sie_ue", statusJson)
		fmt.Println(ui.OcdStatusUpdateStyle.Render(ocdStatus))
	}
}

// Same as ReportAzureStatus, but returns the status string instead of printing it.
func GetAzureStatus(status AzureDeploymentStatus, environment string) string {
	if !IsAzureEnvironment(environment) {
		return ""
	}

	statusJson, err := status.AsJsonString()
	if err != nil {
		logging.GlobalLogger.Error("Failed to marshal status", err)
		return ""
	} else {
		// We add these strings to the output so that the portal can find and parse
		// the JSON status.
		ocdStatus := fmt.Sprintf("ie_us%sie_ue", statusJson)
		return ocdStatus
	}
}

// Attach deployed resource URIs to the one click deployment status if we're in
// the correct environment & we have a resource group name.
func AttachResourceURIsToAzureStatus(
	status *AzureDeploymentStatus,
	resourceGroupName string,
	environment string,
) {
	if !IsAzureEnvironment(environment) {
		logging.GlobalLogger.Info(
			"Not fetching resource URIs because we're not in the OCD environment.",
		)
	}

	if resourceGroupName == "" {
		logging.GlobalLogger.Warn("No resource group name found.")
		return
	}

	resourceURIs := az.FindAllDeployedResourceURIs(resourceGroupName)

	if len(resourceURIs) > 0 {
		logging.GlobalLogger.WithField("resourceURIs", resourceURIs).
			Info("Found deployed resources.")
		status.ResourceURIs = append(status.ResourceURIs, resourceURIs...)
	} else {
		logging.GlobalLogger.Warn("No deployed resources found.")
	}
}
