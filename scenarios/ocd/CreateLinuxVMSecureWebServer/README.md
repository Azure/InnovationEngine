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
export NETWORK_PREFIX="$(($RANDOM % 254 + 1))"
export RANDOM_ID="$(openssl rand -hex 3)"
export MY_RESOURCE_GROUP_NAME="myResourceGroup$RANDOM_ID"
export MY_KEY_VAULT="mykeyvault$RANDOM_ID"
export MY_CERT_NAME="nginxcert$RANDOM_ID"
export REGION="eastus"
export MY_VM_NAME="myVMName$RANDOM_ID"
export MY_VM_ID_NAME="myVMIDName$RANDOM_ID"
export MY_VM_IMAGE='Ubuntu2204'
export MY_VM_USERNAME="azureuser"
export MY_VM_SIZE='Standard_DS2_v2'
export MY_VNET_NAME="myVNet$RANDOM_ID"
export MY_VM_NIC_NAME="myVMNicName$RANDOM_ID"
export MY_NSG_SSH_RULE="Allow-Access$RANDOM_ID"
export MY_VM_NIC_NAME="myVMNicName$RANDOM_ID"
export MY_VNET_PREFIX="10.$NETWORK_PREFIX.0.0/16"
export MY_SN_NAME="mySN$RANDOM_ID"
export MY_SN_PREFIX="10.$NETWORK_PREFIX.0.0/24"
export MY_PUBLIC_IP_NAME="myPublicIP$RANDOM_ID"
export MY_DNS_LABEL="mydnslabel$RANDOM_ID"
export MY_NSG_NAME="myNSGName$RANDOM_ID"
export FQDN="${MY_DNS_LABEL}.${REGION}.cloudapp.azure.com"
```

## Create a Resource Group

Before you can create a secure Linux VM, create a resource group with az group create. The following example creates a resource group equal to the contents of the variable *MY_RESOURCE_GROUP_NAME* in the location specified by the variable contents *REGION*:

```bash
az group create \
    --name $MY_RESOURCE_GROUP_NAME \
    --location $REGION -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f",
  "location": "eastus",
  "managedBy": null,
  "name": "myResourceGroupb1404f",
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
    --location $REGION \
    --address-prefix $MY_VNET_PREFIX \
    --subnet-name $MY_SN_NAME \
    --subnet-prefix $MY_SN_PREFIX -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "newVNet": {
    "addressSpace": {
      "addressPrefixes": [
        "10.168.0.0/16"
      ]
    },
    "bgpCommunities": null,
    "ddosProtectionPlan": null,
    "dhcpOptions": {
      "dnsServers": []
    },
    "enableDdosProtection": false,
    "enableVmProtection": null,
    "encryption": null,
    "extendedLocation": null,
    "flowTimeoutInMinutes": null,
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/virtualNetworks/myVNetb1404f",
    "ipAllocations": null,
    "location": "eastus",
    "name": "myVNetb1404f",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroupb1404f",
    "subnets": [
      {
        "addressPrefix": "10.168.0.0/24",
        "addressPrefixes": null,
        "applicationGatewayIpConfigurations": null,
        "delegations": [],
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/virtualNetworks/myVNetb1404f/subnets/mySNb1404f",
        "ipAllocations": null,
        "ipConfigurationProfiles": null,
        "ipConfigurations": null,
        "name": "mySNb1404f",
        "natGateway": null,
        "networkSecurityGroup": null,
        "privateEndpointNetworkPolicies": "Disabled",
        "privateEndpoints": null,
        "privateLinkServiceNetworkPolicies": "Enabled",
        "provisioningState": "Succeeded",
        "purpose": null,
        "resourceGroup": "myResourceGroupb1404f",
        "resourceNavigationLinks": null,
        "routeTable": null,
        "serviceAssociationLinks": null,
        "serviceEndpointPolicies": null,
        "serviceEndpoints": null,
        "type": "Microsoft.Network/virtualNetworks/subnets"
      }
    ],
    "tags": {},
    "type": "Microsoft.Network/virtualNetworks",
    "virtualNetworkPeerings": []
  }
}
```

Use az network public-ip create to create a standard zone-redundant public IPv4 address named *$MY_PUBLIC_IP_NAME* in *$MY_RESOURCE_GROUP_NAME*.

```bash
az network public-ip create \
    --name $MY_PUBLIC_IP_NAME \
    --location $REGION \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --dns-name $MY_DNS_LABEL \
    --sku Standard \
    --allocation-method static \
    --version IPv4 \
    --zone 1 2 3 -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "publicIp": {
    "ddosSettings": null,
    "deleteOption": null,
    "dnsSettings": {
      "domainNameLabel": "mydnslabelb1404f",
      "fqdn": "mydnslabelb1404f.eastus.cloudapp.azure.com",
      "reverseFqdn": null
    },
    "extendedLocation": null,
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/publicIPAddresses/myPublicIPb1404f",
    "idleTimeoutInMinutes": 4,
    "ipAddress": "20.88.178.210",
    "ipConfiguration": null,
    "ipTags": [],
    "linkedPublicIpAddress": null,
    "location": "eastus",
    "migrationPhase": null,
    "name": "myPublicIPb1404f",
    "natGateway": null,
    "provisioningState": "Succeeded",
    "publicIpAddressVersion": "IPv4",
    "publicIpAllocationMethod": "Static",
    "publicIpPrefix": null,
    "resourceGroup": "myResourceGroupb1404f",
    "servicePublicIpAddress": null,
    "sku": {
      "name": "Standard",
      "tier": "Regional"
    },
    "tags": null,
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
    --location $REGION -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
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
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/networkSecurityGroups/myNSGNameb1404f/defaultSecurityRules/AllowVnetInBound",
        "name": "AllowVnetInBound",
        "priority": 65000,
        "protocol": "*",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroupb1404f",
        "sourceAddressPrefix": "VirtualNetwork",
        "sourceAddressPrefixes": [],
        "sourcePortRange": "*",
        "sourcePortRanges": [],
        "type": "Microsoft.Network/networkSecurityGroups/defaultSecurityRules"
      },
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/networkSecurityGroups/myNSGNameb1404f",
        "location": "eastus",
        "name": "myNSGNameb1404f",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroupb1404f",
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
    --name $MY_NSG_SSH_RULE \
    --access Allow \
    --protocol Tcp \
    --direction Inbound \
    --priority 100 \
    --source-address-prefix '*' \
    --source-port-range '*' \
    --destination-address-prefix '*' \
    --destination-port-range 22 80 443 -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "access": "Allow",
  "description": null,
  "destinationAddressPrefix": "*",
  "destinationAddressPrefixes": [],
  "destinationApplicationSecurityGroups": null,
  "destinationPortRange": null,
  "destinationPortRanges": [
    "22",
    "80",
    "443"
  ],
  "direction": "Inbound",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/networkSecurityGroups/myNSGNameb1404f/securityRules/MY_NSG_SSH_RULE",
  "name": "MY_NSG_SSH_RULE",
  "priority": 100,
  "protocol": "Tcp",
  "provisioningState": "Succeeded",
  "resourceGroup": "myResourceGroupb1404f",
  "sourceAddressPrefix": "*",
  "sourceAddressPrefixes": [],
  "sourceApplicationSecurityGroups": null,
  "sourcePortRange": "*",
  "sourcePortRanges": [],
  "type": "Microsoft.Network/networkSecurityGroups/securityRules"
}
```

And finally create the Network Interface Card (NIC):

```bash
az network nic create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_VM_NIC_NAME \
    --location $REGION \
    --ip-forwarding false \
    --subnet $MY_SN_NAME \
    --vnet-name $MY_VNET_NAME \
    --network-security-group $MY_NSG_NAME \
    --public-ip-address $MY_PUBLIC_IP_NAME -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "NewNIC": {
    "auxiliaryMode": "None",
    "auxiliarySku": "None",
    "disableTcpStateTracking": false,
    "dnsSettings": {
      "appliedDnsServers": [],
      "dnsServers": []
    },
    "enableAcceleratedNetworking": false,
    "enableIPForwarding": false,
    "hostedWorkloads": [],
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/networkInterfaces/myVMNicNameb1404f",
    "ipConfigurations": [
      {
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/networkInterfaces/myVMNicNameb1404f/ipConfigurations/ipconfig1",
        "name": "ipconfig1",
        "primary": true,
        "privateIPAddress": "10.168.0.4",
        "privateIPAddressVersion": "IPv4",
        "privateIPAllocationMethod": "Dynamic",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroupb1404f",
        "subnet": {
          "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/virtualNetworks/myVNetb1404f/subnets/mySNb1404f",
          "resourceGroup": "myResourceGroupb1404f"
        },
        "type": "Microsoft.Network/networkInterfaces/ipConfigurations"
      }
    ],
    "location": "eastus",
    "name": "myVMNicNameb1404f",
    "networkSecurityGroup": {
      "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.Network/networkSecurityGroups/myNSGNameb1404f",
      "resourceGroup": "myResourceGroupb1404f"
    },
    "nicType": "Standard",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroupb1404f",
    "tapConfigurations": [],
    "type": "Microsoft.Network/networkInterfaces",
    "vnetEncryptionSupported": false
  }
}
```

## Generate a certificate and store it in Azure Key Vault

Azure Key Vault safeguards cryptographic keys and secrets, such as certificates or passwords. Key Vault helps streamline the certificate management process and enables you to maintain control of keys that access those certificates. You can create a self-signed certificate inside Key Vault, or upload an existing, trusted certificate that you already own. For this tutorial we'll create self-signed certificates inside the Key Vault and afterwards inject these certificates into a running VM. This process ensures that the most up-to-date certificates are installed on a web server during deployment.

The following example creates an Azure Key Vault named *$MY_KEY_VAULT* in the chosen region *$REGION* with a retention policy of 7 days. This means once a secret, key, certificate, or key vault is deleted, it will remain recoverable for a configurable period of 7 to 90 calendar days.

```bash
az keyvault create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_KEY_VAULT \
    --location $REGION \
    --retention-days 7\
    --enabled-for-deployment true -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.KeyVault/vaults/myKeyVaultb1404f",
  "location": "eastus",
  "name": "myKeyVaultb1404f",
  "properties": {
    "accessPolicies": [
      {
        "applicationId": null,
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
    "softDeleteRetentionInDays": 7,
    "tenantId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "vaultUri": "https://mykeyvaultb1404f.vault.azure.net/"
  },
  "resourceGroup": "myResourceGroupb1404f",
  "systemData": {
    "createdAt": "2023-09-18T12:25:55.208000+00:00",
    "createdBy": "example@microsoft.com",
    "createdByType": "User",
    "lastModifiedAt": "2023-09-18T12:25:55.208000+00:00",
    "lastModifiedBy": "example@microsoft.com",
    "lastModifiedByType": "User"
  },
  "tags": {},
  "type": "Microsoft.KeyVault/vaults"
}
```

## Create a certificate and store in Azure key Vault

Now let's generate a self-signed certificate with az keyvault certificate create that uses the default certificate policy:

```bash
az keyvault certificate create \
    --vault-name $MY_KEY_VAULT \
    --name $MY_CERT_NAME \
    --policy "$(az keyvault certificate get-default-policy)" -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "cancellationRequested": false,
  "csr": "MIICr...",
  "error": null,
  "id": "https://mykeyvault67a7ba.vault.azure.net/certificates/nginxcert67a7ba/pending",
  "issuerParameters": {
    "certificateTransparency": null,
    "certificateType": null,
    "name": "Self"
  },
  "name": "nginxcert67a7ba",
  "status": "completed",
  "statusDetails": null,
  "target": "https://mykeyvault67a7ba.vault.azure.net/certificates/nginxcert67a7ba"
}
```

Finally, we need to prepare the certificate so it can be used during the VM create process. To do so we need to obtain the ID of the certificate with az keyvault secret list-versions, and convert the certificate with az vm secret format. The following example assigns the output of these commands to variables for ease of use in the next steps:

```bash
az identity create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_VM_ID_NAME -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "clientId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourcegroups/myResourceGroupb1404f/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myVMIDNameb1404f",
  "location": "eastus",
  "name": "myVMIDNameb1404f",
  "principalId": "e09ebfce-97f0-4aff-9abd-415ebd6f915c",
  "resourceGroup": "myResourceGroupb1404f",
  "tags": {},
  "tenantId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "type": "Microsoft.ManagedIdentity/userAssignedIdentities"
}
```

```bash
MY_VM_PRINCIPALID=$(az identity show --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_VM_ID_NAME --query principalId -o tsv)

az keyvault set-policy \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_KEY_VAULT \
    --object-id $MY_VM_PRINCIPALID \
    --secret-permissions get list \
    --certificate-permissions get list -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupb1404f/providers/Microsoft.KeyVault/vaults/myKeyVaultb1404f",
  "location": "eastus",
  "name": "myKeyVaultb1404f",
  "properties": {
    "accessPolicies": [
      {
        "applicationId": null,
        "objectId": "ceeb4e98-5831-4d9f-b8ba-2ee14b3cdf80",
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
        "tenantId": "bd7153ee-d085-4a28-a928-2f0ef402f076"
      },
      {
        "applicationId": null,
        "objectId": "e09ebfce-97f0-4aff-9abd-415ebd6f915c",
        "permissions": {
          "certificates": [
            "list",
            "get"
          ],
          "keys": null,
          "secrets": [
            "list",
            "get"
          ],
          "storage": null
        },
        "tenantId": "bd7153ee-d085-4a28-a928-2f0ef402f076"
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
    "softDeleteRetentionInDays": 7,
    "tenantId": "bd7153ee-d085-4a28-a928-2f0ef402f076",
    "vaultUri": "https://mykeyvaultb1404f.vault.azure.net/"
  },
  "resourceGroup": "myResourceGroupb1404f",
  "systemData": {
    "createdAt": "2023-09-18T12:25:55.208000+00:00",
    "createdBy": "ajoian@microsoft.com",
    "createdByType": "User",
    "lastModifiedAt": "2023-09-18T12:48:08.966000+00:00",
    "lastModifiedBy": "ajoian@microsoft.com",
    "lastModifiedByType": "User"
  },
  "tags": {},
  "type": "Microsoft.KeyVault/vaults"
}
```

## Create the VM

Now create a VM with az vm create. Use the --custom-data parameter to pass in the cloud-init config file, named *cloud-init-nginx.txt*.
Cloud-init is a widely used approach to customize a Linux VM as it boots for the first time. You can use cloud-init to install packages and write files, or to configure users and security. As cloud-init runs during the initial boot process, there are no extra steps or required agents to apply your configuration.
When you create a VM, certificates and keys are stored in the protected /var/lib/waagent/ directory. In this example, we are installing and configuring the NGINX web server.

```bash
cat > cloud-init-nginx.txt <<EOF
#cloud-config

# Install, update, and upgrade packages
package_upgrade: true
package_update: true
package_reboot_if_require: true

# Install packages
packages:
  - nginx

write_files:
  - owner: www-data:www-data
  - path: /etc/nginx/sites-available/secure-server
    content: |
      server {
        server_name $FQDN;
        listen 443 ssl http2;
        ssl_certificate /etc/nginx/ssl/${MY_CERT_NAME}.crt;
        ssl_certificate_key /etc/nginx/ssl/${MY_CERT_NAME}.key;
      }
      server {
            listen 80;
            server_name $FQDN;
            return 301 https://$FQDN\$request_uri;
      }

  - owner: root:root
    path: /root/convert_akv_cert.sh
    permissions: "0750"
    content: |
        #!/bin/bash
        runtime="10 minute"; endtime=\$(date -ud "\$runtime" +%s); 
        while [[ \$(date -u +%s) -le \$endtime ]]; do
          if [[ ! -f /etc/nginx/ssl/${MY_KEY_VAULT}.${MY_CERT_NAME} ]]
            then
                sleep 5;
            else
                break;
          fi;
        done
        # Split the file in two (cert and key)
        echo "Creating .key file with private key..."
        openssl rsa -outform pem -in /etc/nginx/ssl/${MY_KEY_VAULT}.${MY_CERT_NAME} -out /etc/nginx/ssl/${MY_CERT_NAME}.key
        echo "Creating .crt file with certificate..."
        openssl x509 -outform pem -in /etc/nginx/ssl/${MY_KEY_VAULT}.${MY_CERT_NAME} -out /etc/nginx/ssl/${MY_CERT_NAME}.crt
runcmd:
  - mkdir /etc/nginx/ssl
  - ln -s /etc/nginx/sites-available/secure-server /etc/nginx/sites-enabled/
  - rm /etc/nginx/sites-enabled/default
  - bash /root/convert_akv_cert.sh
  - (crontab -l 2>/dev/null; echo "0 * * * * /root/convert_akv_cert.sh && service nginx reload") | crontab -
  - service nginx restart
EOF
```

The following example creates a VM named *myVMName$UNIQUE_POSTFIX*:

```bash
MY_VM_ID=$(az identity show --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_VM_ID_NAME --query id -o tsv)

az vm create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_VM_NAME \
    --image $MY_VM_IMAGE \
    --admin-username $MY_VM_USERNAME \
    --generate-ssh-keys \
    --assign-identity $MY_VM_ID \
    --size $MY_VM_SIZE \
    --custom-data cloud-init-nginx.txt \
    --nics $MY_VM_NIC_NAME
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "fqdns": "mydnslabel67a7ba.eastus.cloudapp.azure.com",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup67a7ba/providers/Microsoft.Compute/virtualMachines/myVMName67a7ba",
  "identity": {
    "systemAssignedIdentity": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "userAssignedIdentities": {
      "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourcegroups/myResourceGroup67a7ba/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myVMIDName67a7ba": {
        "clientId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
        "principalId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      }
    }
  },
  "location": "eastus",
  "macAddress": "60-45-BD-D3-B5-29",
  "powerState": "VM running",
  "privateIpAddress": "10.56.0.4",
  "publicIpAddress": "20.231.118.239",
  "resourceGroup": "myResourceGroup67a7ba",
  "zones": ""
}
```

## Deploying AKV extension for VM $vm_name to retrieve cert $cert_name from AKV $akv_name..."

```bash
MY_CERT_ID=$(az keyvault certificate show --vault-name $MY_KEY_VAULT --name $MY_CERT_NAME --query sid -o tsv)
MY_VM_CLIENTID=$(az identity show --resource-group $MY_RESOURCE_GROUP_NAME --name $MY_VM_ID_NAME --query clientId -o tsv)
MY_AKV_EXT_SETTINGS="{\"secretsManagementSettings\":{\"pollingIntervalInS\":\"3600\",\"requireInitialSync\":"true",\"certificateStoreLocation\":\"/etc/nginx/ssl/\",\"observedCertificates\":[\"$MY_CERT_ID\"]},\"authenticationSettings\":{\"msiClientId\":\"${MY_VM_CLIENTID}\"}}"

az vm extension set \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --vm-name $MY_VM_NAME \
    -n "KeyVaultForLinux" \
    --publisher Microsoft.Azure.KeyVault \
    --version 2.0 \
    --enable-auto-upgrade true \
    --settings $MY_AKV_EXT_SETTINGS -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "autoUpgradeMinorVersion": true,
  "enableAutomaticUpgrade": true,
  "forceUpdateTag": null,
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup67a7ba/providers/Microsoft.Compute/virtualMachines/myVMName67a7ba/extensions/KeyVaultForLinux",
  "instanceView": null,
  "location": "eastus",
  "name": "KeyVaultForLinux",
  "protectedSettings": null,
  "protectedSettingsFromKeyVault": null,
  "provisioningState": "Succeeded",
  "publisher": "Microsoft.Azure.KeyVault",
  "resourceGroup": "myResourceGroup67a7ba",
  "settings": {
    "secretsManagementSettings": {
      "certificateStoreLocation": "/etc/nginx/ssl",
      "observedCertificates": [
        "https://mykeyvault67a7ba.vault.azure.net/secrets/nginxcert67a7ba/aac9b30a90c04fc58bc230ae15b1148f"
      ],
      "pollingIntervalInS": "3600"
    }
  },
  "suppressFailures": null,
  "tags": null,
  "type": "Microsoft.Compute/virtualMachines/extensions",
  "typeHandlerVersion": "2.0",
  "typePropertiesType": "KeyVaultForLinux"
}
```

## Enable Azure AD login for a Linux Virtual Machine in Azure

The following example deploys a VM and then installs the extension to enable Azure AD login for a Linux VM. VM extensions are small applications that provide post-deployment configuration and automation tasks on Azure virtual machines.

```bash
az vm extension set \
    --publisher Microsoft.Azure.ActiveDirectory \
    --name AADSSHLoginForLinux \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --vm-name $MY_VM_NAME -o JSON
```

Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "autoUpgradeMinorVersion": true,
  "enableAutomaticUpgrade": null,
  "forceUpdateTag": null,
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupfa636b/providers/Microsoft.Compute/virtualMachines/myVMNamefa636b/extensions/AADSSHLoginForLinux",
  "instanceView": null,
  "location": "eastus",
  "name": "AADSSHLoginForLinux",
  "protectedSettings": null,
  "protectedSettingsFromKeyVault": null,
  "provisioningState": "Succeeded",
  "publisher": "Microsoft.Azure.ActiveDirectory",
  "resourceGroup": "myResourceGroupfa636b",
  "settings": null,
  "suppressFailures": null,
  "tags": null,
  "type": "Microsoft.Compute/virtualMachines/extensions",
  "typeHandlerVersion": "1.0",
  "typePropertiesType": "AADSSHLoginForLinux"
}
```

## Browse your secure website

Validate that the application is running by visiting the application url:

```bash
curl --max-time 120 -k "https://$FQDN"
```

Results:

<!-- expected_similarity=0.3 -->
```html
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```
