package az

import (
	"fmt"

	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/patterns"
	"github.com/Azure/InnovationEngine/internal/shells"
)

// Find all the deployed resources in a resource group.
func FindAllDeployedResourceURIs(resourceGroup string) []string {
	output, err := shells.ExecuteBashCommand(
		"az resource list -g "+resourceGroup,
		shells.BashCommandConfiguration{
			EnvironmentVariables: map[string]string{},
			InheritEnvironment:   true,
			InteractiveCommand:   false,
			WriteToHistory:       true,
		},
	)
	if err != nil {
		logging.GlobalLogger.Error("Failed to list deployments", err)
	}

	matches := patterns.AzResourceURI.FindAllStringSubmatch(output.StdOut, -1)
	results := []string{}
	for _, match := range matches {
		results = append(results, match[1])
	}
	return results
}

// Find the resource group name from the output of an az command.
func FindResourceGroupName(commandOutput string) string {
	matches := patterns.AzResourceGroupName.FindStringSubmatch(commandOutput)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func BuildResourceGroupId(subscription string, resourceGroup string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscription, resourceGroup)
}
