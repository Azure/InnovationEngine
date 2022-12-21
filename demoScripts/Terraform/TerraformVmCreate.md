# How to create a VM using Terraform

TODO:

  * Write descriptive text for this content
  * Enable TerraformRgCreate as a prequisite
  
# Write the Terraform configuration file

We alwas need a Terraform configuration file

## Configure the provider

```bash
mkdir $APPLICATION_NAME
cd $APPLICATION_NAME
cat <<EOF > vm_create.tf
# Configure the provider
provider "azurerm" {
  version = "2.33.0"
  tenant_id = "$TENANT_ID"
  subscription_id = "$SUBSCRIPTION_ID"
  client_id = "$APP_ID"
  client_secret = "$SP_PASSWORD"
}
EOF
```

## Configure the Vm Resource

```bash
mkdir $APPLICATION_NAME
cd $APPLICATION_NAME
cat <<EOF >> vm_create.tf
# # Define the VM resource
resource "azurerm_virtual_machine" "vm" {
  name                  = "my-vm"
  location              = "<your-location>"
  resource_group_name   = "<your-resource-group>"
  network_interface_ids = [azurerm_network_interface.vm.id]
  vm_size               = "Standard_B1s"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "20.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "my-vm-osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "my-vm"
    admin_username = "user"
    admin_password = "<your-admin-password>"
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
mkdir $APPLICATION_NAME
cd $APPLICATION_NAME
cat <<EOF >> vm_create.tf
# Define the network interface resource
resource "azurerm_network_interface" "vm" {
  name                = "my-vm-nic"
  location            = azurerm_virtual_machine.vm.location
  resource_group_name = azurerm_virtual_machine.vm.resource_group_name

  ip_configuration {
    name                          = "my-vm-ipconfig"
    subnet_id                     = "<your-subnet-id>"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vm.id
  }
}
EOF
```

# Configure the public IP address resource
```bash
mkdir $APPLICATION_NAME
cd $APPLICATION_NAME
cat <<EOF >> vm_create.tf
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

