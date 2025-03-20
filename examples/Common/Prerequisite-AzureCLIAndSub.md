This document uses the Azure CLI connected to an active Azure Subscription. The following commands ensure that you have both an active subscription and a current version of the Azure CLI. Assuming you are logged in and have executed these commands the environment variable `ACTIVE_SUBSCRIPTION_ID` will contain the currently active subscription ID.

### Azure CLI

The Azure CLI is used to interact with Azure.

```bash
if ! command -v az &> /dev/null
then
  echo "Azure CLI could not be found, installing..."
  curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
else
  echo "Azure CLI is installed."
fi
```

<!-- expected_similarity=".*installed" -->
```text
Azure CLI is installed.
```
For more details on installing the CLI see [How to install the Azure CLI](/cli/azure/install-azure-cli).


### Azure Subscription

You need to be logged in to an active Azure subscription is required. If you don't have an Azure subscription, you can [create a free account](https://azure.microsoft.com/free/).


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

Once logged in we need to ensure that there is an active refresh token.

```bash
if ! az account get-access-token > /dev/null 2>&1; then
  echo "Azure CLI session has expired. Please login with 'az login --use-device-code' and try again."
else
  echo "Azure CLI session is active."
fi
```

<!-- expected_similarity=0.8 -->
```text
Azure CLI session is active.
```

### Azure Tenant ID

Retrieve the tenant ID associated with the active Azure subscription and store it in an environment variable called `TENANT_ID`.

```bash
export TENANT_ID=$(az account show --query tenantId -o tsv)
echo "Tenant ID: $TENANT_ID"
```

<!-- expected_similarity=0.4 -->
```text
Tenant ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx
```