---
title: Deploy container group to Azure virtual network
description: Learn how to deploy a container group to a new or existing Azure virtual network via the Azure CLI.
ms.topic: how-to
ms.author: tomcassidy
author: tomvcassidy
ms.service: container-instances
ms.date: 06/17/2022
ms.custom: innovation-engine, devx-track-azurecli, linux-related-content
ms.permissions: Microsoft.Resources/deployments/read, Microsoft.Network/profiles/delete, Microsoft.ContainerInstance/containerGroups/restart/action, Microsoft.Network/profiles/read, Microsoft.ContainerInstance/containerGroups/write, Microsoft.Network/virtualNetworks/subnets/join/action, Microsoft.ContainerInstance/containerGroups/delete, Microsoft.Resources/resourceGroups/read, Microsoft.Network/virtualNetworks/networkSecurityGroups/read, Microsoft.Network/virtualNetworks/read, Microsoft.Network/virtualNetworks/subnets/read, Microsoft.ContainerInstance/containerGroups/list, Microsoft.Network/virtualnetworks/delete, Microsoft.Network/virtualNetworks/subnets/write, Microsoft.ContainerInstance/containerGroups/read, Microsoft.ContainerInstance/containerGroups/logs/read
---

# Deploy container instances into an Azure virtual network

[Azure Virtual Network](../virtual-network/virtual-networks-overview.md) provides secure, private networking for your Azure and on-premises resources. By deploying container groups into an Azure virtual network, your containers can communicate securely with other resources in the virtual network.

This article shows how to use the [az container create][az-container-create] command in the Azure CLI to deploy container groups to either a new virtual network or an existing virtual network.

> [!IMPORTANT]
> Before deploying container groups in virtual networks, we suggest checking the limitation first. For networking scenarios and limitations, see [Virtual network scenarios and resources for Azure Container Instances](container-instances-virtual-network-concepts.md).

> [!IMPORTANT]
> Container group deployment to a virtual network is generally available for Linux and Windows containers, in most regions where Azure Container Instances is available. For details, see [available-regions][available-regions].

[!INCLUDE [network profile callout](./includes/network-profile/network-profile-callout.md)]

Examples in this article are formatted for the Bash shell. If you prefer another shell such as PowerShell or Command Prompt, adjust the line continuation characters accordingly.

## Define Environment Variables

The First step in this tutorial is to define environment variables.

```bash
export RANDOM_ID="$(openssl rand -hex 3)"
export MY_RESOURCE_GROUP=myResourceGroup$RANDOM_ID
export LOCATION=EastUS
export MY_APP_CONTAINER_NAME=appcontainer$RANDOM_ID
export MY_COMM_CHECKER_NAME=commchecker$RANDOM_ID
export MY_ACI_VNET_NAME=aci-vnet$RANDOM_ID
export MY_IMAGE_1=mcr.microsoft.com/azuredocs/aci-helloworld
export MY_IMAGE_2=alpine:3.5
export MY_VNET=aci-vnet
export MY_VNETADDRESSPREFIX=10.0.0.0/16
export MY_SUBNET=aci-subnet
export MY_SUBNETADDRESSPREFIX=10.0.0.0/24
export MY_QUERY_1=ipAddress.ip
export MY_OUTPUT_1=tsv
export MY_RESTARTPOLICY=never
```

## Login to Azure using the CLI

In order to run commands against Azure using the CLI you need to login. This is done, very simply, though the `az login` command:

## Create a resource group

A resource group is a container for related resources. All resources must be placed in a resource group. We will create one for this tutorial. The following command creates a resource group with the previously defined $MY_RESOURCE_GROUP and $LOCATION parameters.

```bash
az group create --name $MY_RESOURCE_GROUP --location $LOCATION
```

## Deploy to new virtual network

> [!NOTE]
> If you are using subnet IP range /29 to have only 3 IP addresses. we recommend always to go one range above (never below). For example, use subnet IP range /28 so you can have at least 1 or more IP buffer per container group. By doing this, you can avoid containers in stuck, not able to start, restart or even not able to stop states.

To deploy to a new virtual network and have Azure create the network resources for you automatically, specify the following when you execute [az container create][az-container-create]:

* Virtual network name
* Virtual network address prefix in CIDR format
* Subnet name
* Subnet address prefix in CIDR format

The virtual network and subnet address prefixes specify the address spaces for the virtual network and subnet, respectively. These values are represented in Classless Inter-Domain Routing (CIDR) notation, for example `10.0.0.0/16`. For more information about working with subnets, see [Add, change, or delete a virtual network subnet](../virtual-network/virtual-network-manage-subnet.md).

Once you've deployed your first container group with this method, you can deploy to the same subnet by specifying the virtual network and subnet names, or the network profile that Azure automatically creates for you. Because Azure delegates the subnet to Azure Container Instances, you can deploy *only* container groups to the subnet.

### Example

The following [az container create][az-container-create] command specifies settings for a new virtual network and subnet. Provide the name of a resource group that was created in a region where container group deployments in a virtual network are [available](container-instances-region-availability.md). This command deploys the public Microsoft [aci-helloworld][aci-helloworld] container that runs a small Node.js webserver serving a static web page. In the next section, you'll deploy a second container group to the same subnet, and test communication between the two container instances.

```azurecli-interactive
az container create \
  --name $MY_APP_CONTAINER_NAME \
  --resource-group $MY_RESOURCE_GROUP \
  --image $MY_IMAGE_1 \
  --vnet $MY_ACI_VNET_NAME \
  --vnet-address-prefix $MY_VNETADDRESSPREFIX \
  --subnet $MY_SUBNET \
  --subnet-address-prefix $MY_SUBNETADDRESSPREFIX
```

When you deploy to a new virtual network by using this method, the deployment can take a few minutes while the network resources are created. After the initial deployment, additional container group deployments to the same subnet complete more quickly.

## Deploy to existing virtual network

To deploy a container group to an existing virtual network:

1. Create a subnet within your existing virtual network, use an existing subnet in which a container group is already deployed, or use an existing subnet emptied of *all* other resources and configuration.
1. Deploy a container group with [az container create][az-container-create] and specify one of the following:
   * Virtual network name and subnet name
   * Virtual network resource ID and subnet resource ID, which allows using a virtual network from a different resource group
   * Network profile name or ID, which you can obtain using [az network profile list][az-network-profile-list]

### Example

The following example deploys a second container group to the same subnet created previously, and verifies communication between the two container instances.

First, get the IP address of the first container group you deployed, the *appcontainer*:

```azurecli-interactive
az container show --resource-group $MY_RESOURCE_GROUP \
  --name $MY_APP_CONTAINER_NAME \
  --query $MY_QUERY_1 --output $MY_OUTPUT_1
```

The output displays the IP address of the container group in the private subnet. For example:

```output
10.0.0.4
```

Now, set `CONTAINER_GROUP_IP` to the IP you retrieved with the `az container show` command, and execute the following `az container create` command. This second container, *commchecker*, runs an Alpine Linux-based image and executes `wget` against the first container group's private subnet IP address.

```azurecli-interactive
CONTAINER_GROUP_IP=10.0.0.4

az container create \
  --resource-group $MY_RESOURCE_GROUP \
  --name $MY_COMM_CHECKER_NAME \
  --image $MY_IMAGE_2 \
  --command-line "wget $CONTAINER_GROUP_IP" \
  --restart-policy $MY_RESTARTPOLICY \
  --vnet $MY_ACI_VNET_NAME \
  --subnet $MY_SUBNET
```

After this second container deployment has completed, pull its logs so you can see the output of the `wget` command it executed:

```azurecli-interactive
az container logs --resource-group $MY_RESOURCE_GROUP --name $MY_COMM_CHECKER_NAME
```

If the second container communicated successfully with the first, output is similar to:

```output
Connecting to 10.0.0.4 (10.0.0.4:80)
index.html           100% |*******************************|  1663   0:00:00 ETA
```

The log output should show that `wget` was able to connect and download the index file from the first container using its private IP address on the local subnet. Network traffic between the two container groups remained within the virtual network.

## Next steps

* To deploy a new virtual network, subnet, network profile, and container group using a Resource Manager template, see [Create an Azure container group with VNet](https://github.com/Azure/azure-quickstart-templates/tree/master/quickstarts/microsoft.containerinstance/aci-vnet).

* To deploy Azure Container Instances that can pull images from an Azure Container Registry through a private endpoint, see [Deploy to Azure Container Instances from Azure Container Registry using a managed identity](../container-instances/using-azure-container-registry-mi.md).

<!-- IMAGES -->
[aci-vnet-01]: ./media/container-instances-vnet/aci-vnet-01.png

<!-- LINKS - External -->
[aci-helloworld]: https://hub.docker.com/_/microsoft-azuredocs-aci-helloworld

<!-- LINKS - Internal -->
[az-container-create]: /cli/azure/container#az_container_create
[az-container-show]: /cli/azure/container#az_container_show
[az-network-vnet-create]: /cli/azure/network/vnet#az_network_vnet_create
[az-network-profile-list]: /cli/azure/network/profile#az_network_profile_list
[available-regions]: https://azure.microsoft.com/explore/global-infrastructure/products-by-region/?products=container-instances

<details>
<summary><h2>FAQs</h2></summary>

#### Q. What is the command-specific breakdown of permissions needed to implement this doc?

A. _Format: Commands as they appears in the doc | list of unique permissions needed to run each of those commands_

  - ```azurecli-interactive az container create \ --name $MY_APP_CONTAINER_NAME \ --resource-group $MY_RESOURCE_GROUP \ --image $MY_IMAGE_1 \ --vnet $MY_ACI_VNET_NAME \ --vnet-address-prefix $MY_VNETADDRESSPREFIX \ --subnet $MY_SUBNET \ --subnet-address-prefix $MY_SUBNETADDRESSPREFIX ```

      - Microsoft.ContainerInstance/containerGroups/restart/action
      - Microsoft.ContainerInstance/containerGroups/write
      - Microsoft.ContainerInstance/containerGroups/delete
      - Microsoft.Network/virtualNetworks/read
      - Microsoft.Network/virtualNetworks/subnets/read
      - Microsoft.ContainerInstance/containerGroups/list
      - Microsoft.ContainerInstance/containerGroups/read
  - ```azurecli-interactive az container show --resource-group $MY_RESOURCE_GROUP \ --name $MY_APP_CONTAINER_NAME \ --query $MY_QUERY_1 --output $MY_OUTPUT_1 ```

      - Microsoft.Resources/resourceGroups/read
      - Microsoft.Resources/deployments/read
      - Microsoft.ContainerInstance/containerGroups/read
  - ```azurecli-interactive CONTAINER_GROUP_IP=<container-group-IP-address> az container create \ --resource-group $MY_RESOURCE_GROUP \ --name $MY_COMM_CHECKER_NAME \ --image $MY_IMAGE_2 \ --command-line "wget $CONTAINER_GROUP_IP" \ --restart-policy $MY_RESTARTPOLICY \ --vnet $MY_ACI_VNET_NAME \ --subnet $MY_SUBNET ```

      - Microsoft.Network/virtualNetworks/subnets/join/action
      - Microsoft.ContainerInstance/containerGroups/write
      - Microsoft.Network/virtualNetworks/networkSecurityGroups/read
      - Microsoft.Network/virtualNetworks/read
      - Microsoft.Network/virtualNetworks/subnets/read
      - Microsoft.Network/virtualNetworks/subnets/write
  - ```azurecli-interactive az container logs --resource-group $MY_RESOURCE_GROUP --name $MY_COMM_CHECKER_NAME ```

      - Microsoft.ContainerInstance/containerGroups/read
      - Microsoft.ContainerInstance/containerGroups/logs/read
  - ```azurecli-interactive az container create --resource-group $MY_RESOURCE_GROUP \ --file $MY_FILE ```

      - Microsoft.ContainerInstance/containerGroups/delete
      - Microsoft.ContainerInstance/containerGroups/list
      - Microsoft.ContainerInstance/containerGroups/read
      - Microsoft.ContainerInstance/containerGroups/write
  - ```azurecli-interactive az container delete --resource-group $MY_RESOURCE_GROUP --name $MY_APP_CONTAINER_NAME -y az container delete --resource-group $MY_RESOURCE_GROUP --name $MY_COMM_CHECKER_NAME -y az container delete --resource-group $MY_RESOURCE_GROUP --name $MY_NAME_3 -y ```

      - Microsoft.ContainerInstance/containerGroups/delete
      - Microsoft.ContainerInstance/containerGroups/write
  - ```azurecli-interactive # Replace <my-resource-group> with the name of your resource group # Assumes one virtual network in resource group RES_GROUP=<my-resource-group> # Get network profile ID # Assumes one profile in virtual network NETWORK_PROFILE_ID=$(az network profile list --resource-group $RES_GROUP --query [0].id --output $MY_OUTPUT_1) # Delete the network profile az network profile delete --id $NETWORK_PROFILE_ID -y # Delete virtual network az network vnet delete --resource-group $RES_GROUP --name $MY_ACI_VNET_NAME ```

      - Microsoft.Network/profiles/delete
      - Microsoft.Network/virtualnetworks/delete
      - Microsoft.Network/profiles/read

#### Q. What are the prerequisites for deploying containers to an Azure virtual network?

A. Before deploying container groups to an Azure virtual network, it is important to check the limitations and prerequisites. You can find detailed information about the networking scenarios and limitations for Azure Container Instances in the [Virtual network scenarios and resources for Azure Container Instances](https://docs.microsoft.com/azure/container-instances/container-instances-virtual-network-concepts) documentation.

#### Q. How do I deploy container instances to a new virtual network?

A. To deploy container instances to a new virtual network, you need to use the `az container create` command with appropriate parameters. The parameters include the virtual network name, virtual network address prefix, subnet name, and subnet address prefix. You can find an example command and detailed instructions in the [Deploy to new virtual network](https://docs.microsoft.com/azure/container-instances/container-instances-vnet-deploy#deploy-to-new-virtual-network) section of the documentation.

#### Q. How do I deploy container instances to an existing virtual network?

A. To deploy container instances to an existing virtual network, you need to create a subnet within the virtual network and then use the `az container create` command with the appropriate parameters. The parameters can include the virtual network name and subnet name, virtual network resource ID and subnet resource ID, or network profile name or ID. You can find an example command and detailed instructions in the [Deploy to existing virtual network](https://docs.microsoft.com/azure/container-instances/container-instances-vnet-deploy#deploy-to-existing-virtual-network) section of the documentation.

#### Q. How can I verify communication between container instances deployed to the same virtual network?

A. To verify communication between container instances deployed to the same virtual network, you can use the `az container show` and `az container logs` commands. The `az container show` command retrieves the IP address of a container group, while the `az container logs` command retrieves the logs of a container group. You can find detailed examples and instructions in the [Deploy to existing virtual network](https://docs.microsoft.com/azure/container-instances/container-instances-vnet-deploy#deploy-to-existing-virtual-network) section of the documentation.

#### Q. What is the process to delete container instances and network resources?

A. To delete container instances, you can use the `az container delete` command with the appropriate parameters, such as resource group name and container instance name. To delete network resources, there is a script provided in the documentation. The script assumes that you have created a virtual network and subnet using the example commands provided earlier in the documentation. The script deletes the virtual network and subnet. Before running the script, ensure that you no longer need any of the resources in the virtual network, as the deletion of resources is irreversible. You can find the script and more details in the [Clean up resources](https://docs.microsoft.com/azure/container-instances/container-instances-vnet-deploy#clean-up-resources) section of the documentation.

</details>
