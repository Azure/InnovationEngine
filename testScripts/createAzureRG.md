## Define Environment Variables
Environment variables are commonly used in Linux to centralize configuration data to improve consistency and maintainability of the system. Create the following environment variables to specify the names of resources that will be created later in this tutorial:

```azurecli-interactive
export RESOURCE_GROUP_NAME=createResourceGroupTest$(printf "%08d" $((RANDOM%100000000)))
export LOCATION=eastus
```

## Create a resource group

Create a resource group with the [az group create](/cli/azure/group) command. An Azure resource group is a logical container into which Azure resources are deployed and managed. 

```azurecli-interactive
az group create --name $RESOURCE_GROUP_NAME --location $LOCATION
```
<!--expected_similarity=0.8-->
```json
{
  "id": "/subscriptions/abcdefhijklmnopqrstuvwxyz123456789/resourceGroups/createResourceGroupTest",
  "location": "eastus",
  "managedBy": null,
  "name": "createResourceGroupTest00000000",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```

## Delete the resource group

```azurecli-interactive
az group delete --name $RESOURCE_GROUP_NAME --no-wait --yes --verbose
```
