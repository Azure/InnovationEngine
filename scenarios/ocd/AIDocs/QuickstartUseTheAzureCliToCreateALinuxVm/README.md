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
ms.permissions: Microsoft.Resources/subscriptions/resourceGroups/read, Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/virtualNetworks/subnets/join/action, Microsoft.Authorization/roleAssignments/read, Microsoft.Network/virtualNetworks/subnets/join/action, Microsoft.Network/networkInterfaces/ipConfigurations/read, Microsoft.Compute/virtualMachines/write, Microsoft.Network/publicIPAddresses/read, Microsoft.Resources/subscriptions/read, Microsoft.Resources/subscriptions/tagNames/read, Microsoft.Network/networkInterfaces/delete, Microsoft.Network/networkInterfaces/join/action, Microsoft.Network/virtualNetworks/delete, Microsoft.Authorization/roleAssignments/write, Microsoft.Network/virtualNetworks/read, Microsoft.Network/virtualNetworks/subnets/read, Microsoft.Compute/virtualMachines/read, Microsoft.Authorization/roleDefinitions/read, Microsoft.Resources/subscriptions/tagNames/write, Microsoft.Network/networkInterfaces/write, Microsoft.Resources/subscriptions/tags/read, Microsoft.Resources/subscriptions/tags/write, Microsoft.Network/networkInterfaces/ipConfigurations/delete, Microsoft.Compute/virtualMachines/runCommands/run/action, Microsoft.Compute/virtualMachines/show/action, Microsoft.Resources/subscriptions/resourceGroups/write, Microsoft.Resources/deployments/write, Microsoft.Network/networkInterfaces/read, Microsoft.Network/virtualNetworks/write, Microsoft.Network/virtualNetworks/subnets/write, Microsoft.Network/networkInterfaces/ipConfigurations/write, Microsoft.Network/publicIPAddresses/delete, Microsoft.Authorization/permissions/read, Microsoft.Insights/eventtypes/values/read, Microsoft.Resources/deployments/read, Microsoft.Network/publicIPAddresses/write
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
{"id":"/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}","name":"{resource-group-name}","location":"{location}","properties":{"provisioningState":"Succeeded"}}
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
    "value": [
        {
            "code": null,
            "displayStatus": "Provisioning",
            "level": "Info",
            "message": "VM has reported its status as provisioned. Subsequent status reports will be reported at the frequency specified in diagnostics settings.",
            "time": "2021-03-01T10:10:00.000000+00:00",
            "type": "Status"
        }
    ]
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
{
  "statusCode": 200,
  "headers": {
    "content-length": "12",
    "content-type": "text/html",
    "date": "Sun, 27 Sep 2020 12:00:00 GMT",
    "server": "nginx/1.14.0 (Ubuntu)"
  },
  "body": "Hello World!"
}
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

      - Microsoft.Network/networkInterfaces/read
      - Microsoft.Network/networkInterfaces/write
      - Microsoft.Network/virtualNetworks/subnets/write
      - Microsoft.Network/networkInterfaces/delete
      - Microsoft.Network/networkInterfaces/join/action
      - Microsoft.Network/virtualNetworks/subnets/join/action
      - Microsoft.Network/publicIPAddresses/read
      - Microsoft.Network/virtualNetworks/write
      - Microsoft.Network/networkInterfaces/ipConfigurations/write
      - Microsoft.Network/virtualNetworks/delete
      - Microsoft.Network/publicIPAddresses/delete
      - Microsoft.Network/networkInterfaces/ipConfigurations/read
      - Microsoft.Network/networkInterfaces/ipConfigurations/delete
      - Microsoft.Network/virtualNetworks/read
      - Microsoft.Network/publicIPAddresses/write
      - Microsoft.Compute/virtualMachines/write
      - Microsoft.Network/virtualNetworks/subnets/read
  - ```azurecli-interactive az group create --name $RESOURCE_GROUP_NAME --location $LOCATION ```

      - Microsoft.Resources/subscriptions/resourceGroups/read
      - Microsoft.Resources/subscriptions/resourceGroups/write
  - ```azurecli-interactive az vm create \ --resource-group $RESOURCE_GROUP_NAME \ --name $VM_NAME \ --image $VM_IMAGE \ --admin-username $ADMIN_USERNAME \ --generate-ssh-keys \ --public-ip-sku Standard ```

      - Microsoft.Authorization/permissions/read
      - Microsoft.Network/publicIPAddresses/write
      - Microsoft.Compute/virtualMachines/write
  - ```azurecli-interactive export IP_ADDRESS=$(az vm show --show-details --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --query publicIps --output tsv) ```

      - Microsoft.Resources/subscriptions/resourceGroups/read
      - Microsoft.Compute/virtualMachines/read
      - Microsoft.Resources/deployments/read
      - Microsoft.Compute/virtualMachines/show/action
      - Microsoft.Network/publicIPAddresses/read
  - ```azurecli-interactive az vm run-command invoke \ --resource-group $RESOURCE_GROUP_NAME \ --name $VM_NAME \ --command-id RunShellScript \ --scripts "sudo apt-get update && sudo apt-get install -y nginx" ```

      - Microsoft.Resources/subscriptions/resourceGroups/read
      - Microsoft.Authorization/roleDefinitions/read
      - Microsoft.Resources/subscriptions/read
      - Microsoft.Resources/subscriptions/tagNames/write
      - Microsoft.Authorization/roleAssignments/read
      - Microsoft.Authorization/roleAssignments/write
      - Microsoft.Resources/deployments/write
      - Microsoft.Resources/subscriptions/tags/read
      - Microsoft.Resources/subscriptions/tagNames/read
      - Microsoft.Resources/subscriptions/tags/write
      - Microsoft.Authorization/permissions/read
      - Microsoft.Compute/virtualMachines/runCommands/run/action
      - Microsoft.Insights/eventtypes/values/read
      - Microsoft.Resources/deployments/read
      - Microsoft.Compute/virtualMachines/read
  - ```azurecli-interactive az vm open-port --port 80 --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME ```

      - Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/virtualNetworks/subnets/join/action Microsoft.Network/virtualNetworks/subnets/join/action
  - ```azurecli-interactive az group delete --name $RESOURCE_GROUP_NAME --no-wait --yes --verbose ```

      - Microsoft.Resources/subscriptions/resourceGroups/write

#### Q. What is the purpose of Azure Cloud Shell? 

A. Azure Cloud Shell is a free interactive shell where you can run Azure CLI commands. It is preconfigured with Azure tools and allows you to run the steps in the tutorial. You can access Azure Cloud Shell by selecting 'Try it' from the upper right corner of a code block in the tutorial or by going to [https://shell.azure.com/bash](https://shell.azure.com/bash) in a separate browser tab. Once in Azure Cloud Shell, you can copy and paste the code blocks from the tutorial and run them.


#### Q. How can I install the Azure CLI locally? 

A. If you prefer to install and use the CLI locally, make sure you have Azure CLI version 2.0.30 or later. You can check the version by running the command 'az --version'. If you need to install or upgrade the Azure CLI, refer to the documentation on [Install Azure CLI](https://docs.microsoft.com/cli/azure/install-azure-cli).


#### Q. What are environment variables and how are they used in this tutorial? 

A. Environment variables are used in Linux to centralize configuration data and improve consistency and maintainability of the system. In this tutorial, environment variables are used to specify the names of resources that will be created later. The following environment variables are defined: RESOURCE_GROUP_NAME, LOCATION, VM_NAME, VM_IMAGE, and ADMIN_USERNAME. These variables will be referenced in the subsequent commands for resource creation.


#### Q. What is an Azure resource group and why do we need to create it? 

A. An Azure resource group is a logical container where Azure resources are deployed and managed. It helps organize resources and manage them as a group instead of individually. In this tutorial, we create a resource group using the Azure CLI command 'az group create'. The resource group will provide a context for deploying the virtual machine and other related resources.


#### Q. How do I create a virtual machine with the Azure CLI? 

A. To create a virtual machine using the Azure CLI, you can use the command 'az vm create'. In this tutorial, the command is used to create a virtual machine with the specified resource group, name, image, admin username, and SSH keys. The command also includes the option to generate SSH keys automatically. After running the command, it may take a few minutes for the virtual machine and supporting resources to be created.


#### Q. How can I install the NGINX web server on the virtual machine? 

A. To install the NGINX web server on the virtual machine, the Azure CLI command 'az vm run-command invoke' is used. This command executes a shell script on the virtual machine, which runs the 'sudo apt-get update' and 'sudo apt-get install -y nginx' commands. This will update the package sources and install the latest NGINX package on the virtual machine.


#### Q. How can I open port 80 for web traffic on the virtual machine? 

A. By default, only SSH connections are allowed when creating a Linux virtual machine in Azure. To open TCP port 80 for use with the NGINX web server, you can use the Azure CLI command 'az vm open-port'. In this tutorial, the command is used to open port 80 for the virtual machine with the specified resource group and name.


#### Q. How can I view the web server in action? 

A. To view the NGINX web server in action, you can use a web browser of your choice and navigate to the public IP address of your virtual machine. The public IP address can be obtained by running the command 'az vm show' with the specified resource group and name, and then querying for the 'publicIps' property. Alternatively, you can use the command 'curl $IP_ADDRESS' in the Azure Cloud Shell to see the NGINX welcome page in the terminal.


#### Q. How do I clean up the resources created in this tutorial? 

A. To clean up the resources created in this tutorial, you can use the Azure CLI command 'az group delete'. This command removes the resource group, virtual machine, and all related resources. In the command, specify the name of the resource group to be deleted. The options '--no-wait' and '--yes' are used to skip confirmation prompts during the deletion process.

</details>