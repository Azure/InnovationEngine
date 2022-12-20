This is just a hello world for now... it is based on https://developer.hashicorp.com/terraform/tutorials/azure-get-started/azure-build

# TODO Before publication

  * Descriptive content from that document should come here and we should replace that document with this one, when done.
  * Deploy a VM
  * Provide SSH connection details
  * Demonstrate logging onto the VM

# Prerequisites

Below are the prequisites, but this doesn't work in WSL so... You can use CloudShell.

  * Az CLI installed and logged in
  * [Terraform installed](https://askubuntu.com/questions/983351/how-to-install-terraform-in-ubuntu)
  * `az login` in WSL does not work, it opens in Lynx which doesn't have JS support. So add `--use-device-login` to enable you to manually authenticate via an external browser. Now  login is refused because it's not a managed device, only it is, it's my work laptop

# Setup the Environment

Grab the subscription ID:

```bash
export SUBSCRIPTION_ID=$(az account show --output tsv --query "id")
echo "Subscription ID in use: $SUBSCRIPTION_ID"

export TENANT_ID=$(az account show --output tsv --query "tenantId")
echo "Tenant ID in use: $TENANT_ID"
```

Set the application name you want to use:

```bash
export APPLICATION_NAME="Terraform_HelloWorld"
```

Set the location and resource group information:

```bash
export RESOURCE_GROUP_NAME="RG_MAIN_$APPLICATION_NAME"
export LOCATION="eastus"
```


Create a service principle for the application, grabbing the password from the output. This cannot be retrieved later, so this is important.:

```bash
export SP_PASSWORD=$(az ad sp create-for-rbac --role="Contributor" --name $APPLICATION_NAME --scopes="/subscriptions/$SUBSCRIPTION_ID" --output tsv --query "password")
```

Grab the appliction needed from the service principle:

```bash
az ad sp list --display-name "$APPLICATION_NAME"
```

# Write the Terraform Configuration file

```bash
mkdir $APPLICATION_NAME
cd $APPLICATION_NAME
cat <<EOF > main.tf
# Configure the Azure provider
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
  }

  required_version = ">= 1.1.0"
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "$RESOURCE_GROUP_NAME"
  location = "$LOCATION"
}
EOF
```

# Initialize your Terraform configuration

```bash
terraform init
```

# Format and validate the configuration

```bash
terraform fmt
terraform validate
```

# Apply your Terraform Configuration

```bash
terraform apply -auto-approve
```

# Inspect the state

```bash
terraform show
```
