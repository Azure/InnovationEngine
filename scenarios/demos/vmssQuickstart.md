---
title: Quickstart - Create a Virtual Machine Scale Set with Azure CLI
description: Get started with your deployments by learning how to quickly create a Virtual Machine Scale Set with Azure CLI.
author: ju-shim
ms.author: jushiman
ms.topic: quickstart
ms.service: virtual-machine-scale-sets
ms.date: 11/22/2022
ms.reviewer: mimckitt
ms.custom: mimckitt, devx-track-azurecli, mode-api
---

# Quickstart: Create a Virtual Machine Scale Set with the Azure CLI

**Applies to:** :heavy_check_mark: Linux VMs :heavy_check_mark: Windows VMs :heavy_check_mark: Uniform scale sets

> [!NOTE]
> The following article is for Uniform Virtual Machine Scale Sets. We recommend using Flexible Virtual Machine Scale Sets for new workloads. Learn more about this new orchestration mode in our [Flexible Virtual Machine Scale Sets overview](flexible-virtual-machine-scale-sets.md).

A Virtual Machine Scale Set allows you to deploy and manage a set of auto-scaling virtual machines. You can scale the number of VMs in the scale set manually, or define rules to autoscale based on resource usage like CPU, memory demand, or network traffic. An Azure load balancer then distributes traffic to the VM instances in the scale set. In this quickstart, you create a Virtual Machine Scale Set and deploy a sample application with the Azure CLI.

[!INCLUDE [quickstarts-free-trial-note](../../includes/quickstarts-free-trial-note.md)]

[!INCLUDE [azure-cli-prepare-your-environment.md](~/articles/reusable-content/azure-cli/azure-cli-prepare-your-environment.md)]

- This article requires version 2.0.29 or later of the Azure CLI. If using Azure Cloud Shell, the latest version is already installed. 


## Define Environment Variables

Throughout this document we use environment variables to facilitate cut and paste reuse. 
The default values below will enable you to work through this document in most cases. The meaning of each 
environment variable will be addressed as they are used in the steps below.

```azurecli-interactive
export RESOURCE_GROUP_NAME=vmssQuickstartRG
export RESOURCE_LOCATION=eastus
export SCALE_SET_NAME=vmssQuickstart
export BASE_VM_IMAGE=UbuntuLTS
export ADMIN_USERNAME=azureuser
export LOAD_BALANCER_NAME=vmssQuickstartLB 
export BACKEND_POOL_NAME=vmssQuickstartPool
export LOAD_BALANCER_RULE_NAME=vmssQuickstartRule
export FRONT_END_IP_NAME=vmssQuickstartLoadBalancerFrontEnd 
export CUSTOM_SCRIPT_NAME=vmssQuickstartCustomScript
export SCALE_SET_PUBLIC_IP=vmssQuickstartPublicIP
```

## Create a scale set
Before you can create a scale set, create a resource group with [az group create](/cli/azure/group). The following example creates a resource group named *myResourceGroup* in the *eastus* location:

```azurecli-interactive
az group create --name $RESOURCE_GROUP_NAME --location $RESOURCE_LOCATION
```

<!--expected_similarity=0.2-->
```Output
{
  "id": "/subscriptions/<guid>/resourceGroups/myResourceGroup",
  "location": "eastus",
  "managedBy": null,
  "name": "myResourceGroup",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```

Now create a Virtual Machine Scale Set with [az vmss create](/cli/azure/vmss). The following example creates a scale set named *myScaleSet* that is set to automatically update as changes are applied, and generates SSH keys if they do not exist in *~/.ssh/id_rsa*. These SSH keys are used if you need to log in to the VM instances. To use an existing set of SSH keys, instead use the `--ssh-key-value` parameter and specify the location of your keys.

```azurecli-interactive
az vmss create \
  --resource-group $RESOURCE_GROUP_NAME \
  --name $SCALE_SET_NAME \
  --image $BASE_VM_IMAGE \
  --upgrade-policy-mode automatic \
  --admin-username $ADMIN_USERNAME \
  --generate-ssh-keys
```

It takes a few minutes to create and configure all the scale set resources and VMs.


## Deploy sample application
To test your scale set, install a basic web application. The Azure Custom Script Extension is used to download and run a script that installs an application on the VM instances. This extension is useful for post deployment configuration, software installation, or any other configuration / management task. For more information, see the [Custom Script Extension overview](../virtual-machines/extensions/custom-script-linux.md).

Use the Custom Script Extension to install a basic NGINX web server. Apply the Custom Script Extension that installs NGINX with [az vmss extension set](/cli/azure/vmss/extension) as follows:

```azurecli-interactive
az vmss extension set \
  --publisher Microsoft.Azure.Extensions \
  --version 2.0 \
  --name $CUSTOM_SCRIPT_NAME \
  --resource-group $RESOURCE_GROUP_NAME \
  --vmss-name $SCALE_SET_NAME \
  --settings '{"fileUris":["https://raw.githubusercontent.com/Azure-Samples/compute-automation-configurations/master/automate_nginx.sh"],"commandToExecute":"./automate_nginx.sh"}'
```
<!--expected_similarity=0.25-->
```Output
{
  "vmss": {
    "doNotRunExtensionsOnOverprovisionedVMs": false,
    "orchestrationMode": "Uniform",
    "overprovision": true,
    "provisioningState": "Succeeded",
    "singlePlacementGroup": true,
    "timeCreated": "2023-02-01T22:17:20.1117742+00:00",
    "uniqueId": "38328143-69e8-4a9b-9d55-8a404cdb6d8b",
    "upgradePolicy": {
      "mode": "Automatic",
      "rollingUpgradePolicy": {
        "maxBatchInstancePercent": 20,
        "maxSurge": false,
        "maxUnhealthyInstancePercent": 20,
        "maxUnhealthyUpgradedInstancePercent": 20,
        "pauseTimeBetweenBatches": "PT0S",
        "rollbackFailedInstancesOnPolicyBreach": false
      }
    },
    "virtualMachineProfile": {
      "networkProfile": {
        "networkInterfaceConfigurations": [
          {
            "name": "mysca2132Nic",
            "properties": {
              "disableTcpStateTracking": false,
              "dnsSettings": {
                "dnsServers": []
              },
              "enableAcceleratedNetworking": false,
              "enableIPForwarding": false,
              "ipConfigurations": [
                {
                  "name": "mysca2132IPConfig",
                  "properties": {
                    "loadBalancerBackendAddressPools": [
                      {
                        "id": "/subscriptions/f7a60fca-9977-4899-b907-005a076adbb6/resourceGroups/myResourceGroup/providers/Microsoft.Network/loadBalancers/myScaleSetLB/backendAddressPools/myScaleSetLBBEPool",
                        "resourceGroup": "myResourceGroup"
                      }
                    ],
                    "loadBalancerInboundNatPools": [
                      {
                        "id": "/subscriptions/f7a60fca-9977-4899-b907-005a076adbb6/resourceGroups/myResourceGroup/providers/Microsoft.Network/loadBalancers/myScaleSetLB/inboundNatPools/myScaleSetLBNatPool",
                        "resourceGroup": "myResourceGroup"
                      }
                    ],
                    "privateIPAddressVersion": "IPv4",
                    "subnet": {
                      "id": "/subscriptions/f7a60fca-9977-4899-b907-005a076adbb6/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myScaleSetVNET/subnets/myScaleSetSubnet",
                      "resourceGroup": "myResourceGroup"
                    }
                  }
                }
              ],
              "primary": true
            }
          }
        ]
      },
      "osProfile": {
        "adminUsername": "azureuser",
        "allowExtensionOperations": true,
        "computerNamePrefix": "mysca2132",
        "linuxConfiguration": {
          "disablePasswordAuthentication": true,
          "enableVMAgentPlatformUpdates": false,
          "provisionVMAgent": true,
          "ssh": {
            "publicKeys": [
              {
                "keyData": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCvR1+fGFuVMWS2bAY0SgW4E9QzLZ77ETdbCBUVF46eAyL8JWsLynX214hNSK16l4UYZyC3E6jea5qw2rGHPP4eMp7iif50xqd6qGICS428mqc9Gz29J0LFanM7XpHwLnBiJ6hmKvqvHB5tsGKh44MddW0wv+KiiEHIV1ZdSSvBRJ5MMQhqZoUiqlChHourOhaZxvw2dpJhRCvAEKw1s5RoeoLJAdZ6Qr53ERSkJr3BF7uAoNlGx6gatBVkjV+w9CZXN/YN62b1QQiGnk5/BIXNqEIsyxsa84+GbyieRIN/wYjSEV7ASRxSj60qV7RPexvAI+4JGa9UELYMQDrBElgL",
                "path": "/home/azureuser/.ssh/authorized_keys"
              }
            ]
          }
        },
        "requireGuestProvisionSignal": true,
        "secrets": []
      },
      "storageProfile": {
        "imageReference": {
          "offer": "UbuntuServer",
          "publisher": "Canonical",
          "sku": "18.04-LTS",
          "version": "latest"
        },
        "osDisk": {
          "caching": "ReadWrite",
          "createOption": "FromImage",
          "diskSizeGB": 30,
          "managedDisk": {
            "storageAccountType": "Premium_LRS"
          },
          "osType": "Linux"
        }
      }
    }
  }
}
```

## Allow traffic to application
When the scale set was created, an Azure load balancer was automatically deployed. The load balancer distributes traffic to the VM instances in the scale set. To allow traffic to reach the sample web application, create a load balancer rule with [az network lb rule create](/cli/azure/network/lb/rule). The following example creates a rule named *myLoadBalancerRuleWeb*:

```azurecli-interactive
az network lb rule create \
  --resource-group $RESOURCE_GROUP_NAME \
  --name $LOAD_BALANCER_RULE_NAME \
  --lb-name $LOAD_BALANCER_NAME \
  --backend-pool-name $BACKEND_POOL_NAME \
  --backend-port 80 \
  --frontend-ip-name $FRONT_END_IP_NAME \
  --frontend-port 80 \
  --protocol tcp
```

## Test your scale set
To see your scale set in action, access the sample web application in a web browser. Obtain the public IP address of your load balancer with [az network public-ip show](/cli/azure/network/public-ip). The following example obtains the IP address for *myScaleSetLBPublicIP* created as part of the scale set:

```azurecli-interactive
az network public-ip show \
  --resource-group $RESOURCE_GROUP_NAME \
  --name $SCALE_SET_PUBLIC_IP \
  --query '[ipAddress]' \
  --output tsv
```

Enter the public IP address of the load balancer in to a web browser. The load balancer distributes traffic to one of your VM instances, as shown in the following example:

![Default web page in NGINX](media/virtual-machine-scale-sets-create-cli/running-nginx-site.png)

Or run the following command in a local shell to validate the scale set is set up properly 

```bash
 curl $(az network public-ip show --resource-group $RESOURCE_GROUP_NAME --name $SCALE_SET_PUBLIC_IP --query '[ipAddress]' --output tsv)
```

<!--expected_similarity=0.6-->
```HTML
Hello World from host myscabd00000000 !
```

## Clean up resources
When no longer needed, you can use [az group delete](/cli/azure/group) to remove the resource group, scale set, and all related resources as follows. The `--no-wait` parameter returns control to the prompt without waiting for the operation to complete. The `--yes` parameter confirms that you wish to delete the resources without an additional prompt to do so.

```azurecli-interactive
az group delete --name $RESOURCE_GROUP_NAME --yes --no-wait
```


## Next steps
In this quickstart, you created a basic scale set and used the Custom Script Extension to install a basic NGINX web server on the VM instances. To learn more, continue to the tutorial for how to create and manage Azure Virtual Machine Scale Sets.

> [!div class="nextstepaction"]
> [Create and manage Azure Virtual Machine Scale Sets](tutorial-create-and-manage-cli.md)