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
```azurecli-interactive
export RESOURCE_GROUP_NAME=myResourceGroup
export RESOURCE_LOCATION=eastus
export AKS_CLUSTER_NAME=myAKSCluster
export NAMESPACE_NAME="myNamespace$(printf "%08d" $((RANDOM%100000000)))"
export EVENT_GRID_HUB_NAME=myEventGridHub
export EVENT_GRID_SUBSCRIPTION_NAME=myEventGridSubscription
```

Create an AKS cluster using the [az aks create][az-aks-create] command. The following example creates a resource group *MyResourceGroup* and a cluster named *MyAKS* with one node in the *MyResourceGroup* resource group:

```azurecli-interactive
az group create --name $RESOURCE_GROUP_NAME --location $RESOURCE_LOCATION
az aks create --resource-group $RESOURCE_GROUP_NAME --name $AKS_CLUSTER_NAME --location $RESOURCE_LOCATION  --node-count 1 --generate-ssh-keys
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
az eventhubs eventhub create --name $EVENT_GRID_HUB_NAME --namespace-name $NAMESPACE_NAME --resource-group $RESOURCE_GROUP_NAME
```
<!--expected_similarity=0.7-->
```output
{
  "alternateName": null,
  "clusterArmId": null,
  "createdAt": "2023-02-07T23:41:18.287000+00:00",
  "disableLocalAuth": false,
  "encryption": null,
  "id": "/subscriptions/f7a60fca-9977-4899-b907-005a076adbb6/resourceGroups/myResourceGroup3/providers/Microsoft.EventHub/namespaces/myNamespace00022998",
  "identity": null,
  "isAutoInflateEnabled": false,
  "kafkaEnabled": true,
  "location": "East US",
  "maximumThroughputUnits": 0,
  "metricId": "f7a60fca-9977-4899-b907-005a076adbb6:mynamespace00022998",
  "minimumTlsVersion": "1.2",
  "name": "myNamespace00022998",
  "privateEndpointConnections": null,
  "provisioningState": "Succeeded",
  "publicNetworkAccess": "Enabled",
  "resourceGroup": "myResourceGroup3",
  "serviceBusEndpoint": "https://myNamespace00022998.servicebus.windows.net:443/",
  "sku": {
    "capacity": 1,
    "name": "Standard",
    "tier": "Standard"
  },
  "status": "Active",
  "systemData": null,
  "tags": {},
  "type": "Microsoft.EventHub/Namespaces",
  "updatedAt": "2023-02-07T23:42:11.013000+00:00",
  "zoneRedundant": false
}
{
  "captureDescription": null,
  "createdAt": "2023-02-07T23:42:23.190000+00:00",
  "id": "/subscriptions/f7a60fca-9977-4899-b907-005a076adbb6/resourceGroups/myResourceGroup3/providers/Microsoft.EventHub/namespaces/myNamespace00022998/eventhubs/myEventGridHub",
  "location": "eastus",
  "messageRetentionInDays": 7,
  "name": "myEventGridHub",
  "partitionCount": 4,
  "partitionIds": [
    "0",
    "1",
    "2",
    "3"
  ],
  "resourceGroup": "myResourceGroup3",
  "status": "Active",
  "systemData": null,
  "type": "Microsoft.EventHub/namespaces/eventhubs",
  "updatedAt": "2023-02-07T23:42:23.427000+00:00"
}
```
> [!NOTE]
> The *name* of your namespace must be unique.

Subscribe to the AKS events using [az eventgrid event-subscription create][az-eventgrid-event-subscription-create]:

```azurecli-interactive
SOURCE_RESOURCE_ID=$(az aks show -g $RESOURCE_GROUP_NAME -n $AKS_CLUSTER_NAME --query id --output tsv)
ENDPOINT=$(az eventhubs eventhub show -g $RESOURCE_GROUP_NAME -n $EVENT_GRID_HUB_NAME --namespace-name $NAMESPACE_NAME --query id --output tsv)
az eventgrid event-subscription create --name $EVENT_GRID_SUBSCRIPTION_NAME \
--source-resource-id $SOURCE_RESOURCE_ID \
--endpoint-type eventhub \
--endpoint $ENDPOINT
```

Verify your subscription to AKS events using `az eventgrid event-subscription list`:

```azurecli-interactive
az eventgrid event-subscription list --source-resource-id $SOURCE_RESOURCE_ID
```

The following example output shows you're subscribed to events from the *MyAKS* cluster and those events are delivered to the *MyEventGridHub* event hub:
<!--expected_similarity=0.7-->
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