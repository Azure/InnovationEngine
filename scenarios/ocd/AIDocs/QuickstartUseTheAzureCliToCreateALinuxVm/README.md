---
title: 'Quickstart: Use the Azure CLI to create a Linux VM'
description: In this quickstart, you learn how to use the Azure CLI to create a Linux virtual machine
author: cynthn
ms.service: virtual-machines
ms.collection: linux
ms.topic: quickstart
ms.workload: infrastructure
ms.date: 06/01/2022
ms.author: cynthn
ms.custom: innovation-engine, mvc, seo-javascript-september2019, seo-javascript-october2019, seo-python-october2019, devx-track-azurecli, mode-api
ms.permissions: Microsoft.Resources/deployments/operationStatuses/read, Microsoft.OperationalInsights/workspaces/read, Microsoft.ResourceHealth/availabilityStatuses/read, Microsoft.DevTestLab/labs/virtualmachines/read, Microsoft.Network/virtualNetworks/read Microsoft.Network/networkInterfaces/readWrite Microsoft.Network/networkInterfaces/join/action Microsoft.Network/networkInterfaces/ipConfigurations/readWrite Microsoft.Network/networkInterfaces/ipConfigurations/join/action Microsoft.Network/networkSecurityGroups/read Microsoft.Network/networkSecurityGroups/securityRules/readWrite Microsoft.Network/networkSecurityGroups/securityRules/join/action Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/networkInterfaces/networkSecurityGroups/read Microsoft.Compute/virtualMachines/write, Microsoft.Network/networkInterfaces/read, Microsoft.Authorization/permissions/read, Microsoft.Network/networkInterfaces/write, Microsoft.ManagedIdentity/userAssignedIdentities/read, Microsoft.Network/virtualNetworks/read, Microsoft.Compute/virtualMachines/extensions/read, Microsoft.Authorization/locks/write, Microsoft.AzureAD/userinfo/read, Microsoft.Compute/availabilitySets/read, Microsoft.OperationalInsights/workspaces/sharedKeys/action, Microsoft.Support/supportTickets/read, Microsoft.Network/virtualNetworks/subnets/join/action, Microsoft.Network/networkInterfaces/join/action, Microsoft.Compute/galleries/read, Microsoft.Resources/deployments/read, Microsoft.Compute/virtualMachines/read, Microsoft.Compute/virtualMachines/runCommand/action, Microsoft.Compute/virtualMachines/write, Microsoft.Authorization/roleAssignments/write, Microsoft.ConfigurationManagement/configurationRead/action, Microsoft.Compute/virtualMachines/write Microsoft.Compute/virtualMachines/delete Microsoft.Compute/virtualMachines/start/action Microsoft.Compute/images/read Microsoft.Network/networkInterfaces/write Microsoft.Network/virtualNetworks/read Microsoft.Network/virtualNetworks/subnets/read Microsoft.Network/networkSecurityGroups/write Microsoft.Network/networkSecurityGroups/securityRules/write, Microsoft.Network/publicIPAddresses/read, Microsoft.Storage/storageAccounts/read, Microsoft.Resources/subscriptions/resourceGroups/write, Microsoft.Resources/subscriptions/resourceGroups/delete, Microsoft.Network/publicIPAddresses/write
---

# Quickstart: Create a Linux virtual machine with the Azure CLI

**Applies to:** :heavy_check_mark: Linux VMs

This quickstart shows you how to use the Azure CLI to deploy a Linux virtual machine (VM) in Azure. The Azure CLI is used to create and manage Azure resources via either the command line or scripts.

In this tutorial, we will be installing the latest Debian image. To show the VM in action, you'll connect to it using SSH and install the NGINX web server.

If you don't have an Azure subscription, create a [free account](https://azure.microsoft.com/free/?WT.mc_id=A261C142F) before you begin.

## Launch Azure Cloud Shell

The Azure Cloud Shell is a free interactive shell that you can use to run the steps in this article. It has common Azure tools preinstalled and configured to use with your account. 

To open the Cloud Shell, just select **Try it** from the upper right corner of a code block. You can also open Cloud Shell in a separate browser tab by going to [https://shell.azure.com/bash](https://shell.azure.com/bash). Select **Copy** to copy the blocks of code, paste it into the Cloud Shell, and select **Enter** to run it.

If you prefer to install and use the CLI locally, this quickstart requires Azure CLI version 2.0.30 or later. Run `az --version` to find the version. If you need to install or upgrade, see [Install Azure CLI]( /cli/azure/install-azure-cli).

## Define Environment Variables
Environment variables are commonly used in Linux to centralize configuration data to improve consistency and maintainability of the system. Create the following environment variables to specify the names of resources that will be created later in this tutorial:

```azurecli-interactive
export RESOURCE_GROUP_NAME=myResourceGroup
export LOCATION=eastus
export VM_NAME=myVM
export VM_IMAGE=debian
export ADMIN_USERNAME=azureuser
```

## Create a resource group

Create a resource group with the [az group create](/cli/azure/group) command. An Azure resource group is a logical container into which Azure resources are deployed and managed. 

```azurecli-interactive
az group create --name $RESOURCE_GROUP_NAME --location $LOCATION
```

## Create virtual machine

Create a VM with the [az vm create](/cli/azure/vm) command.

The following example creates a VM and adds a user account. The `--generate-ssh-keys` parameter is used to automatically generate an SSH key, and put it in the default key location (*~/.ssh*). To use a specific set of keys instead, use the `--ssh-key-values` option.

```azurecli-interactive
az vm create \
  --resource-group $RESOURCE_GROUP_NAME \
  --name $VM_NAME \
  --image $VM_IMAGE \
  --admin-username $ADMIN_USERNAME \
  --generate-ssh-keys \
  --public-ip-sku Standard
```

It takes a few minutes to create the VM and supporting resources. The following example output shows the VM create operation was successful.
<!--expected_similarity=0.18-->
```json
{
  "fqdns": "",
  "id": "/subscriptions/<guid>/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM",
  "location": "eastus",
  "macAddress": "00-0D-3A-23-9A-49",
  "powerState": "VM running",
  "privateIpAddress": "10.0.0.4",
  "publicIpAddress": "40.68.254.142",
  "resourceGroup": "myResourceGroup"
}
```

Make a note of the `publicIpAddress` to use later.

You can retrieve and store the IP address in the variable IP_ADDRESS with the following command:

```azurecli-interactive
export IP_ADDRESS=$(az vm show --show-details --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --query publicIps --output tsv)
```

Results:

<!-- expected_similarity=0.3 -->
```json
{"IP_ADDRESS": "XX.XX.XX.XX"}
```

## Install web server 

To see your VM in action, install the NGINX web server. Update your package sources and then install the latest NGINX package. The following command uses run-command to run `sudo apt-get update && sudo apt-get install -y nginx` on the VM:

```azurecli-interactive
az vm run-command invoke \
   --resource-group $RESOURCE_GROUP_NAME \
   --name $VM_NAME \
   --command-id RunShellScript \
   --scripts "sudo apt-get update && sudo apt-get install -y nginx"
```
## Open port 80 for web traffic

By default, only SSH connections are opened when you create a Linux VM in Azure. Use [az vm open-port](/cli/azure/vm) to open TCP port 80 for use with the NGINX web server:

```azurecli-interactive
az vm open-port --port 80 --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME
```

## View the web server in action

Use a web browser of your choice to view the default NGINX welcome page. Use the public IP address of your VM as the web address. The following example shows the default NGINX web site:

![Screenshot showing the N G I N X default web page.](./media/quick-create-cli/nginix-welcome-page-debian.png)

Alternatively, run the following command to see the NGINX welcome page in the terminal

```azurecli-interactive
 curl $IP_ADDRESS
```
 
The following example shows the default NGINX web site in the terminal as successful output:
<!--expected_similarity=0.8-->
```html
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```


## Next steps

In this quickstart, you deployed a simple virtual machine, opened a network port for web traffic, and installed a basic web server. To learn more about Azure virtual machines, continue to the tutorial for Linux VMs.


> [!div class="nextstepaction"]
> [Azure Linux virtual machine tutorials](./tutorial-manage-vm.md)



<details>
<summary><h2>FAQs</h2></summary>

#### Q. What is the command-specific breakdown of permissions needed to implement this doc? 

A. _Format: Commands as they appears in the doc | list of unique permissions needed to run each of those commands_


  - ```azurecli-interactive export RESOURCE_GROUP_NAME=myResourceGroup export LOCATION=eastus export VM_NAME=myVM export VM_IMAGE=debian export ADMIN_USERNAME=azureuser ```

      - Microsoft.Compute/virtualMachines/write Microsoft.Compute/virtualMachines/delete Microsoft.Compute/virtualMachines/start/action Microsoft.Compute/images/read Microsoft.Network/networkInterfaces/write Microsoft.Network/virtualNetworks/read Microsoft.Network/virtualNetworks/subnets/read Microsoft.Network/networkSecurityGroups/write Microsoft.Network/networkSecurityGroups/securityRules/write
  - ```azurecli-interactive az group create --name $RESOURCE_GROUP_NAME --location $LOCATION ```

      - Microsoft.Authorization/roleAssignments/write
      - Microsoft.Resources/subscriptions/resourceGroups/write
  - ```azurecli-interactive az vm create \ --resource-group $RESOURCE_GROUP_NAME \ --name $VM_NAME \ --image $VM_IMAGE \ --admin-username $ADMIN_USERNAME \ --generate-ssh-keys \ --public-ip-sku Standard ```

      - Microsoft.Network/networkInterfaces/write
      - Microsoft.Network/virtualNetworks/subnets/join/action
      - Microsoft.Compute/virtualMachines/write
      - Microsoft.Network/virtualNetworks/read
      - Microsoft.Compute/virtualMachines/extensions/read
      - Microsoft.Compute/availabilitySets/read
      - Microsoft.Network/publicIPAddresses/read
      - Microsoft.Compute/galleries/read
      - Microsoft.Storage/storageAccounts/read
      - Microsoft.Network/networkInterfaces/read
      - Microsoft.Network/publicIPAddresses/write
  - ```azurecli-interactive export IP_ADDRESS=$(az vm show --show-details --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --query publicIps --output tsv) ```

      - Microsoft.Compute/virtualMachines/read
      - Microsoft.Network/networkInterfaces/read
  - ```azurecli-interactive az vm run-command invoke \ --resource-group $RESOURCE_GROUP_NAME \ --name $VM_NAME \ --command-id RunShellScript \ --scripts "sudo apt-get update && sudo apt-get install -y nginx" ```

      - Microsoft.Network/networkInterfaces/write
      - Microsoft.Resources/deployments/operationStatuses/read
      - Microsoft.Support/supportTickets/read
      - Microsoft.ManagedIdentity/userAssignedIdentities/read
      - Microsoft.Compute/virtualMachines/runCommand/action
      - Microsoft.Compute/virtualMachines/write
      - Microsoft.OperationalInsights/workspaces/read
      - Microsoft.ResourceHealth/availabilityStatuses/read
      - Microsoft.Authorization/locks/write
      - Microsoft.DevTestLab/labs/virtualmachines/read
      - Microsoft.Network/networkInterfaces/join/action
      - Microsoft.ConfigurationManagement/configurationRead/action
      - Microsoft.AzureAD/userinfo/read
      - Microsoft.OperationalInsights/workspaces/sharedKeys/action
      - Microsoft.Resources/deployments/read
      - Microsoft.Network/networkInterfaces/read
      - Microsoft.Authorization/permissions/read
  - ```azurecli-interactive az vm open-port --port 80 --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME ```

      - Microsoft.Network/virtualNetworks/read Microsoft.Network/networkInterfaces/readWrite Microsoft.Network/networkInterfaces/join/action Microsoft.Network/networkInterfaces/ipConfigurations/readWrite Microsoft.Network/networkInterfaces/ipConfigurations/join/action Microsoft.Network/networkSecurityGroups/read Microsoft.Network/networkSecurityGroups/securityRules/readWrite Microsoft.Network/networkSecurityGroups/securityRules/join/action Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/networkInterfaces/networkSecurityGroups/read Microsoft.Compute/virtualMachines/write
  - ```azurecli-interactive az group delete --name $RESOURCE_GROUP_NAME --no-wait --yes --verbose ```

      - Microsoft.Resources/subscriptions/resourceGroups/delete

#### Q. What is Azure Cloud Shell? 

A. Azure Cloud Shell is a free interactive shell that allows you to run Azure CLI commands in the browser. It comes with common Azure tools preinstalled and configured to use with your account. You can open the Cloud Shell by selecting 'Try it' from the upper right corner of a code block in the article. Alternatively, you can open it in a separate browser tab by going to [https://shell.azure.com/bash](https://shell.azure.com/bash).


#### Q. How can I install and use Azure CLI locally? 

A. To install and use Azure CLI locally, you need to have Azure CLI version 2.0.30 or later. You can check your current Azure CLI version by running `az --version` command in the command prompt or terminal. If you need to install or upgrade Azure CLI, you can follow the instructions provided in the [Install Azure CLI](https://docs.microsoft.com/cli/azure/install-azure-cli) documentation.


#### Q. How can I create a resource group using Azure CLI? 

A. To create a resource group using Azure CLI, you can use the command `az group create --name <resourceGroupName> --location <location>` where `<resourceGroupName>` is the name you want to give to the resource group, and `<location>` is the Azure region where you want to create the resource group. For example, in the given documentation, the command `az group create --name $RESOURCE_GROUP_NAME --location $LOCATION` is used to create a resource group with the specified name and location.


#### Q. How can I create a virtual machine using Azure CLI? 

A. To create a virtual machine using Azure CLI, you can use the command `az vm create` with the necessary parameters such as `--resource-group`, `--name`, `--image`, and `--admin-username`. For example, in the given documentation, the command `az vm create --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --image $VM_IMAGE --admin-username $ADMIN_USERNAME --generate-ssh-keys --public-ip-sku Standard` is used to create a virtual machine with the specified resource group, name, image, and admin username. The `--generate-ssh-keys` parameter is used to automatically generate an SSH key, and the `--public-ip-sku` parameter is used to specify the SKU for the public IP address associated with the virtual machine.


#### Q. How can I install the NGINX web server on the virtual machine? 

A. To install the NGINX web server on the virtual machine, you can use the `az vm run-command invoke` command with the `--command-id RunShellScript` option and the installation command as the `--scripts` parameter. For example, in the given documentation, the command `az vm run-command invoke --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --command-id RunShellScript --scripts "sudo apt-get update && sudo apt-get install -y nginx"` is used to install NGINX on the virtual machine.


#### Q. How can I open port 80 for web traffic on the virtual machine? 

A. To open port 80 for web traffic on the virtual machine, you can use the `az vm open-port` command with the `--port` option and the resource group and virtual machine name as parameters. For example, in the given documentation, the command `az vm open-port --port 80 --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME` is used to open port 80 on the virtual machine.


#### Q. How can I view the default NGINX web page? 

A. To view the default NGINX web page, you can use a web browser of your choice and enter the public IP address of the virtual machine as the web address. Alternatively, you can use the `curl` command in the terminal and pass the public IP address as a parameter. For example, in the given documentation, the command `curl $IP_ADDRESS` is used to view the default NGINX web page in the terminal.


#### Q. How can I clean up the resources created in this quickstart? 

A. To clean up the resources created in this quickstart, you can use the `az group delete` command with the resource group name as a parameter. For example, in the given documentation, the command `az group delete --name $RESOURCE_GROUP_NAME --no-wait --yes --verbose` is used to delete the resource group and all related resources. Make sure to replace `$RESOURCE_GROUP_NAME` with the actual name of your resource group.


#### Q. Where can I find more tutorials for managing Azure Linux virtual machines? 

A. You can find more tutorials for managing Azure Linux virtual machines in the [Azure Linux virtual machine tutorials](./tutorial-manage-vm.md) documentation. It provides detailed steps and examples for various tasks such as managing VMs, disks, networking, and more.

</details>