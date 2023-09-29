package az

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/logging"
)

// If the correlation ID is set, we need to set the AZURE_HTTP_USER_AGENT
// environment variable so that the Azure CLI will send the correlation ID
// with Azure Resource Manager requests.
func SetCorrelationId(correlationId string, env map[string]string) {
	if correlationId != "" {
		env["AZURE_HTTP_USER_AGENT"] = fmt.Sprintf("innovation-engine-%s", correlationId)
		logging.GlobalLogger.Info("Resource tracking enabled. Tracking ID: " + env["AZURE_HTTP_USER_AGENT"])
	}
}
