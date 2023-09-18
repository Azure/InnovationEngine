# Intro to Create a NGINX Webserver Secured via HTTPS

To secure web servers, a Transport Layer Security (TLS), previously known as Secure Sockets Layer (SSL), certificate can be used to encrypt web traffic. These TLS/SSL certificates can be stored in Azure Key Vault, and allow secure deployments of certificates to Linux virtual machines (VMs) in Azure. In this tutorial you learn how to:

> [!div class="checklist"]

> * Setup and secure Azure Networking
> * Create an Azure Key Vault
> * Generate or upload a certificate to the Key Vault
> * Create a VM and install the NGINX web server
> * Inject the certificate into the VM and configure NGINX with a TLS binding

If you choose to install and use the CLI locally, this tutorial requires that you're running the Azure CLI version 2.0.30 or later. Run `az --version` to find the version. If you need to install or upgrade, see [Install Azure CLI]( https://learn.microsoft.com//cli/azure/install-azure-cli ).

## Variable Declaration

List of all the environment variables you'll need to execute this tutorial:

```bash
export UNIQUE_POSTFIX="$(($RANDOM % 254 + 1))"
export MY_RESOURCE_GROUP_NAME="myResourceGroup$UNIQUE_POSTFIX"
export MY_KEY_VAULT="myKeyVault$UNIQUE_POSTFIX"
export MY_LOCATION="eastus"
export MY_VM_NAME="myVMName$UNIQUE_POSTFIX"
export MY_VM_IMAGE='Canonical:0001-com-ubuntu-minimal-jammy:minimal-22_04-lts-gen2:latest'
export MY_VM_USERNAME="azureadmin"
export MY_VM_SIZE='Standard_DS2_v2'
export MY_VNET_NAME="myVNet$UNIQUE_POSTFIX"
export MY_VNET_PREFIX="10.$UNIQUE_POSTFIX.0.0/16"
export MY_SN_NAME="mySN$UNIQUE_POSTFIX"
export MY_SN_PREFIX="10.$UNIQUE_POSTFIX.0.0/24"
export MY_PUBLIC_IP_NAME="myPublicIP$UNIQUE_POSTFIX"
export MY_DNS_LABEL="mydnslabel$UNIQUE_POSTFIX"
export MY_NSG_NAME="myNSGName$UNIQUE_POSTFIX"
```

## Create a Resource Group

Before you can create a secure Linux VM, create a resource group with az group create. The following example creates a resource group named *myResourceGroup$UNIQUE_POSTFIX* in the *eastus* location:

```bash
az group create \
    --name $MY_RESOURCE_GROUP_NAME \
    --location $MY_LOCATION
```

Results:

```JSON
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242",
  "location": "eastus",
  "managedBy": null,
  "name": "myResourceGroup242",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```
## Set up VM Network

Use az network vnet create to create a virtual network named *$MY_VNET_NAME* with a subnet named *$MY_SN_NAME*in the *$MY_RESOURCE_GROUP_NAME*resource group.

```bash
az network vnet create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_VNET_NAME \
    --location $MY_LOCATION \
    --address-prefix $MY_VNET_PREFIX \
    --subnet-name $MY_SN_NAME \
    --subnet-prefix $MY_SN_PREFIX
```

Results:

```JSON
{
  "newVNet": {
    "addressSpace": {
      "addressPrefixes": [
        "10.242.0.0/16"
      ]
    },
    "enableDdosProtection": false,
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/virtualNetworks/myVNet242",
    "location": "eastus",
    "name": "myVNet242",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup242",
    "subnets": [
      {
        "addressPrefix": "10.242.0.0/24",
        "delegations": [],
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/virtualNetworks/myVNet242/subnets/mySN242",
        "name": "mySN242",
        "privateEndpointNetworkPolicies": "Disabled",
        "privateLinkServiceNetworkPolicies": "Enabled",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup242",
        "type": "Microsoft.Network/virtualNetworks/subnets"
      }
    ],
    "type": "Microsoft.Network/virtualNetworks",
    "virtualNetworkPeerings": []
  }
}
```

Use az network public-ip create to create a standard zone-redundant public IPv4 address named *$MY_PUBLIC_IP_NAME* in *$MY_RESOURCE_GROUP_NAME*.

```bash
az network public-ip create \
    --name $MY_PUBLIC_IP_NAME \
    --location $MY_LOCATION \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --dns-name $MY_DNS_LABEL \
    --sku Standard \
    --allocation-method static \
    --version IPv4 \
    --zone 1 2 3 -o JSON
```

Results:

```JSON
{
  "publicIp": {
    "ddosSettings": {
      "protectionMode": "VirtualNetworkInherited"
    },
    "dnsSettings": {
      "domainNameLabel": "mydnslabel242",
      "fqdn": "mydnslabel242.eastus.cloudapp.azure.com"
    },
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/publicIPAddresses/myPublicIP242",
    "idleTimeoutInMinutes": 4,
    "ipAddress": "20.62.231.145",
    "ipTags": [],
    "location": "eastus",
    "name": "myPublicIP242",
    "provisioningState": "Succeeded",
    "publicIPAddressVersion": "IPv4",
    "publicIPAllocationMethod": "Static",
    "resourceGroup": "myResourceGroup242",
    "sku": {
      "name": "Standard",
      "tier": "Regional"
    },
    "type": "Microsoft.Network/publicIPAddresses",
    "zones": [
      "1",
      "2",
      "3"
    ]
  }
}
```

Security rules in network security groups enable you to filter the type of network traffic that can flow in and out of virtual network subnets and network interfaces. To learn more about network security groups, see [Network security group overview](https://learn.microsoft.com/azure/virtual-network/network-security-groups-overview).

```bash
az network nsg create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_NSG_NAME \
    --location $MY_LOCATION
```

Results:

```JSON
{
  "NewNSG": {
    "defaultSecurityRules": [
      {
        "access": "Allow",
        "description": "Allow inbound traffic from all VMs in VNET",
        "destinationAddressPrefix": "VirtualNetwork",
        "destinationAddressPrefixes": [],
        "destinationPortRange": "*",
        "destinationPortRanges": [],
        "direction": "Inbound",
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/networkSecurityGroups/myNSGName242/defaultSecurityRules/AllowVnetInBound",
        "name": "AllowVnetInBound",
        "priority": 65000,
        "protocol": "*",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup242",
        "sourceAddressPrefix": "VirtualNetwork",
        "sourceAddressPrefixes": [],
        "sourcePortRange": "*",
        "sourcePortRanges": [],
        "type": "Microsoft.Network/networkSecurityGroups/defaultSecurityRules"
      },
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/networkSecurityGroups/myNSGName242",
        "location": "eastus",
        "name": "myNSGName242",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup242",
        "resourceGuid": "1aa33cc2-109c-4144-b6eb-64d0d56711a6",
        "securityRules": [],
        "type": "Microsoft.Network/networkSecurityGroups"
  }
}
```

Open ports 22 (SSH), 80 (HTTP) and 443 (HTTPS) to allow SSH and Web traffic

```bash
az network nsg rule create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --nsg-name $MY_NSG_NAME \
    --name Port_22 \
    --protocol tcp \
    --priority 200\
    --destination-port-range 22 \
    --access allow

az network nsg rule create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --nsg-name $MY_NSG_NAME \
    --name Port_80 \
    --protocol tcp \
    --priority 300\
    --destination-port-range 80 \
    --access allow

az network nsg rule create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --nsg-name $MY_NSG_NAME \
    --name Port_443 \
    --protocol tcp \
    --priority 400\
    --destination-port-range 443 \
    --access allow
```

```JSON

{
  "access": "Allow",
  "destinationAddressPrefix": "*",
  "destinationAddressPrefixes": [],
  "destinationPortRange": "22",
  "destinationPortRanges": [],
  "direction": "Inbound",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/networkSecurityGroups/myNSGName242/securityRules/Port_22",
  "name": "Port_22",
  "priority": 200,
  "protocol": "Tcp",
  "provisioningState": "Succeeded",
  "resourceGroup": "myResourceGroup242",
  "sourceAddressPrefix": "*",
  "sourceAddressPrefixes": [],
  "sourcePortRange": "*",
  "sourcePortRanges": [],
  "type": "Microsoft.Network/networkSecurityGroups/securityRules"
}

{
  "access": "Allow",
  "destinationAddressPrefix": "*",
  "destinationAddressPrefixes": [],
  "destinationPortRange": "80",
  "destinationPortRanges": [],
  "direction": "Inbound",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/networkSecurityGroups/myNSGName242/securityRules/Port_80",
  "name": "Port_80",
  "priority": 300,
  "protocol": "Tcp",
  "provisioningState": "Succeeded",
  "resourceGroup": "myResourceGroup242",
  "sourceAddressPrefix": "*",
  "sourceAddressPrefixes": [],
  "sourcePortRange": "*",
  "sourcePortRanges": [],
  "type": "Microsoft.Network/networkSecurityGroups/securityRules"
}

{
  "access": "Allow",
  "destinationAddressPrefix": "*",
  "destinationAddressPrefixes": [],
  "destinationPortRange": "443",
  "destinationPortRanges": [],
  "direction": "Inbound",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/networkSecurityGroups/myNSGName242/securityRules/Port_443",
  "name": "Port_443",
  "priority": 400,
  "protocol": "Tcp",
  "provisioningState": "Succeeded",
  "resourceGroup": "myResourceGroup242",
  "sourceAddressPrefix": "*",
  "sourceAddressPrefixes": [],
  "sourcePortRange": "*",
  "sourcePortRanges": [],
  "type": "Microsoft.Network/networkSecurityGroups/securityRules"
}
```

Associate NSG to subnet

```bash
az network vnet subnet update \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --vnet-name $MY_VNET_NAME \
    --name $MY_SN_NAME \
    --network-security-group $MY_NSG_NAME
```

Results:

```JSON
{
  "addressPrefix": "10.242.0.0/24",
  "delegations": [],
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/virtualNetworks/myVNet242/subnets/mySN242",
  "name": "mySN242",
  "networkSecurityGroup": {
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.Network/networkSecurityGroups/myNSGName242",
    "resourceGroup": "myResourceGroup242"
  },
  "privateEndpointNetworkPolicies": "Disabled",
  "privateLinkServiceNetworkPolicies": "Enabled",
  "provisioningState": "Succeeded",
  "resourceGroup": "myResourceGroup242",
  "type": "Microsoft.Network/virtualNetworks/subnets"
}
```
## Create an Azure Key Vault

```bash
az keyvault create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_KEY_VAULT \
    --location $MY_LOCATION \
    --retention-days 7\
    --enabled-for-deployment   
```

Results:

```JSON
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup242/providers/Microsoft.KeyVault/vaults/myKeyVault242",
  "location": "eastus",
  "name": "myKeyVault242",
  "properties": {
    "accessPolicies": [
      {
        "applicationId": null,
        "objectId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
        "permissions": {
          "certificates": [
            "all"
          ],
          "keys": [
            "all"
          ],
          "secrets": [
            "all"
          ],
          "storage": [
            "all"
          ]
        },
        "tenantId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      }
    ],
    "createMode": null,
    "enablePurgeProtection": null,
    "enableRbacAuthorization": null,
    "enableSoftDelete": true,
    "enabledForDeployment": true,
    "enabledForDiskEncryption": null,
    "enabledForTemplateDeployment": null,
    "hsmPoolResourceId": null,
    "networkAcls": null,
    "privateEndpointConnections": null,
    "provisioningState": "Succeeded",
    "publicNetworkAccess": "Enabled",
    "sku": {
      "family": "A",
      "name": "standard"
    },
    "softDeleteRetentionInDays": 90,
    "tenantId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "vaultUri": "https://mykeyvault242.vault.azure.net/"
  },
  "resourceGroup": "myResourceGroup242",
  "systemData": {
    "createdAt": "2023-09-08T14:55:45.691000+00:00",
    "createdBy": "dummy.dummy@outlook.pt",
    "createdByType": "User",
    "lastModifiedAt": "2023-09-08T14:55:45.691000+00:00",
    "lastModifiedBy": "dummy.dummy@outlook.pt",
    "lastModifiedByType": "User"
  },
  "tags": {},
  "type": "Microsoft.KeyVault/vaults"
}
```

## Create a certificate and store in Azure key Vault

For this article we’ll use a self signed certificate.

```bash
az keyvault certificate create \
    --vault-name $MY_KEY_VAULT \
    --name nginxcert \
    --policy "$(az keyvault certificate get-default-policy)"
```

```JSON
{
  "cancellationRequested": false,
  "csr": "MIICrjCCAZYCA(...)K6ibPBZqhIH",
  "error": null,
  "id": "https://<MY_KEY_VAULT>.vault.azure.net/certificates/nginxcert/pending",
  "issuerParameters": {
    "certificateTransparency": null,
    "certificateType": null,
    "name": "Self"
  },
  "name": "nginxcert",
  "requestId": "2109088929f3437c9da91bd69827f9a9",
  "status": "completed",
  "statusDetails": null,
  "target": "https://<MY_KEY_VAULT>.vault.azure.net/certificates/nginxcert"
}
```

## Create the VM

Now create a VM with az vm create. Use the --custom-data parameter to pass in your cloud-init config file. Provide the full path to the cloud-init.txt config if you saved the file outside of your present working directory.

```bash
cat > cloud-init-nginx.txt <<EOF
#cloud-config
# Install, update, and upgrade packages
package_upgrade: true
package_update: true
# Install packages
packages:
  - nginx
write_files:
  - owner: www-data:www-data
  - path: /etc/nginx/sites-available/default
    content: |
      server {
        listen 443 ssl;
        ssl_certificate /etc/nginx/ssl/nginxcert.cert;
        ssl_certificate_key /etc/nginx/ssl/nginxcert.prv;
      }
runcmd:
  - mkdir /etc/nginx/ssl
  - service nginx restart
EOF
```

The following example creates a VM named *myVMName$UNIQUE_POSTFIX*:

```bash
az vm create \
  --resource-group $MY_RESOURCE_GROUP_NAME \
  --name $MY_VM_NAME \
  --image $MY_VM_IMAGE \
  --vnet-name $MY_VNET_NAME --subnet $MY_SN_NAME \
  --admin-username $MY_VM_USERNAME \
  --generate-ssh-keys \
  --size $MY_VM_SIZE \
  --custom-data cloud-init-nginx.txt \
  --public-ip-address $MY_PUBLIC_IP_NAME
```

Results:

Enable the system assigned identity on a VM with the 'Contributor' role.
//https://learn.microsoft.com/en-us/cli/azure/vm/identity?view=azure-cli-latest#az-vm-identity-assign

```bash
az vm identity assign \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_VM_NAME \
    --role Contributor \
    --scope /subscriptions/0bb78609-cc8b-4e7d-be30-eee8cf2dbea4/resourceGroups/myResourceGroup3
```

```bash
az keyvault set-policy \
    --name $MY_KEY_VAULT \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --object-id 'a93f8699-ab75-4b63-a3f2-2e59d114561d' \
    --secret-permissions get list delete
```

Start extension deployment

```bash
settings_value=$(cat <<EOF
{
  "secretsManagementSettings": {
    "pollingIntervalInS": "3600",
    "certificateStoreLocation": "/var/lib/waagent/Microsoft.Azure.KeyVault",
    "observedCertificates": [
      "https://mykeyvault210.vault.azure.net/certificates/nginxcert"
    ]
  }
}
EOF
```

```bash
az vm extension set -n "KeyVaultForLinux" \
     --publisher Microsoft.Azure.KeyVault \
     -g $MY_RESOURCE_GROUP_NAME \
     --vm-name $MY_VM_NAME \
     --version 2.0 \
     --enable-auto-upgrade true \
     --settings "$settings_value"
```
