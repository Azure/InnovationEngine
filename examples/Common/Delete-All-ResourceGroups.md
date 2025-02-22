# Delete All Resource Groups Matching a Regular Expression

This executable document will delete all Azure resource groups that match a regular expression provided in an environment variable.

## Prerequisites

- [Azure CLI installed and Logged in](Prerequisite-AzureCLIAndSub.md)

## Setup the Environment

Set the environment variable `RESOURCE_GROUP_REGEX` with the desired regular expression. The example here will delete all resource groups with a name that ends with `DELETE_ME`

```bash
export RESOURCE_GROUP_REGEX=".*DELETE_ME$"
```

## Get All Resource Groups

Get a list of resource groups in the current subscription.

```bash
export resource_groups=$(az group list --query "[].name" -o tsv)
echo "Resource groups to be deleted: ${resource_groups}"
```

## Delete selected Resource Groups

Loop through all resource groups, deleting those that match the regular expression.

```bash
for rg in $resource_groups; do
    if [[ $rg =~ $RESOURCE_GROUP_REGEX ]]; then
    echo "Deleting resource group: $rg"
    az group delete --name $rg --yes --no-wait
    fi
done
```

This script will delete all resource groups that match the regular expression provided in the `RESOURCE_GROUP_REGEX` environment variable.