# Create a Virtual Machine Scale Set with Application Gateway with Linux image

## Define Environment Variables

The First step in this tutorial is to define environment variables.

```bash

export RANDOM_ID="$(openssl rand -hex 3)"
export MY_RESOURCE_GROUP_NAME="myVMSSResourceGroup$RANDOM_ID"
export REGION=EastUS
export MY_VMSS_NAME="myVMSS$RANDOM_ID"
export MY_USERNAME=azureuser
export MY_VM_IMAGE="Canonical:0001-com-ubuntu-minimal-jammy:minimal-22_04-lts-gen2:latest"


export MY_VNET_NAME="myVNet$RANDOM_ID"
export NETWORK_PREFIX="$(($RANDOM % 254 + 1))"
export MY_VNET_PREFIX="10.$NETWORK_PREFIX.0.0/16"
export MY_VM_SN_NAME="myVMSN$RANDOM_ID"
export MY_VM_SN_PREFIX="10.$NETWORK_PREFIX.0.0/24"
export MY_APPGWSN_NAME="myAPPGWSN$RANDOM_ID"
export MY_APPGWSN_PREFIX="10.$NETWORK_PREFIX.1.0/24"
export MY_APPGW_NAME="myAPPGW$RANDOM_ID"
export MY_APPGW_BACKENDPOOL_NAME="myAPPGWBackendpool$RANDOM_ID"

```
# Login to Azure using the CLI

In order to run commands against Azure using the CLI you need to login. This is done, very simply, though the `az login` command:

# Create a resource group

A resource group is a container for related resources. All resources must be placed in a resource group. We will create one for this tutorial. The following command creates a resource group with the previously defined $MY_RESOURCE_GROUP_NAME and $REGION parameters.

```bash
az group create --name $MY_RESOURCE_GROUP_NAME --location $REGION
```

Results:

<!-- expected_similarity=0.3 -->
```json   
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myVMSSResourceGroupxxxxxx",
  "location": "eastus",
  "managedBy": null,
  "name": "myVMSSResourceGroupxxxxxx",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```

## Create Network Resources 

You need to create network resources before you proceed the VMSS steps. In this step you're going to create a VNET, 2 subnets 1 for Application Gateway and 1 for VMs. You also need to have a public IP to attach your Application Gateway to be able to reach your web application from internet. 


# Create network resources 

```bash
az network vnet create  --name $MY_VNET_NAME  --resource-group $MY_RESOURCE_GROUP_NAME --location $REGION  --address-prefix $MY_VNET_PREFIX  --subnet-name $MY_VM_SN_NAME --subnet-prefix $MY_VM_SN_PREFIX 
```

Results:

<!-- expected_similarity=0.3 -->
```json   
{
  "newVNet": {
    "addressSpace": {
      "addressPrefixes": [
        "10.X.0.0/16"
      ]
    },
    "enableDdosProtection": false,
    "etag": "W/\"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\"",
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myVMSSResourceGroupxxxxxx/providers/Microsoft.Network/virtualNetworks/myVNetxxxxxx",
    "location": "eastus",
    "name": "myVNetxxxxxx",
    "provisioningState": "Succeeded",
    "resourceGroup": "myVMSSResourceGroupxxxxxx",
    "resourceGuid": "f00034be-612e-4462-a711-93d0bb263e46",
    "subnets": [
      {
        "addressPrefix": "10.X.0.0/24",
        "delegations": [],
        "etag": "W/\"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx\"",
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myVMSSResourceGroupxxxxxx/providers/Microsoft.Network/virtualNetworks/myVNetxxxxxx/subnets/myVMSNxxxxxx", 
        "name": "myVMSNxxxxxx",
        "privateEndpointNetworkPolicies": "Disabled",
        "privateLinkServiceNetworkPolicies": "Enabled",
        "provisioningState": "Succeeded",
        "resourceGroup": "myVMSSResourceGroupxxxxxx",
        "type": "Microsoft.Network/virtualNetworks/subnets"
      }
    ],
    "type": "Microsoft.Network/virtualNetworks",
    "virtualNetworkPeerings": []
  }
}
```



### Enable Azure AD login for a Linux Virtual Machine in Azure

The following example has deploys a Linux VM and then installs the extension to enable Azure AD login for a Linux VM. VM extensions are small applications that provide post-deployment configuration and automation tasks on Azure virtual machines.

```bash
az vm extension set \
    --publisher Microsoft.Azure.ActiveDirectory \
    --name AADSSHLoginForLinux \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --vm-name $MY_VM_NAME
```

# Store IP Address of VM in order to SSH
run the following command to get the IP Address of the VM and store it as an environment variable

```bash
export IP_ADDRESS=$(az vm show --show-details --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_VM_NAME --query publicIps --output tsv)
```

# SSH Into VM

<!--## Export the SSH configuration for use with SSH clients that support OpenSSH & SSH into the VM.
Login to Azure Linux VMs with Azure AD supports exporting the OpenSSH certificate and configuration. That means you can use any SSH clients that support OpenSSH-based certificates to sign in through Azure AD. The following example exports the configuration for all IP addresses assigned to the VM:-->

<!--
```bash
yes | az ssh config --file ~/.ssh/config --name $MY_VM_NAME --resource-group $MY_RESOURCE_GROUP_NAME
```
-->

You can now SSH into the VM by running the output of the following command in your ssh client of choice

```bash
ssh -o StrictHostKeyChecking=no $MY_USERNAME@$IP_ADDRESS
```

# Next Steps

* [VM Documentation](https://learn.microsoft.com/en-us/azure/virtual-machines/)
* [Use Cloud-Init to initialize a Linux VM on first boot](https://learn.microsoft.com/en-us/azure/virtual-machines/linux/tutorial-automate-vm-deployment)
* [Create custom VM images](https://learn.microsoft.com/en-us/azure/virtual-machines/linux/tutorial-custom-images)
* [Load Balance VMs](https://learn.microsoft.com/en-us/azure/load-balancer/quickstart-load-balancer-standard-public-cli)