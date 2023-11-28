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
ms.permissions: Microsoft.Compute/virtualMachines/read, Microsoft.Network/virtualNetworks/read, Microsoft.Network/networkInterfaces/write, Microsoft.Network/virtualNetworks/subnets/join/action, Microsoft.Authorization/permissions/action, Microsoft.Compute/virtualMachines/runCommand/write, Microsoft.Compute/virtualMachines/instanceView/read, Microsoft.Authorization/permissions/read, Microsoft.Authorization/permissions/write, Microsoft.Authorization/permissions/delete, Microsoft.Compute/virtualMachines/write, Microsoft.Authorization/permissions/write Microsoft.Authorization/permissions/read, Microsoft.Storage/storageAccounts/listKeys/action, Microsoft.Storage/storageAccounts/write, Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/virtualNetworks/join/action, Microsoft.Network/virtualNetworks/write, Microsoft.Resources/subscriptions/resourceGroups/read, Microsoft.Network/publicIPAddresses/read, Microsoft.Network/publicIPAddresses/write
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

Results:

<!-- expected_similarity=0.3 -->
```json
{
    "id": "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}",
    "location": "westus",
    "managedBy": null,
    "name": "{resourceGroupName}",
    "properties": {
        "provisioningState": "Succeeded"
    },
    "tags": null,
    "type": "Microsoft.Resources/resourceGroups"
}
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

Results:

<!-- expected_similarity=0.3 -->
```json
{
  "StatusCode": "Creating",
  "VMId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "VMName": "examplevm",
  "Image": "ubuntu",
  "AdminUsername": "admin",
  "GenerateSSHKeys": true,
  "PublicIpSku": "Standard"
}
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
{"IP_ADDRESS":"x.x.x.x"}
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

Results:

<!-- expected_similarity=0.3 -->
```json
{
    "id": "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{vmName}/extensions/RunCommand",
    "name": "RunCommand",
    "type": "Microsoft.Compute/virtualMachines/extensions",
    "location": "{location}",
    "tags": {
        "$hidden-softdelete-last-modified": "2021-06-22T08:14:34.6612888Z",
        "$hidden-softdelete-tombstone": "2021-06-22T08:14:34.6612888Z"
    },
    "properties": {
        "publisher": "Microsoft.Azure.Extensions",
        "type": "RunShellScript",
        "typeHandlerVersion": "2.0",
        "autoUpgradeMinorVersion": true,
        "settings": {
            "fileUris": [],
            "commandToExecute": "sudo apt-get update && sudo apt-get install -y nginx"
        },
        "protectedSettings": {
            "storageAccountName": "",
            "storageAccountKey": ""
        },
        "provisioningState": "Running",
        "instanceView": {
            "name": "RunShellScript",
            "type": "CustomScriptExtension"
        }
    }
}
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

Results:

<!-- expected_similarity=0.3 -->
```json
{ "message": "curl: (6) Could not resolve host: $IP_ADDRESS" }
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

      - Microsoft.Compute/virtualMachines/write
      - Microsoft.Network/virtualNetworks/write
      - Microsoft.Network/networkInterfaces/write
      - Microsoft.Storage/storageAccounts/listKeys/action
      - Microsoft.Storage/storageAccounts/write
      - Microsoft.Network/publicIPAddresses/write
  - ```azurecli-interactive az group create --name $RESOURCE_GROUP_NAME --location $LOCATION ```

      - Microsoft.Authorization/permissions/write Microsoft.Authorization/permissions/read
  - ```azurecli-interactive az vm create \ --resource-group $RESOURCE_GROUP_NAME \ --name $VM_NAME \ --image $VM_IMAGE \ --admin-username $ADMIN_USERNAME \ --generate-ssh-keys \ --public-ip-sku Standard ```

      - Microsoft.Network/networkInterfaces/write
      - Microsoft.Network/publicIPAddresses/write
      - Microsoft.Compute/virtualMachines/write
      - Microsoft.Network/virtualNetworks/subnets/join/action
  - ```azurecli-interactive export IP_ADDRESS=$(az vm show --show-details --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --query publicIps --output tsv) ```

      - Microsoft.Network/publicIPAddresses/read
      - Microsoft.Compute/virtualMachines/read
      - Microsoft.Resources/subscriptions/resourceGroups/read
  - ```azurecli-interactive az vm run-command invoke \ --resource-group $RESOURCE_GROUP_NAME \ --name $VM_NAME \ --command-id RunShellScript \ --scripts "sudo apt-get update && sudo apt-get install -y nginx" ```

      - Microsoft.Compute/virtualMachines/read
      - Microsoft.Compute/virtualMachines/instanceView/read
      - Microsoft.Authorization/permissions/read
      - Microsoft.Resources/subscriptions/resourceGroups/read
      - Microsoft.Compute/virtualMachines/runCommand/write
  - ```azurecli-interactive az vm open-port --port 80 --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME ```

      - Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/virtualNetworks/join/action
  - ```azurecli-interactive curl $IP_ADDRESS ```

      - Microsoft.Network/publicIPAddresses/read
      - Microsoft.Network/virtualNetworks/read
  - ```azurecli-interactive az group delete --name $RESOURCE_GROUP_NAME --no-wait --yes --verbose ```

      - Microsoft.Authorization/permissions/delete
      - Microsoft.Authorization/permissions/write
      - Microsoft.Authorization/permissions/read
      - Microsoft.Authorization/permissions/action

#### Q. What is the purpose of this document? 

A. This document provides a quickstart guide on how to use the Azure CLI to create a Linux virtual machine (VM) in Azure.


#### Q. How can I access Azure Cloud Shell? 

A. You can access Azure Cloud Shell by selecting 'Try it' from the upper right corner of a code block in the document. Alternatively, you can open Cloud Shell in a separate browser tab by going to [https://shell.azure.com/bash](https://shell.azure.com/bash).


#### Q. What should I do if I prefer to install and use the CLI locally? 

A. If you prefer to install and use the CLI locally, make sure you have Azure CLI version 2.0.30 or later installed. You can check the version by running `az --version`. To install or upgrade Azure CLI, refer to the [Install Azure CLI](https://docs.microsoft.com/cli/azure/install-azure-cli) documentation.


#### Q. How do I create a resource group using the Azure CLI? 

A. You can create a resource group using the `az group create` command. For example: `az group create --name <resourceGroupName> --location <location>`. Replace `<resourceGroupName>` with the desired name for your resource group and `<location>` with the desired location for your resources. This command creates a logical container for your Azure resources.


#### Q. What should I do if the resource group creation fails? 

A. If the resource group creation fails, ensure that the provided resource group name is valid and not already in use. Additionally, check if you have the necessary permissions to create a resource group. Refer to the [az group create](https://docs.microsoft.com/cli/azure/group#create) documentation for more details.


#### Q. How do I create a virtual machine using the Azure CLI? 

A. You can create a virtual machine using the `az vm create` command. For example: 
```
az vm create --resource-group <resourceGroupName> --name <vmName> --image <vmImage> --admin-username <adminUsername> --generate-ssh-keys --public-ip-sku Standard
```
Replace `<resourceGroupName>` with the name of the resource group, `<vmName>` with the desired name for your virtual machine, `<vmImage>` with the desired image for your virtual machine (e.g., `debian`), and `<adminUsername>` with the desired username for the virtual machine's admin. This command creates a virtual machine in the specified resource group.


#### Q. What should I do if the virtual machine creation fails? 

A. If the virtual machine creation fails, ensure that you have provided valid values for all the required parameters. Double-check if the resource group exists and you have sufficient permissions to create a virtual machine. Also, make sure the specified image is available in the desired location. Refer to the [az vm create](https://docs.microsoft.com/cli/azure/vm/create) documentation for more details.


#### Q. How do I install the NGINX web server on the virtual machine? 

A. To install the NGINX web server on the virtual machine, use the `az vm run-command invoke` command. For example: 
```
az vm run-command invoke --resource-group <resourceGroupName> --name <vmName> --command-id RunShellScript --scripts "sudo apt-get update && sudo apt-get install -y nginx"
```
Replace `<resourceGroupName>` with the name of the resource group and `<vmName>` with the name of the virtual machine. This command runs the specified shell script on the virtual machine to install NGINX.


#### Q. What should I do if the NGINX installation fails? 

A. If the NGINX installation fails, ensure that the virtual machine is running and accessible. Double-check the command used to invoke the run-command and make sure there are no syntax errors or typos. Verify that the virtual machine has internet connectivity to download and install the NGINX package. Refer to the [az vm run-command invoke](https://docs.microsoft.com/cli/azure/vm/run-command/invoke) documentation for more details.


#### Q. How do I open port 80 for web traffic on the virtual machine? 

A. To open TCP port 80 for web traffic on the virtual machine, use the `az vm open-port` command. For example: `az vm open-port --port 80 --resource-group <resourceGroupName> --name <vmName>`. Replace `<resourceGroupName>` with the name of the resource group and `<vmName>` with the name of the virtual machine. This command allows inbound traffic on port 80 for the specified virtual machine.


#### Q. What should I do if opening port 80 fails? 

A. If opening port 80 fails, ensure that the specified resource group and virtual machine exist, and you have the necessary permissions to modify the virtual machine's network configuration. Verify that there are no conflicting network security group rules blocking port 80. Refer to the [az vm open-port](https://docs.microsoft.com/cli/azure/vm/open-port) documentation for more details.


#### Q. How can I view the NGINX web server in action? 

A. To view the NGINX web server in action, you can either use a web browser or run a `curl` command. To use a web browser, enter the public IP address of your virtual machine in the address bar. It should display the default NGINX welcome page. Alternatively, you can run `curl <IP_ADDRESS>` in the terminal, where `<IP_ADDRESS>` is the public IP address of your virtual machine. This will show the HTML response of the NGINX welcome page.


#### Q. What should I do if I get a 'Could not resolve host' error when running the `curl` command? 

A. If you get a 'Could not resolve host' error when running the `curl` command, make sure you have replaced `<IP_ADDRESS>` with the actual public IP address of your virtual machine. Verify that the virtual machine is running and accessible, and the public IP address is correct. If the issue persists, check your network configuration and DNS settings. Refer to the [az vm run-command invoke](https://docs.microsoft.com/cli/azure/vm/run-command/invoke) documentation for more details.


#### Q. How do I clean up the created resources? 

A. To clean up the created resources, use the `az group delete` command. For example: `az group delete --name <resourceGroupName> --no-wait --yes --verbose`. Replace `<resourceGroupName>` with the name of the resource group you used. This command deletes the specified resource group, including the virtual machine and all related resources.


#### Q. What should I do if the resource deletion fails? 

A. If the resource deletion fails, ensure that you have provided the correct name of the resource group. Verify that you have the necessary permissions to delete the resource group and its resources. If any resources in the group are still in use or have dependencies, the deletion may fail. Refer to the [az group delete](https://docs.microsoft.com/cli/azure/group/delete) documentation for more details.


#### Q. What are the next steps after following this quickstart? 

A. After following this quickstart, you can explore more Azure Linux virtual machine tutorials by visiting the [Azure Linux virtual machine tutorials](./tutorial-manage-vm.md) page. It provides step-by-step instructions on various aspects of managing Linux virtual machines in Azure.

</details>