# AKS Store Demo

TODO: Move intro content from source document

## Setup The Environment

This document, and those listed in the prerequisites section use environment variables to make reuse easier. Defaults are provided in most documents, but for consistency and clarity we will define the ones we really care about here:

```bash
export RANDOM_ID="$(openssl rand -hex 3)"
export MY_RESOURCE_GROUP_NAME="aks-store-demo-ResourceGroup-$RANDOM_ID"
export REGION="westus2"
export MY_AKS_CLUSTER_NAME="aks-store-demo-$RANDOM_ID"
export MY_DNS_LABEL="aks-store-dns-label-$RANDOM_ID"
```

## Prerequisites

  * Have an [active Azure Subscription (free subscriptions available) and an install of Azure CLI](../../Common/Prerequisite-AzureCLIAndSub.md)
  * Ensure the [`az aks`](../../Common/Prerequisite-AzureCLI-ls AKS.md) commands are installed
  * Install [Helm](../../Common/Prerequisites-Helm.md) - package manager for Kubernetes.
  * Install [Terraform](../../Common/Prerequisites-Terraform.md) - Infrastructure as Code management tool
  * Have an existing [AKS Cluster](https://raw.githubusercontent.com/MicrosoftDocs/azure-aks-docs/refs/heads/main/articles/aks/learn/quick-kubernetes-deploy-cli.md)

## Create the custom-values.yaml file

```bash
cat << EOF > custom-values.yaml
namespace: ${AZURE_AKS_NAMESPACE}
EOF
```

## Add Azure Managed Identity and set to use AzureAD auth 

```bash
if [ -n "${AZURE_IDENTITY_CLIENT_ID}" ] && [ -n "${AZURE_IDENTITY_NAME}" ]; then
  cat << EOF >> custom-values.yaml
useAzureAd: true
managedIdentityName: ${AZURE_IDENTITY_NAME}
managedIdentityClientId: ${AZURE_IDENTITY_CLIENT_ID}
EOF
fi
```

## Add base images

```bash
cat << EOF >> custom-values.yaml
namespace: ${AZURE_AKS_NAMESPACE}
productService:
  image:
    repository: ${AZURE_REGISTRY_URI}/aks-store-demo/product-service
storeAdmin:
  image:
    repository: ${AZURE_REGISTRY_URI}/aks-store-demo/store-admin
storeFront:
  image:
    repository: ${AZURE_REGISTRY_URI}/aks-store-demo/store-front
virtualCustomer:
  image:
    repository: ${AZURE_REGISTRY_URI}/aks-store-demo/virtual-customer
virtualWorker:
  image:
    repository: ${AZURE_REGISTRY_URI}/aks-store-demo/virtual-worker
EOF
```

## Add ai-service if Azure OpenAI endpoint is provided

```bash
if [ -n "${AZURE_OPENAI_ENDPOINT}" ]; then
  cat << EOF >> custom-values.yaml
aiService:
  image:
      repository: ${AZURE_REGISTRY_URI}/aks-store-demo/ai-service
  create: true
  modelDeploymentName: ${AZURE_OPENAI_MODEL_NAME}
  openAiEndpoint: ${AZURE_OPENAI_ENDPOINT}
  useAzureOpenAi: true
EOF

  # If Azure identity does not exists, use the Azure OpenAI API key
  if [ -z "${AZURE_IDENTITY_CLIENT_ID}" ] && [ -z "${AZURE_IDENTITY_NAME}" ]; then
    cat << EOF >> custom-values.yaml
  openAiKey: $(az keyvault secret show --name ${AZURE_OPENAI_KEY} --vault-name ${AZURE_KEY_VAULT_NAME} --query value -o tsv)
EOF
  fi

  # If DALL-E model endpoint and name exists
  if [ -n "${AZURE_OPENAI_DALL_E_ENDPOINT}" ] && [ -n "${AZURE_OPENAI_DALL_E_MODEL_NAME}" ]; then
    cat << EOF >> custom-values.yaml
  openAiDalleEndpoint: ${AZURE_OPENAI_DALL_E_ENDPOINT}
  openAiDalleModelName: ${AZURE_OPENAI_DALL_E_MODEL_NAME}
EOF
  fi
fi
```

## Add order-service
```bash
cat << EOF >> custom-values.yaml
orderService:
  image:
    repository: ${AZURE_REGISTRY_URI}/aks-store-demo/order-service
EOF
```

## Add Azure Service Bus to order-service if provided
```bash
  if [ -n "${AZURE_SERVICE_BUS_HOST}" ]; then
    cat << EOF >> custom-values.yaml
  queueHost: ${AZURE_SERVICE_BUS_HOST}
EOF

  # If Azure identity does not exists, use the Azure Service Bus credentials
  if [ -z "${AZURE_IDENTITY_CLIENT_ID}" ] && [ -z "${AZURE_IDENTITY_NAME}" ]; then
    cat << EOF >> custom-values.yaml
  queuePort: "5671"
  queueTransport: "tls"
  queueUsername: ${AZURE_SERVICE_BUS_SENDER_NAME}
  queuePassword: $(az keyvault secret show --name ${AZURE_SERVICE_BUS_SENDER_KEY} --vault-name ${AZURE_KEY_VAULT_NAME} --query value -o tsv)
EOF
  fi
fi
```

## Add makeline-service

```bash
cat << EOF >> custom-values.yaml
makelineService:
  image:
    repository: ${AZURE_REGISTRY_URI}/aks-store-demo/makeline-service
EOF
```

# Add Azure Service Bus to makeline-service if provided
```bash
if [ -n "${AZURE_SERVICE_BUS_URI}" ]; then
  # If Azure identity exists just set the Azure Service Bus Hostname
  if [ -n "${AZURE_IDENTITY_CLIENT_ID}" ] && [ -n "${AZURE_IDENTITY_NAME}" ]; then
    cat << EOF >> custom-values.yaml
    orderQueueHost: ${AZURE_SERVICE_BUS_HOST}
EOF
  else
    cat << EOF >> custom-values.yaml
  orderQueueUri: ${AZURE_SERVICE_BUS_URI}
  orderQueueUsername: ${AZURE_SERVICE_BUS_LISTENER_NAME}
  orderQueuePassword: $(az keyvault secret show --name ${AZURE_SERVICE_BUS_LISTENER_KEY} --vault-name ${AZURE_KEY_VAULT_NAME} --query value -o tsv)
EOF
  fi
fi
```

## Add Azure Cosmos DB to makeline-service if provided
```bash
if [ -n "${AZURE_COSMOS_DATABASE_URI}" ]; then
  cat << EOF >> custom-values.yaml
  orderDBApi: ${AZURE_DATABASE_API}
  orderDBUri: ${AZURE_COSMOS_DATABASE_URI}
EOF
  # If Azure identity does not exists, use the Azure Cosmos DB credentials
  if [ -z "${AZURE_IDENTITY_CLIENT_ID}" ] && [ -z "${AZURE_IDENTITY_NAME}" ]; then
    cat << EOF >> custom-values.yaml
  orderDBUsername: ${AZURE_COSMOS_DATABASE_NAME}
  orderDBPassword: $(az keyvault secret show --name ${AZURE_COSMOS_DATABASE_KEY} --vault-name ${AZURE_KEY_VAULT_NAME} --query value -o tsv)
EOF
  fi
fi
```

## Do not deploy RabbitMQ when using Azure Service Bus

```bash
if [ -n "${AZURE_SERVICE_BUS_HOST}" ]; then
  cat << EOF >> custom-values.yaml
useRabbitMQ: false
EOF
fi
```

## Do not deploy MongoDB when using Azure Cosmos DB
```bash
if [ -n "${AZURE_COSMOS_DATABASE_URI}" ]; then
  cat << EOF >> custom-values.yaml
useMongoDB: false
EOF
fi
```