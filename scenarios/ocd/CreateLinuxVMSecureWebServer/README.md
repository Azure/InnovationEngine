# Intro to Create a NGINX Webserver Secured via HTTPS
Welcome to this tutorial where we will create a VM. This tutorial assumes you are logged into Azure CLI already and have selected a subscription to use with the CLI. If you have not done this already. Press b and hit ctl c to exit the program. Following that you can enter 

'az login' followed by 'az account list --output table' and 'az account set --subscription "name of subscription to use"'


If you need to install Azure CLI run the following command - curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

testing the changes - by bc

Assuming the pre requisites are met press space bar to proceed

## Create a resource Group
The first thing we need to do is create a resource group. You can do this by running the following command

'az group create --name $RESOURCE_GROUP_NAME --location $RESOURCE_LOCATION'

```bash
az group create --name $RESOURCE_GROUP_NAME --location $RESOURCE_LOCATION
```

Results:
```
{
  "id": "/subscriptions/8c487e6a-8bbb-42bb-81e6-3c122d1bb1c7/resourceGroups/$RESOURCE_GROUP_NAME",
  "location": "eastus",
  "managedBy": null,
  "name": "$RESOURCE_GROUP_NAME",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}

```

## Create a Virtual Machine (VM)
You can do this by running the following command:

'az vm create --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --image $VM_IMAGE --admin-username $VM_ADMIN_USERNAME --generate-ssh-keys'

```bash
az vm create --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --image $VM_IMAGE --admin-username $VM_ADMIN_USERNAME --generate-ssh-keys
```

Results:

```
{
  "fqdns": "",
  "id": "/subscriptions/<guid>/resourceGroups/$RESOURCE_GROUP_NAME2/providers/Microsoft.Compute/virtualMachines/$VM_NAME",
  "location": "eastus",
  "macAddress": "00-0D-3A-23-9A-49",
  "powerState": "VM running",
  "privateIpAddress": "10.0.0.4",
  "publicIpAddress": "52.174.34.95",
  "resourceGroup": "$RESOURCE_GROUP_NAME"
}
```

Congrats you created a VM! Next we will open port 80 and install NGINX. 

## Store IP Address as environment Variable 
The following command will store the IP Address as a environment variable that we can access later to do SSH

'export IP_ADDRESS=$(az vm show --show-details --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --query publicIps --output tsv)'

```bash
export IP_ADDRESS=$(az vm show --show-details --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME --query publicIps --output tsv)
```

## Validate IP_ADDRESS
Let's make sure the IP Address is correctly stored

```bash
echo $IP_ADDRESS
```

# Open Port 80 to allow web traffic 
Open port 80 with the following command:

'az vm open-port --port 80,443 --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME'

```bash
az vm open-port --port 80,443 --resource-group $RESOURCE_GROUP_NAME --name $VM_NAME
```

## Validate SSH Connection
To validate you are connected to your VM you can run the following command: 

'ssh -o StrictHostKeyChecking=no $VM_ADMIN_USERNAME@$IP_ADDRESS hostname'
Note - For the following commands we place ssh before hand as we must connect to the VM each time to run commands

```bash
ssh -o StrictHostKeyChecking=no $VM_ADMIN_USERNAME@$IP_ADDRESS hostname
```

## Ensure VM is up to date
Ensure the VM is up to date by running the following command: 

'sudo apt update'
Note - This may take ~30 seconds to complete

```bash
ssh $VM_ADMIN_USERNAME@$IP_ADDRESS sudo apt update
```

## Install NGINX
Run the following command to install the NGINX webserver

'sudo apt install nginx'
This may take a few minutes...

```bash
ssh $VM_ADMIN_USERNAME@$IP_ADDRESS sudo apt --yes --force-yes install nginx
```

## View Your webserver running

```bash
echo $IP_ADDRESS
```

Congratulations you have now created a Virtual Machine and installed a webserver!

Press 1 to end the tutorial and 2 to secure your webserver via https 

1. Quit the tutorial
2. Secure your webserver via https and add a custom domain

## Select unique custom domain Name 

When prompted to enter a value for CUSTOM_DOMAIN_NAME enter a custom domain for your webserver Note - This must be unique on Azure

```bash
echo $CUSTOM_DOMAIN_NAME
```

install Az CLI extension for Front Door in order to add HTTPS

```bash
az extension add --name front-door
```

## Setting Up HTTPS Terminated EndPoint

The following command will set up a custom domain secured via https. This may take a few minutes 
```bash
az network front-door create --backend-address $IP_ADDRESS --name $CUSTOM_DOMAIN_NAME --resource-group $RESOURCE_GROUP_NAME --accepted-protocols Http Https --forwarding-protocol HttpOnly --protocol Http 
```

## See your webserver running HTTPS

Run the following command to see the url of your webserver.
NOTE - It may take ~5 minutes for backends to update appropriately and for your site to be secured via https.

```bash
az network front-door show --name $CUSTOM_DOMAIN_NAME --resource-group $RESOURCE_GROUP_NAME --query frontendEndpoints[*].hostName --output tsv
```

## Conclusion

You have completed the tutorial! View your resources on portal.azure.com 

# Next Steps

* [VM Documentation](https://learn.microsoft.com/en-us/azure/virtual-machines/)
* [Create Vm Scale Set](https://learn.microsoft.com/en-us/azure/virtual-machine-scale-sets/flexible-virtual-machine-scale-sets-cli)
* [Load Balance VMs](https://learn.microsoft.com/en-us/azure/load-balancer/quickstart-load-balancer-standard-public-cli)
* [Baclup VMs](https://learn.microsoft.com/en-us/azure/virtual-machines/backup-recovery)
