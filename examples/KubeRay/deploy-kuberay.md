
# # Configure and deploy a Ray cluster on Azure Kubernetes Service (AKS)

In this article, you configure and deploy a Ray cluster on Azure Kubernetes Service (AKS) using KubeRay. You also learn how to use the Ray cluster to train a simple machine learning model and display the results on the Ray Dashboard.

## Prerequisites

* Review the [Ray cluster on AKS overview](./ray-overview.md) to understand the components and deployment process.

### Azure CLI

The Azure CLI is used to interact with Azure.

```bash
if ! command -v az &> /dev/null
then
  echo "Azure CLI could not be found, installing..."
  curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
fi

echo "Azure CLI is installed."
```

<!-- expected_similarity=".*installed" -->
```text
Azure CLI is installed.
```
For more details on installing the CLI see [How to install the Azure CLI](/cli/azure/install-azure-cli).


### Azure Subscription

You need to be logged in to an active Azure subscription is required. If you don't have an Azure subscription, you can create a free account [here](https://azure.microsoft.com/free/).

```bash
if ! az account show > /dev/null 2>&1; then
    echo "Please login to Azure CLI using 'az login' before running this script."
else
    export ACTIVE_SUBSCRIPTION_ID=$(az account show --query id -o tsv)
    echo "Currently logged in to Azure CLI. Using subscription ID: $ACTIVE_SUBSCRIPTION_ID."
fi
```

<!-- expected_similarity=0.8 -->
```text
Currently logged in to Azure CLI. Using subscription ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx
```

### Draft for Azure Kubernetes Service (AKS)

TODO: Is Draft really needed - not sure it is since I ran some tests sucessfully without it.

[Draft](https://github.com/Azure/draft) is an open-source project that streamlines Kubernetes development by taking a non-containerized application and generating the Dockerfiles, Kubernetes manifests, Helm charts, Kustomize configurations, and other artifacts associated with a containerized application. 

```bash
if ! command -v draft &> /dev/null
then
  echo "Draft could not be found, installing..."
  curl -fsSL https://raw.githubusercontent.com/Azure/draft/main/scripts/install.sh | bash
fi

echo "Draft is installed."
```

<!-- expected_similarity=".*installed" -->
```text
Draft is installed.
```

For more details on installing Draft see [Azure Kubernetes Service Preview extension](/azure/aks/draft#install-the-aks-preview-azure-cli-extension).

### Helm

Helm is a package manager for Kubernetes. It is one ofthe best ways to find, share, and use software built for Kubernetes.

```bash
if ! command -v helm &> /dev/null
then
  echo "Helm could not be found, installing..."
  curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
fi

echo "Helm is installed."
```

<!-- expected_similarity=".*installed" -->
```text
Helm is installed.
```

[Helm documentation](https://helm.sh/docs/intro/install/) provides more information on installing Helm.

### Terraform

[Terraform client tools](https://developer.hashicorp.com/terraform/install) or [OpenTofu](https://opentofu.org/) need to be installed. This article uses Terraform, but the modules used should be compatible with OpenTofu.

```bash
if ! command -v terraform &> /dev/null
then
  echo "Terraform could not be found, installing..."
  curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
  sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
  sudo apt-get update && sudo apt-get install terraform
fi

echo "Terraform is installed."
```

<!-- expected_similarity=".*installed" -->
```text
Terraform is installed.
```

## Create an AKS cluster 

We will create a new AKS cluster and deploy Ray onto it. To do this we will use Terraform. The Terraform files are provided in a sample project on GitHub, clone or update it for local use using:

```bash
if [ -d "aks-ray-sample" ]; then
    cd aks-ray-sample
    git pull
else
    git clone https://github.com/Azure-Samples/aks-ray-sample
    cd aks-ray-sample
fi
pwd
```

<!-- expected_similarity=".*/aks-ray-sample" -->
```
/home/foo/projects/aks-ray-sample
```

### Setup the environment

For convenience we will define some environment variables that Terraform uses. 

```bash
export TF_VAR_resource_group_owner="TODO_DefineOwnerEnvironmentVariable"
export TF_VAR_subscription_id=$ACTIVE_SUBSCRIPTION_ID
```

Note that the value for `ACTIVE_SUBSCRIPTION_ID` was set in the prerequisites section using `export ACTIVE_SUBSCRIPTION_ID=$(az account show --query id -o tsv)`.

### Create the Terraform plan

As is always the case with Terraform the first thing we need to do is create the Terraform plan.

TODO: the requirement of the cd here is a limitation of current version of IE. This needs to be resolved.

```bash
cd aks-ray-sample
terraform init
terraform plan -out main.tfplan
```

The output of this command will tell you about the changes that Terraform believes are necessary. This section, if succesful, will end with:

<!-- expected__similarity=".*Saved the plan to: main.tfplan.*" -->
```text
Saved the plan to: main.tfplan

To perform exactly these actions, run the following command to apply:
    terraform apply "main.tfplan"
```

### Apply the Terraform plan

Now that we have the Terraform plan we can go ahead and apply it. This should create our AKS cluster for us.

TODO: the requirement of the cd here is a limitation of current version of IE. This needs to be resolved.

```bash
cd aks-ray-sample
terraform apply "main.tfplan"
```

The output of this command should show details of the deployment and include the following text:

<!-- expected_similarity=".*Apply complete!.* -->
```
[Details of the Deploymet]

Apply complete!

[Summary of outputs from the command]
```

## Deploy the KubeRay onto the AKS cluster

Now that we have an AKS cluster we can deploy KubeRay to it.

### Configure the environment

It is good practice to store values Retrieve the Terraform outputs and store in variables for later use:

```bash
cd aks-ray-sample

export kuberay_namespace="kuberay"
echo "The kuberay namespace is `$kuberay_namespace`"

export resource_group_name=$(terraform output -raw resource_group_name)
echo "Resource group name is '$resource_group_name'"

export system_node_pool_name=$(terraform output -raw system_node_pool_name)
echo "System Node Pool name is '$system_node_pool_name'"

export aks_cluster_name=$(terraform output -raw kubernetes_cluster_name)
echo "AKS Cluster name is '$aks_cluster_name'"
```

<!-- expected_similarity=0.5 -->
```text
The kuberay namespace is `kuberay`
Resource group name is 'ExampleResourceGroupName'
System Node Pool name is 'ExampleNodePollName'
AKS Cluster name is 'ExampleClusterName'
```

### Get AKS credentials for the cluster

In order to configure the environment to work with the AKS cluster we need to grab the credentials for it. The Azure CLI provides a way to do this.

```bash
az aks get-credentials \
    --resource-group $resource_group_name \
    --name $aks_cluster_name \
    --overwrite-existing
```

### Create the kuberay namespace

We will want KubeRay to be in its own namespace.

```bash
if ! kubectl get namespace $kuberay_namespace; then
    kubectl create namespace $kuberay_namespace
else
    echo "Namespace $kuberay_namespace already exists"
fi
```

### Add the KubeRay Helm repository

We will deploy KubeRay using Helm, so we need to ensure that the KubeRay repository is available.

```bash
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo update
```

### Install or upgrade the KubeRay operator using Helm

We are now ready to install Kuberay operator. This command will upgrade any existing installation, or install a new instance as necessary.

```bash
helm upgrade \
  --install \
  --cleanup-on-fail \
  --wait \
  --timeout 10m0s \
  --namespace "$kuberay_namespace" \
  --create-namespace kuberay-operator kuberay/kuberay-operator \
  --version 1.1.1 \
  --set clusterRoles.namespace="$kuberay_namespace"
```

### Deploy the RayJob specification

Now that we have an AKS cluster with Kuberay on it we can get to work and train our model.

Fashion MNIST is a dataset of Zalando's article images consisting of a training set of 60,000 examples and a test set of 10,000 examples. Each example is a 28x28 grayscale image associated with a label from ten classes. We will use this dataset to train a simple PyTorch modelusing the Ray cluster we have created.

To use this dataset we will leverate a [Ray Job specification](https://github.com/ray-project/kuberay/blob/master/ray-operator/config/samples/pytorch-mnist/ray-job.pytorch-mnist.yaml). This is a YAML file that describes the resources required to run the job, including the Docker image, the command to run, and the number of workers to use. 

### Download the PyTorch MNIST job YAML file

Looking at the Ray Job specification, you might need to modify some fields to match your environment, but you can proceed with the default setup if you simply want to get things working.

* The `replicas` field under the `workerGroupSpecs` section in `rayClusterSpec` specifies the number of worker pods that KubeRay schedules to the Kubernetes cluster. Each worker pod requires *3 CPUs* and *4 GB of memory*. The head pod requires *1 CPU* and *4 GB of memory*. Setting the `replicas` field to *2* requires *8 vCPUs* in the node pool used to implement the RayCluster for the job.
* The `NUM_WORKERS` field under `runtimeEnvYAML` in `spec` specifies the number of Ray actors to launch. Each Ray actor must be serviced by a worker pod in the Kubernetes cluster, so this field must be less than or equal to the `replicas` field. In this example, we set `NUM_WORKERS` to *2*, which matches the `replicas` field.
* The `CPUS_PER_WORKER` field must be set to *less than or equal the number of CPUs allocated to each worker pod minus 1*. In this example, the CPU resource request per worker pod is *3*, so `CPUS_PER_WORKER` is set to *2*.

To summarize, you need a total of *8 vCPUs* in the node pool to run the PyTorch model training job. Since we added a taint on the system node pool so that no user pods can be scheduled on it, we must create a new node pool with at least *8 vCPUs* to host the Ray cluster.

First lets download the specification.

```bash
curl -LO https://raw.githubusercontent.com/ray-project/kuberay/master/ray-operator/config/samples/pytorch-mnist/ray-job.pytorch-mnist.yaml
```

If you want to make any changes, you can do so in your local copy.

### Train a PyTorch Model on Fashion MNIST

We can submit this specification to our KubeRay cluster as follows.

```bash
kubectl apply -n $kuberay_namespace -f ray-job.pytorch-mnist.yaml
```

This can take some time to run, you can monitor progress with the following command.

```bash
job_status=$(kubectl get rayjobs -n $kuberay_namespace -o jsonpath='{.items[0].status.jobDeploymentStatus}')
echo "KubeRay job status is '$job_status'"
```

Completion is indicated by the value of job_status being `Complete`. We can cause a script to wait for completion in a while loop. 

```bash
while [ "$job_status" != "Complete" ]; do
    echo -ne "Job Status: $job_status\\r"
    sleep 30
    job_status=$(kubectl get rayjobs -n $kuberay_namespace -o jsonpath='{.items[0].status.jobDeploymentStatus}')
done
echo "Job Status: $job_status"
```

Once the job has completed we can check to ensure that it succeeded and abort the script of there was a problem.

```bash
job_status=$(kubectl get rayjobs -n $kuberay_namespace -o jsonpath='{.items[0].status.jobStatus}')

if [ "$job_status" != "SUCCEEDED" ]; then
    echo "Job Failed!"
    exit 1
else
  echo "Job succeeded."
fi
```

<!-- expected_similarity=1.0 -->
```text
Job succeeded.
```

## Provide access to the KubeRay dashboard via a browser

When the RayJob successfully completes, you can view the training results on the Ray Dashboard. The Ray Dashboard provides real-time monitoring and visualizations of Ray clusters. You can use the Ray Dashboard to monitor the status of Ray clusters, view logs, and visualize the results of machine learning jobs.

To access the Ray Dashboard, you need to expose the Ray head service to the public internet by creating a *service shim* to expose the Ray head service on port 80 instead of port 8265.

TODO: this says "real-time monitoring" wouldn't it make more sense to do this before submitting the job? That way the user can view the dashboard while the job is running, is this useful?

Get the name of the Ray head service and save it in a shell variable using.

```bash
export rayclusterhead=$(kubectl get service -n $kuberay_namespace | grep 'rayjob-pytorch-mnist-raycluster' | grep 'ClusterIP' | awk '{print $1}')
echo "Ray Cluster Head is '$rayclusterhead'"
```

Create the service shim to expose the Ray head service on port 80 using the `kubectl expose service` command.

```bash
if ! kubectl get service ray-dash -n $kuberay_namespace; then
  kubectl expose service $rayclusterhead \
    -n $kuberay_namespace \
    --port=80 \
    --target-port=8265 \
    --type=NodePort \
    --name=ray-dash
else
  echo "Service ray-dash already exists"
fi
```

Create the ingress rule to expose the service shim using the ingress controller using the following command:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ray-dash
  namespace: kuberay
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: webapprouting.kubernetes.azure.com
  rules:
  - http:
      paths:
      - backend:
          service:
            name: ray-dash
            port:
              number: 80
        path: /
        pathType: Prefix
EOF
```

Finally we can find the find the public IP address of the ingress controller and output it for the user.

```bash
export lb_public_ip=$(kubectl get svc -n app-routing-system -o jsonpath='{.items[?(@.metadata.name == "nginx")].status.loadBalancer.ingress[0].ip}')

echo "KubeRay Dashboard URL: http://$lb_public_ip/"
```

You can now view the dashboard using the IP displayed.

## Clean up resources

To clean up the resources created in this guide, you can delete the Azure resource group that contains the AKS cluster. As a reminder the group name is stored in an environment variable for easy access:

```bash
echo "Resource group name is '$resource_group_name'"
```

## Debugging

It's always good to have a way to inspect the status of a failed deployment. The following commands will help you inspect the status of your KubeRay Job should it fail.

TODO: the original doc had the following commands that might be good to move here, but in the context of an automated script didn't do anything useful.

`kubectl get pods -n kuberay`
`kubectl logs -n kuberay rayjob-pytorch-mnist-fc959`

## Next steps

To learn more about AI and machine learning workloads on AKS, see the following articles:

* [Deploy an application that uses OpenAI on Azure Kubernetes Service (AKS)](./open-ai-quickstart.md)
* [Build and deploy data and machine learning pipelines with Flyte on Azure Kubernetes Service (AKS)](./use-flyte.md)
* [Deploy an AI model on Azure Kubernetes Service (AKS) with the AI toolchain operator (preview)](./ai-toolchain-operator.md)

## Contributors

*Microsoft maintains this article. The following contributors originally wrote it:*

* Russell de Pina | Principal TPM
* Ken Kilty | Principal TPM
* Erin Schaffer | Content Developer 2