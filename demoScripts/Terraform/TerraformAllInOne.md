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

# Write the Terraform Configuration files

First we need a providers file that will define the required providers for our configuration:

```bash
mkdir $APPLICATION_NAME
cd $APPLICATION_NAME
cat <<EOF > providers.tf
# Configure the Azure provider
terraform {
  required_version = ">= 1.1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
  }
}

provider "azurerm" {
  features {}
}
EOF
```

All resources in Azure are placed in a resource group so here's a cofiguration to create the resources.

```bash
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

Assuming that things go well you will see an output similar to this.

<!-- expectedi_similarity=0.8 -->
```
Initializing the backend...

Initializing provider plugins...
- Finding hashicorp/azurerm versions matching "~> 3.0.2"...
- Installing hashicorp/azurerm v3.0.2...
- Installed hashicorp/azurerm v3.0.2 (signed by HashiCorp)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

# Format and validate the configuration

```bash
terraform fmt
terraform validate
```

It is essential that the Terraform config is correctly validate. This is indicated as follows:

<!-- expected_similrity=1.0 -->

```
Success! The configuration is valid.
```

# Apply your Terraform Configuration

Now it is time to apply the configuration. This will deploy the resources described in the configuration file.

```bash
terraform apply -auto-approve
```

Upon completion you will see an output similar to the below:

<!-- expected_similarity=0.7 -->
```
Terraform used the selected providers to generate the following execution plan.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # azurerm_resource_group.rg will be created
  + resource "azurerm_resource_group" "rg" {
      + id       = (known after apply)
      + location = "eastus"
      + name     = "RG_MAIN_Terraform_HelloWorld"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
azurerm_resource_group.rg: Creating...
azurerm_resource_group.rg: Creation complete after 2s [id=/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/RG_MAIN_Terraform_HelloWorld]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

# Inspect the state

At this point, assuming the previous step was succesful, you can be fairly sure that the resources were deployed. If you want to confirm then run:

```bash
terraform show
```

This will output something like:

<!-- expected_similarity=0.4 -->
```
# azurerm_resource_group.rg:
resource "azurerm_resource_group" "rg" {
    id       = "/subscriptions/325e7c34-99fb-4190-aa87-1df746c67705/resourceGroups/RG_MAIN_Terraform_HelloWorld"
    location = "eastus"
    name     = "RG_MAIN_Terraform_HelloWorld"
}
```

# How to create a VM using Terraform

TODO:

  * Write descriptive text for this content
  * Enable TerraformRgCreate as a prequisite
  * APP_ID is an empty environment variable

# Setup the Environment

The following are parameters that you are likely to want to configure for the VM:

```bash
export VM_NAME="TerraformCreatedVM"
export STORAGE_OS_DISK_NAME="OS_Disk"
export ADMIN_USER="admin"
export ADMIN_PASSWORD="pa$$w0rd"
```

For the network you will want to think about, at least, these parameters:

```bash
export IP_NAME=
export SUBNET_ID=
```

# Write the Terraform configuration file

We alwas need a Terraform configuration file

## Configure the Vm Resource

```bash
cat <<EOF >> vm.tf
# Define the VM resource
resource "azurerm_virtual_machine" "vm" {
  name                  = "$VM_NAME"
  location              = "$LOCATION"
  resource_group_name   = "$RESOURCE_GROUP_NAME"
  network_interface_ids = [azurerm_network_interface.vm.id]
  vm_size               = "Standard_B1s"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "20.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "$STORAGE_OS_DISK_NAME"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "$VM_NAME"
    admin_username = "$ADMIN_USER"
    admin_password = "$ADMIN_PASSWORD"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "production"
  }
}

EOF
```

##  Configure the Network Resource

```bash
cat <<EOF >> vm.tf
# Define the network interface resource
resource "azurerm_network_interface" "vm" {
  name                = "$VM_NIC_NAME"
  location            = azurerm_virtual_machine.vm.location
  resource_group_name = azurerm_virtual_machine.vm.resource_group_name

  ip_configuration {
    name                          = "$IP_NAME"
    subnet_id                     = "$SUBNET_ID"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vm.id
  }
}

EOF
```

# Configure the public IP address resource
```bash
cat <<EOF >> vm.tf
# Define the public IP address resource
resource "azurerm_public_ip" "vm" {
  name                = "my-vm-pip"
  location            = azurerm_virtual_machine.vm.location
  resource_group_name = azurerm_virtual_machine.vm.resource_group_name
  allocation_method   = "Dynamic"
  domain_name_label   = "my-vm-dns"
}

EOF
```

