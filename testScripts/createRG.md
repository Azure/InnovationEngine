
<!-- 
```variables
export MY_RESOURCE_GROUP_NAME=myResourceGroup
export MY_LOCATION=eastus
export MY_VM_NAME=myVM
export MY_VM_IMAGE=debian
export MY_ADMIN_USERNAME=azureuser
```
-->

## Create a resource group

Create a resource group with the [az group create](/cli/azure/group) command. An Azure resource group is a logical container into which Azure resources are deployed and managed. The following example creates a resource group named *myResourceGroup* in the *eastus* location:

```bash
az group create --name $MY_RESOURCE_GROUP_NAME --location $MY_LOCATION
```