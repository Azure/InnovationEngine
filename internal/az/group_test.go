package az

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindingResourceGroups(t *testing.T) {
	testCases := []struct {
		resourceGroupString       string
		expectedResourceGroupName string
	}{
		// RG string that ends with a slash
		{
			resourceGroupString:       "resourceGroups/rg1/",
			expectedResourceGroupName: "rg1",
		},
		// RG string that ends with a space and starts new text
		{
			resourceGroupString:       "resourceGroups/rg1 /subscriptions/",
			expectedResourceGroupName: "rg1",
		},

		// RG string that includes nested resources and extraneous text.
		{
			resourceGroupString:       "/subscriptions/9b70acd9-975f-44ba-bad6-255a2c8bda37/resourceGroups/myResourceGroup-rg/providers/Microsoft.ContainerRegistry/registries/mydnsrandomnamebbbhe  ffc55a9e-ed2a-4b60-b034-45228dfe7db5  2024-06-11T09:41:36.631310+00:00",
			expectedResourceGroupName: "myResourceGroup-rg",
		},
		// RG string that is surrounded by quotes.
		{
			resourceGroupString:       `"id": "/subscriptions/0a2c89a7-a44e-4cd0-b6ec-868432ad1d13/resourceGroups/myResourceGroup"`,
			expectedResourceGroupName: "myResourceGroup",
		},
		// RG string that has no match.
		{
			resourceGroupString:       "NoMatch",
			expectedResourceGroupName: "",
		},
	}

	for _, tc := range testCases {
		resourceGroupName := FindResourceGroupName(tc.resourceGroupString)
		assert.Equal(t, tc.expectedResourceGroupName, resourceGroupName)
	}
}
