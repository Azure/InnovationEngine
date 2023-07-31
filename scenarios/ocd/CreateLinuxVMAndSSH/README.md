# Create a Linux VM and SSH On Azure

## Define Environment Variables

The First step in this tutorial is to define environment variables 

```bash
export MY_RESOURCE_GROUP_NAME=myResourceGroup
export MY_LOCATION=EastUS
export MY_VM_NAME=myVM
export MY_USERNAME=azureuser
export MY_VM_IMAGE=UbuntuLTS
```

# Login to Azure using the CLI

In order to run commands against Azure using the CLI you need to login. This is done, very simply, though the `az login` command:

# Create a resource group

A resource group is a container for related resources. All resources must be placed in a resource group. We will create one for this tutorial. The following command creates a resource group with the previously defined $MY_RESOURCE_GROUP_NAME and $MY_LOCATION parameters.

```bash
az group create --name $MY_RESOURCE_GROUP_NAME --location $MY_LOCATION
```

Results:

```expected_similarity=0.3
  "id": "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/myResourceGroup",
  "location": "eastus",
  "managedBy": null,
  "name": "myResourceGroup",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
```

## Create the Virtual Machine

To create a VM in this resource group we need to run a simple command, here we have provided the `--generate-ssh-keys` flag, this will cause the CLI to look for an avialable ssh key in `~/.ssh`, if one is found it will be used, otherwise one will be generated and stored in `~/.ssh`. We also provide the `--public-ip-sku Standard` flag to ensure that the machine is accessible via a public IP. Finally, we are deploying an `UbuntuLTS` image. 

All other values are configured using environment variables.

```bash
az vm create --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_VM_NAME --image $MY_VM_IMAGE --admin-username $MY_USERNAME --generate-ssh-keys --public-ip-sku Standard
```

Results:

```expected_similarity=0.3
  "fqdns": "",
  "id": "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM",
  "location": "eastus",
  "macAddress": "00-0D-3A-10-4F-70",
  "powerState": "VM running",
  "privateIpAddress": "10.0.0.4",
  "publicIpAddress": "52.147.208.85",
  "resourceGroup": "myResourceGroup",
  "zones": ""
```

# Store IP Address of VM in order to SSH
run the following command to get the IP Address of the VM and store it as an environment variable

```bash
export IP_ADDRESS=$(az vm show --show-details --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_VM_NAME --query publicIps --output tsv)
```
# SSH Into VM

You can now SSH into the VM by running the output of the following command in your ssh client of choice

```bash
ssh -o StrictHostKeyChecking=no $MY_USERNAME@$IP_ADDRESS
```

# Next Steps

* [VM Documentation](https://learn.microsoft.com/en-us/azure/virtual-machines/)
* [Use Cloud-Init to initialize a Linux VM on first boot](https://learn.microsoft.com/en-us/azure/virtual-machines/linux/tutorial-automate-vm-deployment)
* [Create custom VM images](https://learn.microsoft.com/en-us/azure/virtual-machines/linux/tutorial-custom-images)
* [Load Balance VMs](https://learn.microsoft.com/en-us/azure/load-balancer/quickstart-load-balancer-standard-public-cli)