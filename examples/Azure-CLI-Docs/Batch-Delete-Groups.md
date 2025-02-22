<!-- TODO: Convert the better version of this content found at https://learn.microsoft.com/en-us/cli/azure/delete-azure-resources-at-scale -->

# Batch Delete Resource Groups

When working with new infrastructure configurations it is common to have a number of unused resource groups left behind. If these have resources in them then you will be spending money needlessly. It can therefore be useful to automates the deletion of these groups.

## Environment Setup

It is a good practice is to use a common prefix for resource groups within a particular work unit. This allows us to query the list of resouces on Azure. So lets create a variable for that prefix.

```bash
export COMMON_PREFIX="Tutorial_Content"
```

## Create some dummy Resource Groups

In order to demonstrate this method we need to create some dummy resource groups. We can use the following code block. This script will create three resource groups using the `COMMON_PREFIX` for the name.

```bash
for i in 1 2 3; do
    az group create --name "${COMMON_PREFIX}_RG_$i" --location "eastus"
done
```

## Getting the list of Groups to Delete

We can now query the resource groups in the subscription, filter on our prefix and store the result in an environment variable, these are the candidates for deletion:

```bash
export RG_TO_DELETE=$(az group list --query "[?starts_with(name, '$COMMON_PREFIX')].name" -o tsv | tr '\n' ',' | sed 's/,$//')
echo $RG_TO_DELETE
```

## Deleting the Resource Groups

Now that we have identified the resource groups to delete, we can proceed with the deletion process. The following script will iterate over each resource group name stored in the `RGS_TO_DELETE` variable and delete them one by one.

```bash
IFS=',' read -ra RG_ARRAY <<< "$RG_TO_DELETE"
for rg in "${RG_ARRAY[@]}"; do
    az group delete --name $rg --yes --no-wait
    echo "$rg is being deleted"
done
```

This script uses the `--yes` flag to confirm the deletion without prompting and the `--no-wait` flag to return immediately without waiting for the operation to complete. This allows the script to proceed with deleting the next resource group without delay.