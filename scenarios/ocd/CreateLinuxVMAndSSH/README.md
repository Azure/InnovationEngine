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

```json expected-similarity=0.7
{
  "id": "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/myResourceGroup",
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

## Create the Virtual Machine

To create a VM in this resource group we need to run a simple command, here we have provided the `--generate-ssh-keys` flag, this will cause the CLI to look for an avialable ssh key in `~/.ssh`, if one is found it will be used, otherwise one will be generated and stored in `~/.ssh`. We also provide the `--public-ip-sku Standard` flag to ensure that the machine is accessible via a public IP. Finally, we are deploying an `UbuntuLTS` image. 

All other values are configured using environment variables.

```bash
az vm create --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_VM_NAME --image $MY_VM_IMAGE --assign-identity --admin-username $MY_USERNAME --generate-ssh-keys --public-ip-sku Standard
```

Results:

```json expected-similarity=0.7
{
  "fqdns": "",
  "id": "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM",
  "location": "eastus",
  "macAddress": "00-0D-3A-10-4F-70",
  "powerState": "VM running",
  "privateIpAddress": "10.0.0.4",
  "publicIpAddress": "52.147.208.85",
  "resourceGroup": "myResourceGroup",
  "zones": ""
}
```

## Add VM AAD Extension

In order to use Azure AD Login for a Linux VM an extension needs to be installed on the VM. VM extensions are small applications that provide post-deployment configuration and automation tasks on Azure virtual machines

The following installs the extension to enable Azure AD Login on the recently deployed VM

```bash
az vm extension set --publisher Microsoft.Azure.ActiveDirectory --name AADSSHLoginForLinux --resource-group $MY_RESOURCE_GROUP_NAME --vm-name $MY_VM_NAME
```
## Configure role assignments for the VM

Now that you've created the VM, you need to configure an Azure RBAC policy to determine who can log in to the VM. Two Azure roles are used to authorize VM login:

Virtual Machine Administrator Login: Users who have this role assigned can log in to an Azure virtual machine with administrator privileges.
Virtual Machine User Login: Users who have this role assigned can log in to an Azure virtual machine with regular user privileges.
To allow a user to log in to a VM over SSH, you must assign the Virtual Machine Administrator Login or Virtual Machine User Login role on the resource group that contains the VM and its associated virtual network, network interface, public IP address, or load balancer resources.

An Azure user who has the Owner or Contributor role assigned for a VM doesn't automatically have privileges to Azure AD login to the VM over SSH. There's an intentional (and audited) separation between the set of people who control virtual machines and the set of people who can access virtual machines.

The following example uses az role assignment create to assign the Virtual Machine Administrator Login role to the VM for your current Azure user. You obtain the username of your current Azure account by using az account show, and you set the scope to the VM created in a previous step by using az vm show.

You can also assign the scope at a resource group or subscription level. Normal Azure RBAC inheritance permissions apply

```bash
USERNAME=$(az account show --query user.name --output tsv)
RESOURCE_GROUP_ID=$(az group show --resource-group $MY_RESOURCE_GROUP_NAME --query id -o tsv)
az role assignment create --role "Virtual Machine Administrator Login" --assignee $USERNAME --scope $RESOURCE_GROUP_ID
```

# SSH Into VM

You can now SSH into the VM by running the output of the following command in your ssh client of choice

```bash
echo az ssh vm --name $MY_VM_NAME --resource-group $MY_RESOURCE_GROUP_NAME
```

# Next Steps

* [VM Documentation](https://learn.microsoft.com/en-us/azure/virtual-machines/)
* [Use Cloud-Init to initialize a Linux VM on first boot](https://learn.microsoft.com/en-us/azure/virtual-machines/linux/tutorial-automate-vm-deployment)
* [Create custom VM images](https://learn.microsoft.com/en-us/azure/virtual-machines/linux/tutorial-custom-images)
* [Load Balance VMs](https://learn.microsoft.com/en-us/azure/load-balancer/quickstart-load-balancer-standard-public-cli)