# Quickstart: Deploy a Scalable & Secure WordPress instance on AKS

Welcome to this tutorial where we will take you step by step in creating an Azure Kubernetes Web Application that is secured via https. This tutorial assumes you are logged into Azure CLI already and have selected a subscription to use with the CLI. It also assumes that you have Helm installed ([Instructions can be found here](https://helm.sh/docs/intro/install/)).

## Define Environment Variables

The First step in this tutorial is to define environment variables.

```bash
export SSL_EMAIL_ADDRESS="$(az account show --query user.name --output tsv)"
export NETWORK_PREFIX="$(($RANDOM % 253 + 1))"
export RANDOM_ID="$(openssl rand -hex 3)"
export MY_RESOURCE_GROUP_NAME="myResourceGroup$RANDOM_ID"
export REGION="eastus"
export MY_AKS_CLUSTER_NAME="myAKSCluster$RANDOM_ID"
export MY_PUBLIC_IP_NAME="myPublicIP$RANDOM_ID"
export MY_DNS_LABEL="mydnslabel$RANDOM_ID"
export MY_VNET_NAME="myVNet$RANDOM_ID"
export MY_VNET_PREFIX="10.$NETWORK_PREFIX.0.0/16"
export MY_SN_NAME="mySN$RANDOM_ID"
export MY_SN_PREFIX="10.$NETWORK_PREFIX.0.0/22"
export MY_MYSQL_DB_NAME="mydb$RANDOM_ID"
export MY_MYSQL_ADMIN_USERNAME="dbadmin$RANDOM_ID"
export MY_MYSQL_ADMIN_PW="$(openssl rand -base64 32)"
export MY_MYSQL_SN_NAME="myMySQLSN$RANDOM_ID"
export MY_MYSQL_HOSTNAME="$MY_MYSQL_DB_NAME.mysql.database.azure.com"
export MY_WP_ADMIN_PW="$(openssl rand -base64 32)"
export MY_WP_ADMIN_USER="wpcliadmin"
export FQDN="${MY_DNS_LABEL}.${REGION}.cloudapp.azure.com"
```

## Create a resource group

A resource group is a container for related resources. All resources must be placed in a resource group. We will create one for this tutorial. The following command creates a resource group with the previously defined $MY_RESOURCE_GROUP_NAME and $REGION parameters.

```bash
az group create \
    --name $MY_RESOURCE_GROUP_NAME \
    --location $REGION
```

Results:

<!-- expected_similarity=0.3 -->
```json
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup210",
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
    --location $REGION \
    --name $MY_VNET_NAME \
    --address-prefix $MY_VNET_PREFIX \
    --subnet-name $MY_SN_NAME \
    --subnet-prefixes $MY_SN_PREFIX
```

Results:

<!-- expected_similarity=0.3 -->
```json
{
  "newVNet": {
    "addressSpace": {
      "addressPrefixes": [
        "10.210.0.0/16"
      ]
    },
    "enableDdosProtection": false,
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/myResourceGroup210/providers/Microsoft.Network/virtualNetworks/myVNet210",
    "location": "eastus",
    "name": "myVNet210",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup210",
    "subnets": [
      {
        "addressPrefix": "10.210.0.0/22",
        "delegations": [],
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/myResourceGroup210/providers/Microsoft.Network/virtualNetworks/myVNet210/subnets/mySN210",
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

## Create an Azure Database for MySQL - Flexible Server

Azure Database for MySQL - Flexible Server is a managed service that you can use to run, manage, and scale highly available MySQL servers in the cloud. Create a flexible server with the [az mysql flexible-server create](https://learn.microsoft.com/cli/azure/mysql/flexible-server#az-mysql-flexible-server-create) command. A server can contain multiple databases. The following command creates a server using service defaults and variable values from your Azure CLI's local environment:

```bash
echo "Your MySQL user $MY_MYSQL_ADMIN_USERNAME password is: $MY_WP_ADMIN_PW" 
```

```bash
az mysql flexible-server create \
    --admin-password $MY_MYSQL_ADMIN_PW \
    --admin-user $MY_MYSQL_ADMIN_USERNAME \
    --auto-scale-iops Disabled \
    --high-availability Disabled \
    --iops 500 \
    --location $REGION \
    --name $MY_MYSQL_DB_NAME \
    --database-name wordpress \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --sku-name Standard_B2s \
    --storage-auto-grow Disabled \
    --storage-size 20 \
    --subnet $MY_MYSQL_SN_NAME \
    --private-dns-zone $MY_DNS_LABEL.private.mysql.database.azure.com \
    --tier Burstable \
    --version 8.0.21 \
    --vnet $MY_VNET_NAME \
    --yes -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```json
{
  "databaseName": "wp001",
  "host": "mydb6ad2bc.mysql.database.azure.com",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.DBforMySQL/flexibleServers/mydb6ad2bc",
  "location": "East US",
  "resourceGroup": "myResourceGroup6ad2bc",
  "skuname": "Standard_B2s",
  "subnetId": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/virtualNetworks/myVNet6ad2bc/subnets/myMySQLSN6ad2bc",
  "username": "dbadmin6ad2bc",
  "version": "8.0.21"
}
```

The server created has the below attributes:

- The server name, admin username, admin password, resource group name, location are already specified in local context environment of the cloud shell, and will be created in the same location as your the resource group and the other Azure components.
- Service defaults for remaining server configurations: compute tier (Burstable), compute size/SKU (Standard_B2s), backup retention period (7 days), and MySQL version (8.0.21)
- The default connectivity method is Private access (VNet Integration) with a linked virtual network and a auto-generated subnet.

> [!NOTE]
> The connectivity method cannot be changed after creating the server. For example, if you selected `Private access (VNet Integration)` during create then you cannot change to `Public access (allowed IP addresses)` after create. We highly recommend creating a server with Private access to securely access your server using VNet Integration. Learn more about Private access in the [concepts article](https://learn.microsoft.com/azure/mysql/flexible-server/concepts-networking-vnet).

If you'd like to change any defaults, please refer to the Azure CLI [reference documentation](https://learn.microsoft.com/cli/azure//mysql/flexible-server) for the complete list of configurable CLI parameters.

## Check the Azure Database for MySQL - Flexible Server status

It takes a few minutes to create the Azure Database for MySQL - Flexible Server and supporting resources.

```bash
runtime="10 minute"; endtime=$(date -ud "$runtime" +%s); while [[ $(date -u +%s) -le $endtime ]]; do STATUS=$(az mysql flexible-server show -g $MY_RESOURCE_GROUP_NAME -n $MY_MYSQL_DB_NAME --query state -o tsv); echo $STATUS; if [ "$STATUS" = 'Ready' ]; then break; else sleep 10; fi; done
```

## Configure server parameters in Azure Database for MySQL - Flexible Server

You can manage Azure Database for MySQL - Flexible Server configuration using server parameters. The server parameters are configured with the default and recommended value when you create the server.

Show server parameter details
To show details about a particular parameter for a server, run the [az mysql flexible-server parameter show](https://learn.microsoft.com/cli/azure/mysql/flexible-server/parameter) command.

### Disable Azure Database for MySQL - Flexible Server SSL connection parameter for WordPress integration

You can also modify the value of certain server parameters, which updates the underlying configuration values for the MySQL server engine. To update the server parameter, use the [az mysql flexible-server parameter set](https://learn.microsoft.com/cli/azure/mysql/flexible-server/parameter#az-mysql-flexible-server-parameter-set) command.

```bash
az mysql flexible-server parameter set \
    -g $MY_RESOURCE_GROUP_NAME \
    -s $MY_MYSQL_DB_NAME \
    -n require_secure_transport -v "OFF" -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```json
{
  "allowedValues": "ON,OFF",
  "currentValue": "OFF",
  "dataType": "Enumeration",
  "defaultValue": "ON",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.DBforMySQL/flexibleServers/mydb6ad2bc/configurations/require_secure_transport",
  "isConfigPendingRestart": "False",
  "isDynamicConfig": "True",
  "isReadOnly": "False",
  "name": "require_secure_transport",
  "resourceGroup": "myResourceGroup6ad2bc",
  "source": "user-override",
  "systemData": null,
  "type": "Microsoft.DBforMySQL/flexibleServers/configurations",
  "value": "OFF"
}
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
    --location $REGION \
    --node-count 1 \
    --min-count 1 \
    --max-count 3 \
    --network-plugin azure \
    --network-policy azure \
    --vnet-subnet-id $MY_SN_ID \
    --no-ssh-key \
    --node-vm-size Standard_DS2_v2 \
    --service-cidr 10.255.0.0/24 \
    --dns-service-ip 10.255.0.10 \
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

You can configure your ingress controller with a static public IP address. The static public IP address remains if you delete your ingress controller. The IP address doesn't remain if you delete your AKS cluster.
When you upgrade your ingress controller, you must pass a parameter to the Helm release to ensure the ingress controller service is made aware of the load balancer that will be allocated to it. For the HTTPS certificates to work correctly, you use a DNS label to configure an FQDN for the ingress controller IP address.
Your FQDN should follow this form: $MY_DNS_LABEL.AZURE_REGION_NAME.cloudapp.azure.com.

```bash
export MY_STATIC_IP=$(az network public-ip create --resource-group MC_${MY_RESOURCE_GROUP_NAME}_${MY_AKS_CLUSTER_NAME}_${REGION} --location ${REGION} --name ${MY_PUBLIC_IP_NAME} --dns-name ${MY_DNS_LABEL} --sku Standard --allocation-method static --version IPv4 --zone 1 2 3 --query publicIp.ipAddress -o tsv)
```

Add the --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-dns-label-name"="<DNS_LABEL>" parameter. The DNS label can be set either when the ingress controller is first deployed, or it can be configured later. Add the --set controller.service.loadBalancerIP="<STATIC_IP>" parameter. Specify your own public IP address that was created in the previous step.

1. Add the ingress-nginx Helm repository

    ```bash
    helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
    ```

2. Update local Helm Chart repository cache

    ```bash
    helm repo update
    ```

3. Install ingress-nginx addon via Helm by running the following:

  ```bash
  helm upgrade --install --cleanup-on-fail --atomic ingress-nginx ingress-nginx/ingress-nginx \
      --namespace ingress-nginx \
      --create-namespace \
      --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-dns-label-name"=$MY_DNS_LABEL \
      --set controller.service.loadBalancerIP=$MY_STATIC_IP \
      --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-load-balancer-health-probe-request-path"=/healthz \
      --wait --timeout 10m0s
  ```

## Add HTTPS termination to custom domain

At this point in the tutorial you have an AKS web app with NGINX as the Ingress controller and a custom domain you can use to access your application. The next step is to add an SSL certificate to the domain so that users can reach your application securely via https.

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

3. Install Cert-Manager addon via Helm by running the following:

    ```bash
    helm upgrade --install --cleanup-on-fail --atomic \
        --namespace cert-manager \
        --version v1.7.0 \
        --wait --timeout 10m0s \
        cert-manager jetstack/cert-manager
    ```

4. Apply Certificate Issuer YAML File

    ClusterIssuers are Kubernetes resources that represent certificate authorities (CAs) that are able to generate signed certificates by honoring certificate signing requests. All cert-manager certificates require a referenced issuer that is in a ready condition to attempt to honor the request.
    The issuer we are using can be found in the `cluster-issuer-prod.yml file`

    ```bash
    cluster_issuer_variables=$(<cluster-issuer-prod.yaml)
    echo "${cluster_issuer_variables//\$SSL_EMAIL_ADDRESS/$SSL_EMAIL_ADDRESS}" | kubectl apply -f -
    ```

## Create a custom storage class

The default storage classes suit the most common scenarios, but not all. For some cases, you might want to have your own storage class customized with your own parameters. For example, use the following manifest to configure the mountOptions of the file share.
The default value for fileMode and dirMode is 0755 for Kubernetes mounted file shares. You can specify the different mount options on the storage class object.

```bash
kubectl apply -f wp-azurefiles-sc.yaml
```

## Deploy WordPress to AKS cluster

For this document, we are using an existing Helm Chart for WordPress built by Bitnami. For example Bitnami Helm chart uses local MariaDB as the database and we need to override these values to use the app with Azure Database for MySQL. All of the override values You can override the values and the custom settings can be found in the file `helm-wp-aks-values.yaml`

1. Add the Wordpress Bitnami Helm repository

    ```bash
    helm repo add bitnami https://charts.bitnami.com/bitnami
    ```

2. Update local Helm Chart repository cache

    ```bash
    helm repo update
    ```

3. Install Wordpress workload via Helm by running the following:

    ```bash
    helm upgrade --install --cleanup-on-fail \
        --wait --timeout 10m0s \
        --namespace wordpress \
        --create-namespace \
        --set wordpressUsername="$MY_WP_ADMIN_USER" \
        --set wordpressPassword="$MY_WP_ADMIN_PW" \
        --set wordpressEmail="$SSL_EMAIL_ADDRESS" \
        --set externalDatabase.host="$MY_MYSQL_HOSTNAME" \
        --set externalDatabase.user="$MY_MYSQL_ADMIN_USERNAME" \
        --set externalDatabase.password="$MY_MYSQL_ADMIN_PW" \
        --set ingress.hostname="$FQDN" \
        --values helm-wp-aks-values.yaml \
        wordpress bitnami/wordpress
    ```

Results:

<!-- expected_similarity=0.3 -->
```text
Release "wordpress" does not exist. Installing it now.
NAME: wordpress
LAST DEPLOYED: Tue Oct 24 16:19:35 2023
NAMESPACE: wordpress
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: wordpress
CHART VERSION: 18.0.8
APP VERSION: 6.3.2

** Please be patient while the chart is being deployed **

Your WordPress site can be accessed through the following DNS name from within your cluster:

    wordpress.wordpress.svc.cluster.local (port 80)

To access your WordPress site from outside the cluster follow the steps below:

1. Get the WordPress URL and associate WordPress hostname to your cluster external IP:

    export CLUSTER_IP=$(minikube ip) # On Minikube. Use: `kubectl cluster-info` on others K8s clusters
    echo "WordPress URL: https://mydnslabel1f2c0c.eastus.cloudapp.azure.com/"
    echo "$CLUSTER_IP  mydnslabel1f2c0c.eastus.cloudapp.azure.com" | sudo tee -a /etc/hosts

2. Open a browser and access WordPress using the obtained URL.

3. Login with the following credentials below to see your blog:

    echo Username: wpcliadmin
    echo Password: $(kubectl get secret --namespace wordpress wordpress -o jsonpath="{.data.wordpress-password}" | base64 -d)
```

## Browse your AKS Deployment Secured via HTTPS

Run the following command to get the HTTPS endpoint for your application:

> [!NOTE]
> It often takes 2-3 minutes for the SSL certificate to propogate and the site to be reachable via https.

```bash
runtime="5 minute"; endtime=$(date -ud "$runtime" +%s); while [[ $(date -u +%s) -le $endtime ]]; do STATUS=$(kubectl get svc --namespace=ingress-nginx ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}'); echo $STATUS; if [ "$STATUS" = "$MY_STATIC_IP" ]; then break; else sleep 10; fi; done

echo "You can now visit your web server at https://$FQDN"
```
