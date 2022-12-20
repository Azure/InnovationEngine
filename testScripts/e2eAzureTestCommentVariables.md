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
ms.custom: mvc, seo-javascript-september2019, seo-javascript-october2019, seo-python-october2019, devx-track-azurecli, mode-api
---

<!-- 
```variables
export MY_RESOURCE_GROUP_NAME=myResourceGroup22323
export MY_LOCATION=eastus
export MY_VM_NAME=myVM
export MY_VM_IMAGE=debian
export MY_ADMIN_USERNAME=azureuser
```
-->

# Quickstart: Create a Linux virtual machine with the Azure CLI

**Applies to:** :heavy_check_mark: Linux VMs

This quickstart shows you how to use the Azure CLI to deploy a Linux virtual machine (VM) in Azure. The Azure CLI is used to create and manage Azure resources via either the command line or scripts.

In this tutorial, we will be installing the latest Debian image. To show the VM in action, you'll connect to it using SSH and install the NGINX web server.

If you don't have an Azure subscription, create a [free account](https://azure.microsoft.com/free/?WT.mc_id=A261C142F) before you begin.

## Launch Azure Cloud Shell

The Azure Cloud Shell is a free interactive shell that you can use to run the steps in this article. It has common Azure tools preinstalled and configured to use with your account. 

To open the Cloud Shell, just select **Try it** from the upper right corner of a code block. You can also open Cloud Shell in a separate browser tab by going to [https://shell.azure.com/bash](https://shell.azure.com/bash). Select **Copy** to copy the blocks of code, paste it into the Cloud Shell, and select **Enter** to run it.

If you prefer to install and use the CLI locally, this quickstart requires Azure CLI version 2.0.30 or later. Run `az --version` to find the version. If you need to install or upgrade, see [Install Azure CLI]( /cli/azure/install-azure-cli).



## Create a resource group

Create a resource group with the [az group create](/cli/azure/group) command. An Azure resource group is a logical container into which Azure resources are deployed and managed. The following example creates a resource group named *myResourceGroup* in the *eastus* location:

```bash
az group create --name $MY_RESOURCE_GROUP_NAME --location $MY_LOCATION
```

## Create virtual machine

Create a VM with the [az vm create](/cli/azure/vm) command.

The following example creates a VM named *myVM* and adds a user account named *azureuser*. The `--generate-ssh-keys` parameter is used to automatically generate an SSH key, and put it in the default key location (*~/.ssh*). To use a specific set of keys instead, use the `--ssh-key-values` option.

```bash
az vm create \
  --resource-group $MY_RESOURCE_GROUP_NAME \
  --name $MY_VM_NAME \
  --image $MY_VM_IMAGE \
  --admin-username $MY_ADMIN_USERNAME \
  --generate-ssh-keys
```

It takes a few minutes to create the VM and supporting resources. The following example output shows the VM create operation was successful.
<!--expected_similarity=0.2-->
```Output
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

## Install web server

To see your VM in action, install the NGINX web server. Update your package sources and then install the latest NGINX package.

```bash
az vm run-command invoke \
   -g $MY_RESOURCE_GROUP_NAME \
   -n $MY_VM_NAME \
   --command-id RunShellScript \
   --scripts "sudo apt-get update && sudo apt-get install -y nginx"
```

## Open port 80 for web traffic

By default, only SSH connections are opened when you create a Linux VM in Azure. Use [az vm open-port](/cli/azure/vm) to open TCP port 80 for use with the NGINX web server:

```bash
az vm open-port --port 80 --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_VM_NAME
```

## View the web server in action

Use a web browser of your choice to view the default NGINX welcome page. Use the public IP address of your VM as the web address. The following example shows the default NGINX web site:

![Screenshot showing the N G I N X default web page.](./media/quick-create-cli/nginix-welcome-page-debian.png)

Or Run the following command to see the NGINX welcome page in terminal

```bash
 curl $(az vm show -d -g $MY_RESOURCE_GROUP_NAME -n $MY_VM_NAME --query "publicIps" -o tsv)
```

<!--expected_similarity=0.8-->
```HTML
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

## Clean up resources

When no longer needed, you can use the [az group delete](/cli/azure/group) command to remove the resource group, VM, and all related resources. 

```bash
az group delete --name $MY_RESOURCE_GROUP_NAME --no-wait --yes --verbose
```

## Next steps

In this quickstart, you deployed a simple virtual machine, opened a network port for web traffic, and installed a basic web server. To learn more about Azure virtual machines, continue to the tutorial for Linux VMs.


> [!div class="nextstepaction"]
> [Azure Linux virtual machine tutorials](./tutorial-manage-vm.md)
