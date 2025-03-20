Before starting to deploy a VM it is a good idea to check that availability exists in the region desired. This document explains how to do that.

## Configure the Environment

We use enviroment variables to simplify commands, some of them will have been set in the above prerequisites, and echoed below for convenience. The remaining ones are set with defaults:

```bash
export LOCATION=eastus
export VM_SKU=Standard_D2als_v6
# export VM_SKU=Standard_L8s # this is an invalid VM_SKU for most users deliberately selected to create a failure in validation
```

## Prerequisites

The VM SKU chosen in the previous section is one that is not usually available to customers. This is deliberate so that we can demonstrate using prerequisites to valdidate the options chosen.

* Have an [active Azure Subscription (free subscriptions available) and an install of Azure CLI](../Common/Prerequisite-AzureCLIAndSub.md)
* Have a [valid configuration](../Common/Prerequisite-Validation.md) in the selected region for the requested VM SKU size
