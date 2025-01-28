Before starting to deploy a VM it is a good idea to check that availability exists in the region desired. This document explains how to do that.

# Prerequisites

* Have an [active Azure Subscription (free subscriptions available) and an install of Azure CLI](../Common/Prerequisites-AzureCLIAndSub.md)

# Configure the Environment

We use enviroment variables to simplify commands, some of them will have been set in the above prerequisites, and echoed below for convenience. The remaining ones are set with defaults:

```bash
echo "ACTIVE_SUBSCRIPTION_ID=$ACTIVE_SUBSCRIPTION_ID"
export AZURE_LOCATION=eastus
export VM_SKU=Standard_D2_v2
```

# Check Availability

We can use the az CLI to check availability of the desired SKU in the location selected with the currently active subscription as follows:

```bash
az vm list-skus --location $AZURE_LOCATION --subscription $ACTIVE_SUBSCRIPTION_ID --size $VM_SKU --output table
```
