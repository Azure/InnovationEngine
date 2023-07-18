## Requirements

- An **Azure Subscription** (e.g. [Free](https://aka.ms/azure-free-account) or [Student](https://aka.ms/azure-student-account) account)
- The [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- Bash shell (e.g. macOS, Linux, [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/about), [Multipass](https://multipass.run/), [Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/cloud-shell/quickstart), [GitHub Codespaces](https://github.com/features/codespaces), etc)

## 1. Clone Sample

```bash
git clone https://github.com/asw101/python-fastapi-pypy.git
```

Change into the directory of the cloned project
```bash
cd python-fastapi-pypy/
```

## 2. Install Azure CLI Extension and Register Resource Providers

If this is the first time you have used Azure Container Apps from the Azure CLI, or with your Azure Account, you will need to install the `containerapp` extension, and register the resource providers for `Microsoft.App` and `Microsoft.OperationalInsights` using the following commands.

```bash
az extension add --name containerapp
```

```bash
az provider register --namespace Microsoft.App --wait
```

```bash
az provider register --namespace Microsoft.OperationalInsights --wait
```

## 3. Setup Additional Environment Variables

```bash
SUBSCRIPTION_ID=$(az account show --query id --out tsv)
```

```bash
SCOPE="/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${MY_RESOURCE_GROUP_NAME}"
```

```bash
RANDOM_STR=$(echo -n "$SCOPE" | shasum | head -c 6)
```

```bash
ACR_NAME="acr${RANDOM_STR}"
```

## 4. Create Resource Group

```bash
az group create --name $MY_RESOURCE_GROUP_NAME --location $MY_LOCATION
```

## 5. Create Azure Container Registry

[Quickstart (docs.microsoft.com)](https://docs.microsoft.com/en-us/azure/container-registry/container-registry-get-started-azure-cli)

```bash
az acr create --resource-group $MY_RESOURCE_GROUP_NAME --name $ACR_NAME --sku Basic --admin-enabled true
```

```bash
az acr build -t $ACR_IMAGE_NAME -r $ACR_NAME .
```

```bash
CONTAINER_IMAGE="${ACR_NAME}.azurecr.io/${ACR_IMAGE_NAME}"
```

```bash
REGISTRY_SERVER="${ACR_NAME}.azurecr.io"
```

```bash
REGISTRY_USERNAME="${ACR_NAME}"
```

```bash
REGISTRY_PASSWORD=$(az acr credential show -n $ACR_NAME --query 'passwords[0].value' --out tsv)
```

```bash
echo "$CONTAINER_IMAGE"
```

## 6. Create Azure Container Apps Environment

[Quickstart (docs.microsoft.com)](https://docs.microsoft.com/en-us/azure/container-apps/get-started-existing-container-image?tabs=bash&pivots=container-apps-private-registry)

```bash
az containerapp env create --name $MY_CONTAINER_APPS_ENVIRONMENT --resource-group $MY_RESOURCE_GROUP_NAME --location $MY_LOCATION
```

## 7. Create Container App

```bash
az containerapp create --name my-container-app --resource-group $MY_RESOURCE_GROUP_NAME --environment $MY_CONTAINER_APPS_ENVIRONMENT --image "$CONTAINER_IMAGE" --registry-server "$REGISTRY_SERVER" --registry-username "$REGISTRY_USERNAME" --registry-password "$REGISTRY_PASSWORD" --target-port 80 --ingress 'external'
```

## 8. Test Container App with curl

```bash
CONTAINERAPP_FQDN=$(az containerapp show --resource-group $MY_RESOURCE_GROUP_NAME --name my-container-app --query properties.configuration.ingress.fqdn --out tsv)
```

```bash
echo "https://${CONTAINERAPP_FQDN}"
```

Results:

```expected_similarity=0.5
https://my-container-app.nicegrass-fdb78751.canadacentral.azurecontainerapps.io/
```

```bash
curl "https://${CONTAINERAPP_FQDN}/"
```

Congratulations! You have now created an Azure Container App for an existing Python Application!

## 9. Clean Up

Delete the Resource Group by running the following:

```bash
az group delete --name $MY_RESOURCE_GROUP_NAME
```

## Notes

- The sample in section 1 is originally from <https://github.com/tonybaloney/ants-azure-demos>, which is also referenced in step 1 of the video walkthrough. The updated sample is at: <https://github.com/asw101/python-fastapi-pypy>.
