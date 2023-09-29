---
title: Subscribe to Azure Kubernetes Service events with Azure Event Grid
description: Use Azure Event Grid to subscribe to Azure Kubernetes Service events
services: container-service
author: zr-msft
ms.topic: article
ms.date: 07/12/2021
ms.author: zarhoads
---

# Quickstart: Subscribe to Azure Kubernetes Service (AKS) events with Azure Event Grid

Azure Event Grid is a fully managed event routing service that provides uniform event consumption using a publish-subscribe model.

In this quickstart, you'll create an AKS cluster and subscribe to AKS events.

## Prerequisites

* An Azure subscription. If you don't have an Azure subscription, you can create a [free account](https://azure.microsoft.com/free).
* [Azure CLI][azure-cli-install] or [Azure PowerShell][azure-powershell-install] installed.

## Create an AKS cluster

### [Azure CLI](#tab/azure-cli)

## Define Environment Variables

This document uses environment variables for all parameters to facilitate reuse. The default values provided here should work in most test environments. For production work you will obviously need to modify these values.

```azurecli-interactive
export RESOURCE_GROUP_NAME=aksQuickstartResourceGroup
export RESOURCE_LOCATION=eastus
export AKS_CLUSTER_NAME=aksQuickstartCluster
export NAMESPACE_NAME="aksQuickstartNamespace$(printf "%08d" $((RANDOM%100000000)))"
export EVENT_GRID_HUB_NAME=aksQuickstartEventGridHub
export EVENT_GRID_SUBSCRIPTION_NAME=aksQuickstartEventGridSubscription
```

## Create an AKS Cluster

Create an AKS cluster using the [az aks create][az-aks-create] command. The following example creates a resource group and a cluster with one node. They will be named according to the environment variables set above:

```azurecli-interactive
az group create --name $RESOURCE_GROUP_NAME --location $RESOURCE_LOCATION
```

<!-- expected_similarity=0.2 -->
```output
{
  "id": "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/aksQuickstartResourceGroup",
  "location": "eastus",
  "managedBy": null,
  "name": "aksQuickstartResourceGroup",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```

Now we can create an AKS cluster within that resource group.

```azurecli-interactive
az aks create --resource-group $RESOURCE_GROUP_NAME --name $AKS_CLUSTER_NAME --location $RESOURCE_LOCATION  --node-count 1 --generate-ssh-keys
```

This will take a little while to run, when it completes you should see an output that looks something like this:

<!-- expected_similarity=0.2 -->
```output
{
  "aadProfile": null,
  "addonProfiles": null,
  "agentPoolProfiles": [
    {
      "availabilityZones": null,
      "count": 1,
      "creationData": null,
      "currentOrchestratorVersion": "1.24.9",
      "enableAutoScaling": false,
      "enableEncryptionAtHost": false,
      "enableFips": false,
      "enableNodePublicIp": false,
      "enableUltraSsd": false,
      "gpuInstanceProfile": null,
      "hostGroupId": null,
      "kubeletConfig": null,
      "kubeletDiskType": "OS",
      "linuxOsConfig": null,
      "maxCount": null,
      "maxPods": 110,
      "minCount": null,
      "mode": "System",
      "name": "nodepool1",
      "nodeImageVersion": "AKSUbuntu-1804gen2containerd-2023.01.20",
      "nodeLabels": null,
      "nodePublicIpPrefixId": null,
      "nodeTaints": null,
      "orchestratorVersion": "1.24.9",
      "osDiskSizeGb": 128,
      "osDiskType": "Managed",
      "osSku": "Ubuntu",
      "osType": "Linux",
      "podSubnetId": null,
      "powerState": {
        "code": "Running"
      },
      "provisioningState": "Succeeded",
      "proximityPlacementGroupId": null,
      "scaleDownMode": null,
      "scaleSetEvictionPolicy": null,
      "scaleSetPriority": null,
      "spotMaxPrice": null,
      "tags": null,
      "type": "VirtualMachineScaleSets",
      "upgradeSettings": {
        "maxSurge": null
      },
      "vmSize": "Standard_DS2_v2",
      "vnetSubnetId": null,
      "workloadRuntime": null
    }
  ],
  "apiServerAccessProfile": null,
  "autoScalerProfile": null,
  "autoUpgradeProfile": null,
  "azurePortalFqdn": "aksquickst-aksquickstartres-325e7c-784c55cf.portal.hcp.eastus.azmk8s.io",
  "currentKubernetesVersion": "1.24.9",
  "disableLocalAccounts": false,
  "diskEncryptionSetId": null,
  "dnsPrefix": "aksQuickst-aksQuickstartRes-325e7c",
  "enablePodSecurityPolicy": null,
  "enableRbac": true,
  "extendedLocation": null,
  "fqdn": "aksquickst-aksquickstartres-325e7c-784c55cf.hcp.eastus.azmk8s.io",
  "fqdnSubdomain": null,
  "httpProxyConfig": null,
  "id": "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourcegroups/aksQuickstartResourceGroup/providers/Microsoft.ContainerService/managedClusters/aksQuickstartCluster",
  "identity": {
    "principalId": "REDACTED",
    "tenantId": "REDACTED",
    "type": "SystemAssigned",
    "userAssignedIdentities": null
  },
  "identityProfile": {
    "kubeletidentity": {
      "clientId": "REDACTED",
      "objectId": "REDACTED",
      "resourceId": "/subscriptions/REDACTED/resourcegroups/MC_aksQuickstartResourceGroup_aksQuickstartCluster_eastus/providers/Microsoft.ManagedIdentity/userAssignedIdentities/aksQuickstartCluster-agentpool"
    }
  },
  "kubernetesVersion": "1.24.9",
  "linuxProfile": {
    "adminUsername": "azureuser",
    "ssh": {
      "publicKeys": [
        {
          "keyData": "ssh-rsa REDACTED"
        }
      ]
    }
  },
  "location": "eastus",
  "maxAgentPools": 100,
  "name": "aksQuickstartCluster",
  "networkProfile": {
    "dnsServiceIp": "10.0.0.10",
    "dockerBridgeCidr": "172.17.0.1/16",
    "ipFamilies": [
      "IPv4"
    ],
    "loadBalancerProfile": {
      "allocatedOutboundPorts": null,
      "effectiveOutboundIPs": [
        {
          "id": "/subscriptions/REDACTED/resourceGroups/MC_aksQuickstartResourceGroup_aksQuickstartCluster_eastus/providers/Microsoft.Network/publicIPAddresses/e19ddc6c-0842-45d5-814d-702cc95945ce",
          "resourceGroup": "MC_aksQuickstartResourceGroup_aksQuickstartCluster_eastus"
        }
      ],
      "enableMultipleStandardLoadBalancers": null,
      "idleTimeoutInMinutes": null,
      "managedOutboundIPs": {
        "count": 1,
        "countIpv6": null
      },
      "outboundIPs": null,
      "outboundIpPrefixes": null
    },
    "loadBalancerSku": "Standard",
    "natGatewayProfile": null,
    "networkMode": null,
    "networkPlugin": "kubenet",
    "networkPolicy": null,
    "outboundType": "loadBalancer",
    "podCidr": "10.244.0.0/16",
    "podCidrs": [
      "10.244.0.0/16"
    ],
    "serviceCidr": "10.0.0.0/16",
    "serviceCidrs": [
      "10.0.0.0/16"
    ]
  },
  "nodeResourceGroup": "MC_aksQuickstartResourceGroup_aksQuickstartCluster_eastus",
  "oidcIssuerProfile": {
    "enabled": false,
    "issuerUrl": null
  },
  "podIdentityProfile": null,
  "powerState": {
    "code": "Running"
  },
  "privateFqdn": null,
  "privateLinkResources": null,
  "provisioningState": "Succeeded",
  "publicNetworkAccess": null,
  "resourceGroup": "aksQuickstartResourceGroup",
  "securityProfile": {
    "azureKeyVaultKms": null,
    "defender": null
  },
  "servicePrincipalProfile": {
    "clientId": "msi",
    "secret": null
  },
  "sku": {
    "name": "Basic",
    "tier": "Free"
  },
  "storageProfile": {
    "blobCsiDriver": null,
    "diskCsiDriver": {
      "enabled": true
    },
    "fileCsiDriver": {
      "enabled": true
    },
    "snapshotController": {
      "enabled": true
    }
  },
  "systemData": null,
  "tags": null,
  "type": "Microsoft.ContainerService/ManagedClusters",
  "windowsProfile": null
}
```

### [Azure PowerShell](#tab/azure-powershell)

Create an AKS cluster using the [New-AzAksCluster][new-azakscluster] command. The following example creates a resource group *MyResourceGroup* and a cluster named *MyAKS* with one node in the *MyResourceGroup* resource group:

```azurepowershell-interactive
New-AzResourceGroup -Name MyResourceGroup -Location eastus
New-AzAksCluster -ResourceGroupName MyResourceGroup -Name MyAKS -Location eastus -NodeCount 1 -GenerateSshKey
```

---

## Subscribe to AKS events

### [Azure CLI](#tab/azure-cli)

Create a namespace and event hub using [az eventhubs namespace create][az-eventhubs-namespace-create] and [az eventhubs eventhub create][az-eventhubs-eventhub-create]. The following example creates a namespace *MyNamespace* and an event hub *MyEventGridHub* in *MyNamespace*, both in the *MyResourceGroup* resource group.

```azurecli-interactive
az eventhubs namespace create --location $RESOURCE_LOCATION --name $NAMESPACE_NAME --resource-group $RESOURCE_GROUP_NAME
```

<!--expected_similarity=0.7-->
```output
{
  "alternateName": null,
  "clusterArmId": null,
  "createdAt": "2023-02-11T00:27:48.977000+00:00",
  "disableLocalAuth": false,
  "encryption": null,
  "id": "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/aksQuickstartResourceGroup/providers/Microsoft.EventHub/namespaces/aksQuickstartNamespace00021677",
  "identity": null,
  "isAutoInflateEnabled": false,
  "kafkaEnabled": true,
  "location": "East US",
  "maximumThroughputUnits": 0,
  "metricId": "325e7c34-99fb-4190-aa87-1df746c67705:aksquickstartnamespace00021677",
  "minimumTlsVersion": "1.2",
  "name": "aksQuickstartNamespace00021677",
  "privateEndpointConnections": null,
  "provisioningState": "Succeeded",
  "publicNetworkAccess": "Enabled",
  "resourceGroup": "aksQuickstartResourceGroup",
  "serviceBusEndpoint": "https://aksQuickstartNamespace00021677.servicebus.windows.net:443/",
  "sku": {
    "capacity": 1,
    "name": "Standard",
    "tier": "Standard"
  },
  "status": "Active",
  "systemData": null,
  "tags": {},
  "type": "Microsoft.EventHub/Namespaces",
  "updatedAt": "2023-02-11T00:28:40.050000+00:00",
  "zoneRedundant": false
}
```

```azurecli-interactive
az eventhubs eventhub create --name $EVENT_GRID_HUB_NAME --namespace-name $NAMESPACE_NAME --resource-group $RESOURCE_GROUP_NAME
```

<!--expected_similarity=0.5-->
```output
{
  "alternateName": null,
  "clusterArmId": null,
  "createdAt": "2023-02-11T00:27:48.977000+00:00",
  "disableLocalAuth": false,
  "encryption": null,
  "id": "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/aksQuickstartResourceGroup/providers/Microsoft.EventHub/namespaces/aksQuickstartNamespace00021677",
  "identity": null,
  "isAutoInflateEnabled": false,
  "kafkaEnabled": true,
  "location": "East US",
  "maximumThroughputUnits": 0,
  "metricId": "325e7c34-99fb-4190-aa87-1df746c67705:aksquickstartnamespace00021677",
  "minimumTlsVersion": "1.2",
  "name": "aksQuickstartNamespace00021677",
  "privateEndpointConnections": null,
  "provisioningState": "Succeeded",
  "publicNetworkAccess": "Enabled",
  "resourceGroup": "aksQuickstartResourceGroup",
  "serviceBusEndpoint": "https://aksQuickstartNamespace00021677.servicebus.windows.net:443/",
  "sku": {
    "capacity": 1,
    "name": "Standard",
    "tier": "Standard"
  },
  "status": "Active",
  "systemData": null,
  "tags": {},
  "type": "Microsoft.EventHub/Namespaces",
  "updatedAt": "2023-02-11T00:29:54.450000+00:00",
  "zoneRedundant": false
}
```

> [!NOTE]
> The *name* of your namespace must be unique. In the defaults above we set a random postfix to try to ensure it is unique, but this is not guaranteed.

Subscribe to the AKS events using [az eventgrid event-subscription create][az-eventgrid-event-subscription-create]:

First we need the resource ID and endpoint, which we will store in an environment variables for later use:

```azurecli-interactive
SOURCE_RESOURCE_ID=$(az aks show -g $RESOURCE_GROUP_NAME -n $AKS_CLUSTER_NAME --query id --output tsv)
ENDPOINT=$(az eventhubs eventhub show -g $RESOURCE_GROUP_NAME -n $EVENT_GRID_HUB_NAME --namespace-name $NAMESPACE_NAME --query id --output tsv)
```

Now we can actually subscribe to the events:

```azurecli-interactive
az eventgrid event-subscription create --name $EVENT_GRID_SUBSCRIPTION_NAME \
  --source-resource-id $SOURCE_RESOURCE_ID \
  --endpoint-type eventhub \
  --endpoint $ENDPOINT
```

<!--expected_similarity=0.5-->
```output
{
  "deadLetterDestination": null,
  "deadLetterWithResourceIdentity": null,
  "deliveryWithResourceIdentity": null,
  "destination": {
    "deliveryAttributeMappings": null,
    "endpointType": "EventHub",
    "resourceId": "/subscriptions/REDACTED/resourceGroups/aksQuickstartResourceGroup/providers/Microsoft.EventHub/namespaces/aksQuickstartNamespace00006800/eventhubs/aksQuickstartEventGridHub"
  },
  "eventDeliverySchema": "EventGridSchema",
  "expirationTimeUtc": null,
  "filter": {
    "advancedFilters": null,
    "enableAdvancedFilteringOnArrays": null,
    "includedEventTypes": [
      "Microsoft.ContainerService.NewKubernetesVersionAvailable"
    ],
    "isSubjectCaseSensitive": null,
    "subjectBeginsWith": "",
    "subjectEndsWith": ""
  },
  "id": "/subscriptions/REDACTED/resourceGroups/aksQuickstartResourceGroup/providers/Microsoft.ContainerService/managedClusters/aksQuickstartCluster/providers/Microsoft.EventGrid/eventSubscriptions/aksQuickstartEventGridSubscription",
  "labels": null,
  "name": "aksQuickstartEventGridSubscription",
  "provisioningState": "Succeeded",
  "resourceGroup": "aksQuickstartResourceGroup",
  "retryPolicy": {
    "eventTimeToLiveInMinutes": 1440,
    "maxDeliveryAttempts": 30
  },
  "systemData": null,
  "topic": "/subscriptions/REDACTED/resourceGroups/aksquickstartresourcegroup/providers/microsoft.containerservice/managedclusters/aksquickstartcluster",
  "type": "Microsoft.EventGrid/eventSubscriptions"
}
```

Verify your subscription to AKS events using `az eventgrid event-subscription list`:

```azurecli-interactive
az eventgrid event-subscription list --source-resource-id $SOURCE_RESOURCE_ID
```

The following example output shows you're subscribed to events from the *MyAKS* cluster and those events are delivered to the *MyEventGridHub* event hub:
<!--expected_similarity=0.5-->
```output
[
  {
    "deadLetterDestination": null,
    "deadLetterWithResourceIdentity": null,
    "deliveryWithResourceIdentity": null,
    "destination": {
      "deliveryAttributeMappings": null,
      "endpointType": "EventHub",
      "resourceId": "/subscriptions/SUBSCRIPTION_ID/resourceGroups/MyResourceGroup/providers/Microsoft.EventHub/namespaces/MyNamespace/eventhubs/MyEventGridHub"
    },
    "eventDeliverySchema": "EventGridSchema",
    "expirationTimeUtc": null,
    "filter": {
      "advancedFilters": null,
      "enableAdvancedFilteringOnArrays": null,
      "includedEventTypes": [
        "Microsoft.ContainerService.NewKubernetesVersionAvailable"
      ],
      "isSubjectCaseSensitive": null,
      "subjectBeginsWith": "",
      "subjectEndsWith": ""
    },
    "id": "/subscriptions/SUBSCRIPTION_ID/resourceGroups/MyResourceGroup/providers/Microsoft.ContainerService/managedClusters/MyAKS/providers/Microsoft.EventGrid/eventSubscriptions/MyEventGridSubscription",
    "labels": null,
    "name": "MyEventGridSubscription",
    "provisioningState": "Succeeded",
    "resourceGroup": "MyResourceGroup",
    "retryPolicy": {
      "eventTimeToLiveInMinutes": 1440,
      "maxDeliveryAttempts": 30
    },
    "systemData": null,
    "topic": "/subscriptions/SUBSCRIPTION_ID/resourceGroups/MyResourceGroup/providers/microsoft.containerservice/managedclusters/MyAKS",
    "type": "Microsoft.EventGrid/eventSubscriptions"
  }
]
```

### [Azure PowerShell](#tab/azure-powershell)

Create a namespace and event hub using [New-AzEventHubNamespace][new-azeventhubnamespace] and [New-AzEventHub][new-azeventhub]. The following example creates a namespace *MyNamespace* and an event hub *MyEventGridHub* in *MyNamespace*, both in the *MyResourceGroup* resource group.

```azurepowershell-interactive
New-AzEventHubNamespace -Location eastus -Name MyNamespace -ResourceGroupName $RESOURCE_GROUP_NAME
New-AzEventHub -Name MyEventGridHub -Namespace MyNamespace -ResourceGroupName $RESOURCE_GROUP_NAME
```

> [!NOTE]
> The *name* of your namespace must be unique.

Subscribe to the AKS events using [New-AzEventGridSubscription][new-azeventgridsubscription]:

```azurepowershell-interactive
$SOURCE_RESOURCE_ID = (Get-AzAksCluster -ResourceGroupName MyResourceGroup -Name MyAKS).Id
$ENDPOINT = (Get-AzEventHub -ResourceGroupName MyResourceGroup -EventHubName MyEventGridHub -Namespace MyNamespace).Id
$params = @{
    EventSubscriptionName = 'MyEventGridSubscription'
    ResourceId            = $SOURCE_RESOURCE_ID
    EndpointType          = 'eventhub'
    Endpoint              = $ENDPOINT 
}
New-AzEventGridSubscription @params
```

Verify your subscription to AKS events using `Get-AzEventGridSubscription`:

```azurepowershell-interactive
Get-AzEventGridSubscription -ResourceId $SOURCE_RESOURCE_ID | Select-Object -ExpandProperty PSEventSubscriptionsList
```

The following example output shows you're subscribed to events from the *MyAKS* cluster and those events are delivered to the *MyEventGridHub* event hub:

```Output
EventSubscriptionName : MyEventGridSubscription
Id                    : /subscriptions/SUBSCRIPTION_ID/resourceGroups/MyResourceGroup/providers/Microsoft.ContainerService/managedClusters/MyAKS/providers/Microsoft.EventGrid/eventSubscriptions/MyEventGridSubscription
Type                  : Microsoft.EventGrid/eventSubscriptions
Topic                 : /subscriptions/SUBSCRIPTION_ID/resourceGroups/myresourcegroup/providers/microsoft.containerservice/managedclusters/myaks
Filter                : Microsoft.Azure.Management.EventGrid.Models.EventSubscriptionFilter
Destination           : Microsoft.Azure.Management.EventGrid.Models.EventHubEventSubscriptionDestination
ProvisioningState     : Succeeded
Labels                : 
EventTtl              : 1440
MaxDeliveryAttempt    : 30
EventDeliverySchema   : EventGridSchema
ExpirationDate        : 
DeadLetterEndpoint    : 
Endpoint              : /subscriptions/SUBSCRIPTION_ID/resourceGroups/MyResourceGroup/providers/Microsoft.EventHub/namespaces/MyNamespace/eventhubs/MyEventGridHub
```

---

When AKS events occur, you'll see those events appear in your event hub. For example, when the list of available Kubernetes versions for your clusters changes, you'll see a `Microsoft.ContainerService.NewKubernetesVersionAvailable` event. For more information on the events AKS emits, see [Azure Kubernetes Service (AKS) as an Event Grid source][aks-events].

## Delete the cluster and subscriptions

### [Azure CLI](#tab/azure-cli)

Use the [az group delete][az-group-delete] command to remove the resource group, the AKS cluster, namespace, and event hub, and all related resources.

```azurecli-interactive
az group delete --name $RESOURCE_GROUP_NAME --yes --no-wait
```

### [Azure PowerShell](#tab/azure-powershell)

Use the [Remove-AzResourceGroup][remove-azresourcegroup] cmdlet to remove the resource group, the AKS cluster, namespace, and event hub, and all related resources.

```azurepowershell-interactive
Remove-AzResourceGroup -Name MyResourceGroup
```

---

> [!NOTE]
> When you delete the cluster, the Azure Active Directory service principal used by the AKS cluster is not removed. For steps on how to remove the service principal, see [AKS service principal considerations and deletion][sp-delete].
> 
> If you used a managed identity, the identity is managed by the platform and does not require removal.

## Next steps

In this quickstart, you deployed a Kubernetes cluster and then subscribed to AKS events in Azure Event Hubs.

To learn more about AKS, and walk through a complete code to deployment example, continue to the Kubernetes cluster tutorial.

> [!div class="nextstepaction"]
> [AKS tutorial][aks-tutorial]

[azure-cli-install]: /cli/azure/install-azure-cli
[azure-powershell-install]: /powershell/azure/install-az-ps
[aks-events]: ../event-grid/event-schema-aks.md
[aks-tutorial]: ./tutorial-kubernetes-prepare-app.md
[az-aks-create]: /cli/azure/aks#az_aks_create
[new-azakscluster]: /powershell/module/az.aks/new-azakscluster
[az-eventhubs-namespace-create]: /cli/azure/eventhubs/namespace#az-eventhubs-namespace-create
[new-azeventhubnamespace]: /powershell/module/az.eventhub/new-azeventhubnamespace
[az-eventhubs-eventhub-create]: /cli/azure/eventhubs/eventhub#az-eventhubs-eventhub-create
[new-azeventhub]: /powershell/module/az.eventhub/new-azeventhub
[az-eventgrid-event-subscription-create]: /cli/azure/eventgrid/event-subscription#az-eventgrid-event-subscription-create
[new-azeventgridsubscription]: /powershell/module/az.eventgrid/new-azeventgridsubscription
[az-group-delete]: /cli/azure/group#az_group_delete
[sp-delete]: kubernetes-service-principal.md#other-considerations
[remove-azresourcegroup]: /powershell/module/az.resources/remove-azresourcegroup