# Quickstart: Deploy a Scalable & Secure Azure Kubernetes Service cluster using the Azure CLI
Welcome to this tutorial where we will take you step by step in creating an Azure Kubernetes Web Application that is secured via https. This tutorial assumes you are logged into Azure CLI already and have selected a subscription to use with the CLI. It also assumes that you have Helm installed (Instructions can be found here https://helm.sh/docs/intro/install/).

## Define Environment Variables

The First step in this tutorial is to define environment variables 

```bash
export UNIQUE_POSTFIX="$(($RANDOM % 254 + 1))"
export MY_RESOURCE_GROUP_NAME="myResourceGroup$UNIQUE_POSTFIX"
export MY_LOCATION="eastus"
export MY_AKS_CLUSTER_NAME="myAKSCluster$UNIQUE_POSTFIX"
export MY_PUBLIC_IP_NAME="myPublicIP$UNIQUE_POSTFIX"
export MY_DNS_LABEL="mydnslabel$UNIQUE_POSTFIX"
export MY_VNET_NAME="myVNet$UNIQUE_POSTFIX"
export MY_VNET_PREFIX="10.$UNIQUE_POSTFIX.0.0/16"
export MY_SN_NAME="mySN$UNIQUE_POSTFIX"
export MY_SN_PREFIX="10.$UNIQUE_POSTFIX.0.0/22"
```

# Create a resource group

A resource group is a container for related resources. All resources must be placed in a resource group. We will create one for this tutorial. The following command creates a resource group with the previously defined $MY_RESOURCE_GROUP_NAME and $MY_LOCATION parameters.

```bash
az group create --name $MY_RESOURCE_GROUP_NAME --location $MY_LOCATION
```
Results:

```expected_similarity=0.3
{
  "id": "/subscriptions/bb318642-28fd-482d-8d07-79182df07999/resourceGroups/myResourceGroup210",
  "location": "eastus",
  "managedBy": null,
  "name": "testResourceGroup",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```

## Create a virtual network and subnet

A virtual network is the fundamental building block for private networks in Azure. Azure Virtual Network enables Azure resources like VMs to securely communicate with each other and the internet.

```bash
az network vnet create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --location $MY_LOCATION \
    --name $MY_VNET_NAME \
    --address-prefix $MY_VNET_PREFIX \
    --subnet-name $MY_SN_NAME \
    --subnet-prefixes $MY_SN_PREFIX
```
Results:

```expected_similarity=0.3
{
  "newVNet": {
    "addressSpace": {
      "addressPrefixes": [
        "10.210.0.0/16"
      ]
    },
    "enableDdosProtection": false,
    "etag": "W/\"1e065114-2ae3-4dee-91eb-c69667e60afb\"",
    "id": "/subscriptions/bb318642-28fd-482d-8d07-79182df07999/myResourceGroup210/providers/Microsoft.Network/virtualNetworks/myVNet210",
    "location": "eastus",
    "name": "myVNet210",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup210",
    "resourceGuid": "3e54a2e8-32fa-4157-b817-f4e4507dbac9",
    "subnets": [
      {
        "addressPrefix": "10.210.0.0/22",
        "delegations": [],
        "etag": "W/\"1e065114-2ae3-4dee-91eb-c69667e60afb\"",
        "id": "/subscriptions/bb318642-28fd-482d-8d07-79182df07999/myResourceGroup210/providers/Microsoft.Network/virtualNetworks/myVNet210/subnets/mySN210",
        "name": "mySN210",
        "privateEndpointNetworkPolicies": "Disabled",
        "privateLinkServiceNetworkPolicies": "Enabled",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup210",
        "type": "Microsoft.Network/virtualNetworks/subnets"
      }
    ],
    "type": "Microsoft.Network/virtualNetworks",
    "virtualNetworkPeerings": []
  }
}
```

## Register to AKS Azure Resource Providers
Verify Microsoft.OperationsManagement and Microsoft.OperationalInsights providers are registered on your subscription. These are Azure resource providers required to support [Container insights](https://docs.microsoft.com/en-us/azure/azure-monitor/containers/container-insights-overview). To check the registration status, run the following commands

```bash
az provider register --namespace Microsoft.OperationsManagement
az provider register --namespace Microsoft.OperationalInsights
```

## Create AKS Cluster 
Create an AKS cluster using the az aks create command with the --enable-addons monitoring parameter to enable Container insights. The following example creates an autoscaling, availability zone enabled cluster named myAKSCluster:

This will take a few minutes
```bash
export MY_SN_ID=$(az network vnet subnet list --resource-group $MY_RESOURCE_GROUP_NAME --vnet-name $MY_VNET_NAME --query "[0].id" --output tsv)

az aks create \
  --resource-group $MY_RESOURCE_GROUP_NAME \
  --name $MY_AKS_CLUSTER_NAME \
  --auto-upgrade-channel stable \
  --enable-cluster-autoscaler \
  --enable-addons monitoring \
  --location $MY_LOCATION \
  --node-count 1 \
  --min-count 1 \
  --max-count 3 \
  --network-plugin azure \
  --network-policy azure \
  --vnet-subnet-id $MY_SN_ID \
  --no-ssh-key \
  --node-vm-size Standard_DS2_v2 \
  --zones 1 2 3
```

## Connect to the cluster
To manage a Kubernetes cluster, use the Kubernetes command-line client, kubectl. kubectl is already installed if you use Azure Cloud Shell.

1. Install az aks CLI locally using the az aks install-cli command

```bash
if ! [ -x "$(command -v kubectl)" ]; then az aks install-cli; fi
```

2. Configure kubectl to connect to your Kubernetes cluster using the az aks get-credentials command. The following command:
    - Downloads credentials and configures the Kubernetes CLI to use them.
    - Uses ~/.kube/config, the default location for the Kubernetes configuration file. Specify a different location for your Kubernetes configuration file using --file argument. 

> [!WARNING]
> This will overwrite any existing credentials with the same entry

```bash
az aks get-credentials --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_AKS_CLUSTER_NAME --overwrite-existing
```

3. Verify the connection to your cluster using the kubectl get command. This command returns a list of the cluster nodes.

```bash
kubectl get nodes
```

## Install NGINX Ingress Controller

```bash
export MY_STATIC_IP=$(az network public-ip create --resource-group MC_${MY_RESOURCE_GROUP_NAME}_${MY_AKS_CLUSTER_NAME}_${MY_LOCATION} --location  ${MY_LOCATION} --name ${MY_PUBLIC_IP_NAME} --dns-name ${MY_DNS_LABEL} --sku Standard --allocation-method static --version IPv4 --zone 1 2 3 --query publicIp.ipAddress -o tsv)
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-dns-label-name"=$MY_DNS_LABEL \
  --set controller.service.loadBalancerIP=$MY_STATIC_IP \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-load-balancer-health-probe-request-path"=/healthz
```

## Deploy the Application 
A Kubernetes manifest file defines a cluster's desired state, such as which container images to run.

In this quickstart, you will use a manifest to create all objects needed to run the Azure Vote application. This manifest includes two Kubernetes deployments:

- The sample Azure Vote Python applications.
- A Redis instance.

Two Kubernetes Services are also created:

- An internal service for the Redis instance.
- An external service to access the Azure Vote application from the internet.

Finally, an Ingress resource is created to route traffic to the Azure Vote application.

A test voting app YML file is already prepared. To deploy this app run the following command 
```bash
kubectl apply -f azure-vote-start.yml
```

## Test The Application

Validate that the application is running by either visiting the public ip or the application url. The application url can be found by running the following command:
```bash
curl "http://${MY_DNS_LABEL}.${MY_LOCATION}.cloudapp.azure.com"
```

# Add Application Gateway Ingress Controller
The Application Gateway Ingress Controller (AGIC) is a Kubernetes application, which makes it possible for Azure Kubernetes Service (AKS) customers to leverage Azure's native Application Gateway L7 load-balancer to expose cloud software to the Internet. AGIC monitors the Kubernetes cluster it is hosted on and continuously updates an Application Gateway, so that selected services are exposed to the Internet

AGIC helps eliminate the need to have another load balancer/public IP in front of the AKS cluster and avoids multiple hops in your datapath before requests reach the AKS cluster. Application Gateway talks to pods using their private IP directly and does not require NodePort or KubeProxy services. This also brings better performance to your deployments.

## Deploy a new Application Gateway 
1. Create a Public IP for Application Gateway by running the following:
```bash
az network public-ip create --name $PUBLIC_IP_NAME --resource-group $RESOURCE_GROUP_NAME --allocation-method Static --sku Standard
```

2. Create a Virtual Network(Vnet) for Application Gateway by running the following:
```bash
az network vnet create --name $VNET_NAME --resource-group $RESOURCE_GROUP_NAME --address-prefix 11.0.0.0/8 --subnet-name $SUBNET_NAME --subnet-prefix 11.1.0.0/16 
```

3. Create Application Gateway by running the following:

> [!NOTE] 
> This will take around 5 minutes 
```bash
az network application-gateway create --name $APPLICATION_GATEWAY_NAME --location $RESOURCE_LOCATION --resource-group $RESOURCE_GROUP_NAME --sku Standard_v2 --public-ip-address $PUBLIC_IP_NAME --vnet-name $VNET_NAME --subnet $SUBNET_NAME --priority 100
```

## Enable the AGIC add-on in existing AKS cluster 

1. Store Application Gateway ID by running the following:
```bash
APPLICATION_GATEWAY_ID=$(az network application-gateway show --name $APPLICATION_GATEWAY_NAME --resource-group $RESOURCE_GROUP_NAME --output tsv --query "id") 
```

2. Enable Application Gateway Ingress Addon by running the following:

> [!NOTE]
> This will take a few minutes
```bash
az aks enable-addons --name $AKS_CLUSTER_NAME --resource-group $RESOURCE_GROUP_NAME --addon ingress-appgw --appgw-id $APPLICATION_GATEWAY_ID
```

3. Store the node resource as an environment variable group by running the following:
```bash
NODE_RESOURCE_GROUP=$(az aks show --name myAKSCluster --resource-group $RESOURCE_GROUP_NAME --output tsv --query "nodeResourceGroup")
```
4. Store the Vnet name as an environment variable by running the following:
```bash
AKS_VNET_NAME=$(az network vnet list --resource-group $NODE_RESOURCE_GROUP --output tsv --query "[0].name")
```

5. Store the Vnet ID as an environment variable by running the following:
```bash
AKS_VNET_ID=$(az network vnet show --name $AKS_VNET_NAME --resource-group $NODE_RESOURCE_GROUP --output tsv --query "id")
```
## Peer the two virtual networks together 
Since we deployed the AKS cluster in its own virtual network and the Application Gateway in another virtual network, you'll need to peer the two virtual networks together in order for traffic to flow from the Application Gateway to the pods in the cluster. Peering the two virtual networks requires running the Azure CLI command two separate times, to ensure that the connection is bi-directional. The first command will create a peering connection from the Application Gateway virtual network to the AKS virtual network; the second command will create a peering connection in the other direction.

1. Create the peering from Application Gateway to AKS by runnig the following:
```bash
az network vnet peering create --name $APPGW_TO_AKS_PEERING_NAME --resource-group $RESOURCE_GROUP_NAME --vnet-name $VNET_NAME --remote-vnet $AKS_VNET_ID --allow-vnet-access 
```

2. Store Id of Application Gateway Vnet As enviornment variable by running the following:
```bash
APPLICATION_GATEWAY_VNET_ID=$(az network vnet show --name $VNET_NAME --resource-group $RESOURCE_GROUP_NAME --output tsv --query "id")
```
3. Create Vnet Peering from AKS to Application Gateway
```bash
az network vnet peering create --name $AKS_TO_APPGW_PEERING_NAME --resource-group $NODE_RESOURCE_GROUP --vnet-name $AKS_VNET_NAME --remote-vnet $APPLICATION_GATEWAY_VNET_ID --allow-vnet-access
```
4. Store New IP address as environment variable by running the following command:
```bash
runtime="2 minute"; endtime=$(date -ud "$runtime" +%s); while [[ $(date -u +%s) -le $endtime ]]; do export IP_ADDRESS=$(az network public-ip show --resource-group $RESOURCE_GROUP_NAME --name $PUBLIC_IP_NAME --query ipAddress --output tsv); if ! [ -z $IP_ADDRESS ]; then break; else sleep 10; fi; done
```

## Apply updated application YAML complete with AGIC
In order to use the Application Gateway Ingress Controller we deployed we need to re-deploy an update Voting App YML file. The following command will update the application:

The full updated YML file can be viewed at `azure-vote-agic-yml`
```bash
kubectl apply -f azure-vote-agic.yml
```

## Check that the application is reachable
Now that the Application Gateway is set up to serve traffic to the AKS cluster, let's verify that your application is reachable. 

Check that the sample application you created is up and running by either visiting the IP address of the Application Gateway that get from running the following command or check with curl. It may take Application Gateway a minute to get the update, so if the Application Gateway is still in an "Updating" state on Portal, then let it finish before trying to reach the IP address. Run the following to check the status:
```bash
kubectl get ingress
```

## Add custom subdomain to AGIC
Now Application Gateway Ingress has been added to the application gateway the next step is to add a custom domain. This will allow the endpoint to be reached by a human readable URL as well as allow for SSL Termination at the endpoint.

1. Store Unique ID of the Public IP Address as an environment variable by running the following:
```bash
export PUBLIC_IP_ID=$(az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '$IP_ADDRESS')].[id]" --output tsv)
```

2. Update public IP to respond to custom domain requests by running the following:
```bash
az network public-ip update --ids $PUBLIC_IP_ID --dns-name $CUSTOM_DOMAIN_NAME
```

3. Validate the resource is reachable via the custom domain.
```bash
az network public-ip show --ids $PUBLIC_IP_ID --query "[dnsSettings.fqdn]" --output tsv
```

4. Store the custom domain as en enviornment variable. This will be used later when setting up https termination.
```bash
export FQDN=$(az network public-ip show --ids $PUBLIC_IP_ID --query "[dnsSettings.fqdn]" --output tsv)
```

# Add HTTPS termination to custom domain 
At this point in the tutorial you have an AKS web app with Application Gateway as the Ingress controller and a custom domain you can use to access your application. The next step is to add an SSL certificate to the domain so that users can reach your application securely via https.  

## Set Up Cert Manager
In order to add HTTPS we are going to use Cert Manager. Cert Manager is an open source tool used to obtain and manage SSL certificate for Kubernetes deployments. Cert Manager will obtain certificates from a variety of Issuers, both popular public Issuers as well as private Issuers, and ensure the certificates are valid and up-to-date, and will attempt to renew certificates at a configured time before expiry.

1. In order to install cert-manager, we must first create a namespace to run it in. This tutorial will install cert-manager into the cert-manager namespace. It is possible to run cert-manager in a different namespace, although you will need to make modifications to the deployment manifests.
```bash
kubectl create namespace cert-manager
```

2. We can now install cert-manager. All resources are included in a single YAML manifest file. This can be installed by running the following:
```bash
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.7.0/cert-manager.crds.yaml
```


3. Add the certmanager.k8s.io/disable-validation: "true" label to the cert-manager namespace by running the following. This will allow the system resources that cert-manager requires to bootstrap TLS to be created in its own namespace.
```bash
kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true
```

## Obtain certificate via Helm Charts
Helm is a Kubernetes deployment tool for automating creation, packaging, configuration, and deployment of applications and services to Kubernetes clusters.

Cert-manager provides Helm charts as a first-class method of installation on Kubernetes.

1. Add the Jetstack Helm repository
This repository is the only supported source of cert-manager charts. There are some other mirrors and copies across the internet, but those are entirely unofficial and could present a security risk.
```bash
helm repo add jetstack https://charts.jetstack.io
```

2. Update local Helm Chart repository cache 
```bash
helm repo update
```

3. Install Cert-Manager addon via helm by running the following:
```bash
helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v1.7.0
```

4. Apply Certificate Issuer YAML File

    ClusterIssuers are Kubernetes resources that represent certificate authorities (CAs) that are able to generate signed certificates by honoring certificate signing requests. All cert-manager certificates require a referenced issuer that is in a ready condition to attempt to honor the request.

    The issuer we are using can be found in the `cluster-issuer-prod.yaml file`
```bash
envsubst < cluster-issuer-prod.yaml | kubectl apply -f -
```

5. Upate Voting App Application to use Cert-Manager to obtain an SSL Certificate. 

    The full YAML file can be found in `azure-vote-agic-ssl-yml`
```bash
envsubst < azure-vote-agic-ssl.yml | kubectl apply -f -
```
## Validate application is working

Wait for SSL certificate to issue. The following command will query the status of the SSL certificate for 3 minutes.
 In rare occasions it may take up to 15 minutes for Lets Encrypt to issue a successful challenge and the ready state to be 'True'
```bash
runtime="10 minute"; endtime=$(date -ud "$runtime" +%s); while [[ $(date -u +%s) -le $endtime ]]; do STATUS=$(kubectl get certificate --output jsonpath={..status.conditions[0].status}); echo $STATUS; if [ "$STATUS" = 'True' ]; then break; else sleep 10; fi; done
```

Validate SSL certificate is True by running the follow command:
```bash
kubectl get certificate --output jsonpath={..status.conditions[0].status}
```

Results:

```expected_similarity=0.8
True
```

## Browse your AKS Deployment Secured via HTTPS!
Run the following command to get the HTTPS endpoint for your application:

>[!Note]
> It often takes 2-3 minutes for the SSL certificate to propogate and the site to be reachable via https 
```bash
echo https://$FQDN
```
Paste this into the browser to validate your deployment.

## Next Steps

* [Azure Kubernetes Service Documentation](https://learn.microsoft.com/en-us/azure/aks/)
* [Create an Azure Container Registry](https://learn.microsoft.com/en-us/azure/aks/tutorial-kubernetes-prepare-acr?tabs=azure-cli)
* [Scale your Applciation in AKS](https://learn.microsoft.com/en-us/azure/aks/tutorial-kubernetes-scale?tabs=azure-cli)
* [Update your application in AKS](https://learn.microsoft.com/en-us/azure/aks/tutorial-kubernetes-app-update?tabs=azure-cli)
