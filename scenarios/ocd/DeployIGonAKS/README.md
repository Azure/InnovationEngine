# Quickstart: Deploy Inspektor Gadget in an Azure Kubernetes Service cluster

Welcome to this tutorial where we will take you step by step in deploying [Inspektor Gadget](https://www.inspektor-gadget.io/) in an Azure Kubernetes Service (AKS) cluster with the kubectl plugin: `gadget`. This tutorial assumes you are logged into Azure CLI already and have selected a subscription to use with the CLI.

## Define Environment Variables

The First step in this tutorial is to define environment variables:

```bash
export RANDOM_ID="$(openssl rand -hex 3)"
export MY_RESOURCE_GROUP_NAME="myResourceGroup$RANDOM_ID"
export REGION="eastus"
export MY_AKS_CLUSTER_NAME="myAKSCluster$RANDOM_ID"
```

## Create a resource group

A resource group is a container for related resources. All resources must be placed in a resource group. We will create one for this tutorial. The following command creates a resource group with the previously defined $MY_RESOURCE_GROUP_NAME and $REGION parameters.

```bash
az group create --name $MY_RESOURCE_GROUP_NAME --location $REGION
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
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

## Create AKS Cluster

Create an AKS cluster using the az aks create command.

This will take a few minutes.

```bash
az aks create \
  --resource-group $MY_RESOURCE_GROUP_NAME \
  --name $MY_AKS_CLUSTER_NAME \
  --location $REGION \
  --no-ssh-key
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

## Install Inspektor Gadget

The Inspektor Gadget installation is composed of two steps:

1. Installing the kubectl plugin in the user's system.
2. Installing Inspektor Gadget in the cluster.

    > [!NOTE]
    > There are additional mechanisms for deploying and utilizing Inspektor Gadget, each tailored to specific use cases and requirements. Using the `kubectl gadget` plugin covers many of them, but not all. For instance, deploying Inspektor Gadget with the `kubectl gadget` plugin depends on the Kubernetes API server's availability. So, if you can’t depend on such a component because its availability could be sometimes compromised, then it is recommended to not use the `kubectl gadget`deployment mechanism. Please check [ig documentation](https://github.com/inspektor-gadget/inspektor-gadget/blob/main/docs/ig.md) to know what to do in that, and other use cases.

### Installing the kubectl plugin: `gadget`

Install the latest version of the kubectl plugin from the releases page, uncompress and move the `kubectl-gadget` executable to `$HOME/.local/bin`:

> [!NOTE]
> If you want to install it using [`krew`](https://sigs.k8s.io/krew) or compile it from the source, please follow the official documentation: [installing kubectl gadget](https://github.com/inspektor-gadget/inspektor-gadget/blob/main/docs/install.md#installing-kubectl-gadget).

```bash
IG_VERSION=$(curl -s https://api.github.com/repos/inspektor-gadget/inspektor-gadget/releases/latest | jq -r .tag_name)
IG_ARCH=amd64
mkdir -p $HOME/.local/bin
export PATH=$PATH:$HOME/.local/bin
curl -sL https://github.com/inspektor-gadget/inspektor-gadget/releases/download/${IG_VERSION}/kubectl-gadget-linux-${IG_ARCH}-${IG_VERSION}.tar.gz  | tar -C $HOME/.local/bin -xzf - kubectl-gadget
```

Now, let’s verify the installation by running the `version` command:

```bash
kubectl gadget version
```

The `version` command will display the version of the client (kubectl gadget plugin) and show that it is not yet installed in the server (the cluster):

<!--expected_similarity="(?m)^Client version: v\d+\.\d+\.\d+$\n^Server version: not installed$"-->
```text
Client version: vX.Y.Z
Server version: not installed
```

### Installing Inspektor Gadget in the cluster

The following command will deploy the DaemonSet:

> [!NOTE]
> Several options are available to customize the deployment: use a specific container image, deploy to specific nodes, and many others. To know all of them, please check the official documentation: [installing in the cluster](https://github.com/inspektor-gadget/inspektor-gadget/blob/main/docs/install.md#installing-in-the-cluster).

```bash
kubectl gadget deploy
```

Now, let’s verify the installation by running the `version` command again:

```bash
kubectl gadget version
```

This time, the client and server will be correctly installed:

<!--expected_similarity="(?m)^Client version: v\d+\.\d+\.\d+$\n^Server version: v\d+\.\d+\.\d+$"-->
```text
Client version: vX.Y.Z
Server version: vX.Y.Z
```

You can now start running the gadgets:

```bash
kubectl gadget help
```

<!--
## Clean Up

### Undeploy Inspektor Gadget

```bash
kubectl gadget undeploy
```

### Clean up Azure resources

When no longer needed, you can use `az group delete` to remove the resource group, cluster, and all related resources as follows. The `--no-wait` parameter returns control to the prompt without waiting for the operation to complete. The `--yes` parameter confirms that you wish to delete the resources without an additional prompt to do so.

```bash
az group delete --name $MY_RESOURCE_GROUP_NAME --no-wait --yes
```
-->