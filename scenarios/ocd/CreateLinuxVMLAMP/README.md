# Variable declaration

```bash
export NETWORK_PREFIX="$(($RANDOM % 254 + 1))"
export RANDOM_ID="$(openssl rand -hex 3)"
export MY_RESOURCE_GROUP_NAME="myResourceGroup$RANDOM_ID"
export MY_LOCATION="eastus"
export MY_VM_NAME="myVMName$RANDOM_ID"
export MY_VM_USERNAME="azureadmin"
export MY_VM_SIZE='Standard_DS2_v2'
export MY_VM_IMAGE='Canonical:0001-com-ubuntu-minimal-jammy:minimal-22_04-lts-gen2:latest'
export MY_PUBLIC_IP_NAME="myPublicIP$RANDOM_ID"
export MY_DNS_LABEL="mydnslabel$RANDOM_ID"
export MY_NSG_NAME="myNSGName$RANDOM_ID"
export MY_NSG_SSH_RULE="Allow-Access$RANDOM_ID"
export MY_VM_NIC_NAME="myVMNicName$RANDOM_ID"
export MY_VNET_NAME="myVNet$RANDOM_ID"
export MY_VNET_PREFIX="10.$NETWORK_PREFIX.0.0/22"
export MY_SN_NAME="mySN$RANDOM_ID"
export MY_SN_PREFIX="10.$NETWORK_PREFIX.0.0/24"
export MY_MYSQL_DB_NAME="mydb$RANDOM_ID"
export MY_MYSQL_ADMIN_USERNAME="dbadmin$RANDOM_ID"
export MY_MYSQL_ADMIN_PW="$(openssl rand -base64 32)"
export MY_MYSQL_SN_NAME="myMySQLSN$RANDOM_ID"
export MY_WP_ADMIN_PW="$(openssl rand -base64 32)"
export MY_WP_ADMIN_USER="wpcliadmin"
export MY_SSL_EMAIL="$(az account show --query user.name --output tsv)"
export FQDN="${MY_DNS_LABEL}.${MY_LOCATION}.cloudapp.azure.com"
```

# Create RG
```bash
az group create \
    --name $MY_RESOURCE_GROUP_NAME \
    --location $MY_LOCATION
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc",
  "location": "eastus",
  "managedBy": null,
  "name": "myResourceGroup6ad2bc",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```

# Setup LAMP networking
```bash
az network vnet create \
    --name $MY_VNET_NAME \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --location $MY_LOCATION \
    --address-prefix $MY_VNET_PREFIX \
    --subnet-name $MY_SN_NAME \
    --subnet-prefixes $MY_SN_PREFIX
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "newVNet": {
    "addressSpace": {
      "addressPrefixes": [
        "10.19.0.0/22"
      ]
    },
    "enableDdosProtection": false,
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/virtualNetworks/myVNet6ad2bc",
    "location": "eastus",
    "name": "myVNet6ad2bc",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup6ad2bc",
    "subnets": [
      {
        "addressPrefix": "10.19.0.0/24",
        "delegations": [],
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/virtualNetworks/myVNet6ad2bc/subnets/mySN6ad2bc",
        "name": "mySN6ad2bc",
        "privateEndpointNetworkPolicies": "Disabled",
        "privateLinkServiceNetworkPolicies": "Enabled",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup6ad2bc",
        "type": "Microsoft.Network/virtualNetworks/subnets"
      }
    ],
    "type": "Microsoft.Network/virtualNetworks",
    "virtualNetworkPeerings": []
  }
}
```

```bash
az network public-ip create \
    --name $MY_PUBLIC_IP_NAME \
    --location $MY_LOCATION \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --dns-name $MY_DNS_LABEL \
    --sku Standard \
    --allocation-method static \
    --version IPv4 \
    --zone 1 2 3
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "publicIp": {
    "ddosSettings": {
      "protectionMode": "VirtualNetworkInherited"
    },
    "dnsSettings": {
      "domainNameLabel": "mydnslabel6ad2bc",
      "fqdn": "mydnslabel6ad2bc.eastus.cloudapp.azure.com"
    },
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/publicIPAddresses/myPublicIP6ad2bc",
    "idleTimeoutInMinutes": 4,
    "ipTags": [],
    "location": "eastus",
    "name": "myPublicIP6ad2bc",
    "provisioningState": "Succeeded",
    "publicIPAddressVersion": "IPv4",
    "publicIPAllocationMethod": "Static",
    "resourceGroup": "myResourceGroup6ad2bc",
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

```bash
az network nsg create \
    --name $MY_NSG_NAME \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --location $MY_LOCATION
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
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup104/providers/Microsoft.Network/networkSecurityGroups/protect-vms/defaultSecurityRules/AllowVnetInBound",
        "name": "AllowVnetInBound",
        "priority": 65000,
        "protocol": "*",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup104",
        "sourceAddressPrefix": "VirtualNetwork",
        "sourceAddressPrefixes": [],
        "sourcePortRange": "*",
        "sourcePortRanges": [],
        "type": "Microsoft.Network/networkSecurityGroups/defaultSecurityRules"
      },
    ],
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup104/providers/Microsoft.Network/networkSecurityGroups/protect-vms",
    "location": "eastus",
    "name": "protect-vms",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup104",
    "securityRules": [],
    "type": "Microsoft.Network/networkSecurityGroups"
  }
}
```

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
    --destination-port-range 22 80 443
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "access": "Allow",
  "destinationAddressPrefix": "*",
  "destinationAddressPrefixes": [],
  "destinationPortRanges": [
    "22",
    "80",
    "443"
  ],
  "direction": "Inbound",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/networkSecurityGroups/myNSGName6ad2bc/securityRules/Allow-Access6ad2bc",
  "name": "Allow-Access6ad2bc",
  "priority": 100,
  "protocol": "Tcp",
  "provisioningState": "Succeeded",
  "resourceGroup": "myResourceGroup6ad2bc",
  "sourceAddressPrefix": "*",
  "sourceAddressPrefixes": [],
  "sourcePortRange": "*",
  "sourcePortRanges": [],
  "type": "Microsoft.Network/networkSecurityGroups/securityRules"
}
```

```bash
az network nic create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_VM_NIC_NAME \
    --location $MY_LOCATION \
    --ip-forwarding false \
    --subnet $MY_SN_NAME \
    --vnet-name $MY_VNET_NAME \
    --network-security-group $MY_NSG_NAME
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
      "dnsServers": [],
    },
    "enableAcceleratedNetworking": false,
    "enableIPForwarding": false,
    "hostedWorkloads": [],
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/networkInterfaces/myVMNicName6ad2bc",
    "ipConfigurations": [
      {
        "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/networkInterfaces/myVMNicName6ad2bc/ipConfigurations/ipconfig1",
        "name": "ipconfig1",
        "primary": true,
        "privateIPAddress": "10.19.0.4",
        "privateIPAddressVersion": "IPv4",
        "privateIPAllocationMethod": "Dynamic",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup6ad2bc",
        "subnet": {
          "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/virtualNetworks/myVNet6ad2bc/subnets/mySN6ad2bc",
          "resourceGroup": "myResourceGroup6ad2bc"
        },
        "type": "Microsoft.Network/networkInterfaces/ipConfigurations"
      }
    ],
    "location": "eastus",
    "name": "myVMNicName6ad2bc",
    "networkSecurityGroup": {
      "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/networkSecurityGroups/myNSGName6ad2bc",
      "resourceGroup": "myResourceGroup6ad2bc"
    },
    "nicType": "Standard",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup6ad2bc",
    "tapConfigurations": [],
    "type": "Microsoft.Network/networkInterfaces",
    "vnetEncryptionSupported": false
  }
}
```

```bash
az network nic ip-config create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --nic-name $MY_VM_NIC_NAME \
    --name ipconfig1 \
    --private-ip-address-version IPv4 \
    --subnet $MY_SN_NAME \
    --vnet-name $MY_VNET_NAME \
    --make-primary true \
    --public-ip-address $MY_PUBLIC_IP_NAME
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/networkInterfaces/myVMNicName6ad2bc/ipConfigurations/ipconfig1",
  "name": "ipconfig1",
  "primary": true,
  "privateIPAddress": "10.19.0.4",
  "privateIPAddressVersion": "IPv4",
  "privateIPAllocationMethod": "Dynamic",
  "provisioningState": "Succeeded",
  "publicIPAddress": {
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/publicIPAddresses/myPublicIP6ad2bc",
    "resourceGroup": "myResourceGroup6ad2bc"
  },
  "resourceGroup": "myResourceGroup6ad2bc",
  "subnet": {
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/virtualNetworks/myVNet6ad2bc/subnets/mySN6ad2bc",
    "resourceGroup": "myResourceGroup6ad2bc"
  },
  "type": "Microsoft.Network/networkInterfaces/ipConfigurations"
}
```

# Cloud-init

```bash
cat << EOF > cloud-init.txt
#cloud-config

# Install, update, and upgrade packages
package_upgrade: true
package_update: true
package_reboot_if_require: true

# Install packages
packages:
  - vim
  - certbot
  - python3-certbot-nginx
  - bash-completion
  - nginx
  - mysql-client
  - php
  - php-cli
  - php-bcmath
  - php-curl
  - php-imagick
  - php-intl
  - php-json
  - php-mbstring
  - php-mysql
  - php-gd
  - php-xml
  - php-xmlrpc
  - php-zip
  - php-fpm

write_files:
  - owner: www-data:www-data
    path: /etc/nginx/sites-available/default.conf
    content: |
        server {
            listen 80 default_server;
            listen [::]:80 default_server;
            root /var/www/html;
            server_name $FQDN;
        }

write_files:
  - owner: www-data:www-data
    path: /etc/nginx/sites-available/$FQDN.conf
    content: |
        upstream php {
            server unix:/run/php/php8.1-fpm.sock;
        }
        server {
            listen 443 ssl http2;
            listen [::]:443 ssl http2;

            server_name $FQDN;

            ssl_certificate /etc/letsencrypt/live/$FQDN/fullchain.pem;
            ssl_certificate_key /etc/letsencrypt/live/$FQDN/privkey.pem;

            root /var/www/$FQDN;
            index index.php;

            location / {
                try_files \$uri \$uri/ /index.php?\$args;
            }
            location ~ \.php$ {
                include fastcgi_params;
                fastcgi_intercept_errors on;
                fastcgi_pass php;
                fastcgi_param  SCRIPT_FILENAME \$document_root\$fastcgi_script_name;
            }
            location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
                    expires max;
                    log_not_found off;
            }
            location = /favicon.ico {
                    log_not_found off;
                    access_log off;
            }

            location = /robots.txt {
                    allow all;
                    log_not_found off;
                    access_log off;
            }
        }
        server {
            listen 80;
            listen [::]:80;
            server_name $FQDN;
            return 301 https://$FQDN\$request_uri;
        }

runcmd:
  - sed -i "s/;cgi.fix_pathinfo.*/cgi.fix_pathinfo = 1/" /etc/php/8.1/fpm/php.ini
  - sed -i 's/^max_execution_time \= .*/max_execution_time \= 300/g' /etc/php/8.1/fpm/php.ini
  - sed -i 's/^upload_max_filesize \= .*/upload_max_filesize \= 64M/g' /etc/php/8.1/fpm/php.ini
  - sed -i 's/^post_max_size \= .*/post_max_size \= 64M/g' /etc/php/8.1/fpm/php.ini
  - systemctl restart php8.1-fpm
  - systemctl restart nginx
  - certbot --nginx certonly --non-interactive --agree-tos -d $FQDN -m dummy@dummy.com --redirect
  - ln -s /etc/nginx/sites-available/$FQDN.conf /etc/nginx/sites-enabled/
  - rm /etc/nginx/sites-enabled/default
  - systemctl restart nginx
  - curl --url https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar --output /tmp/wp-cli.phar
  - mv /tmp/wp-cli.phar /usr/local/bin/wp
  - chmod +x /usr/local/bin/wp
  - wp cli update
  - mkdir -m 0755 -p /var/www/$FQDN
  - chown -R azureadmin:www-data /var/www/$FQDN
  - sudo -u azureadmin -i -- wp core download --path=/var/www/$FQDN
  - sudo -u azureadmin -i -- wp config create --dbhost=$MY_MYSQL_DB_NAME.mysql.database.azure.com --dbname=wp001 --dbuser=$MY_MYSQL_ADMIN_USERNAME --dbpass="$MY_MYSQL_ADMIN_PW" --path=/var/www/$FQDN
  - sudo -u azureadmin -i -- wp core install --url=$FQDN --title="Azure hosted blog" --admin_user=$MY_WP_ADMIN_USER --admin_password="$MY_WP_ADMIN_PW" --admin_email=$MY_SSL_EMAIL --path=/var/www/$FQDN 
  - sudo -u azureadmin -i -- wp plugin update --all --path=/var/www/$FQDN
  - chmod 600 /var/www/$FQDN/wp-config.php
  - mkdir -p -m 0775 /var/www/$FQDN/wp-content/uploads
  - chgrp www-data /var/www/$FQDN/wp-content/uploads
EOF
```

# create private dns zone
```bash
az network private-dns zone create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_DNS_LABEL.private.mysql.database.azure.com
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myresourcegroup6ad2bc/providers/Microsoft.Network/privateDnsZones/mydnslabel6ad2bc.private.mysql.database.azure.com",
  "location": "global",
  "maxNumberOfRecordSets": 25000,
  "maxNumberOfVirtualNetworkLinks": 1000,
  "maxNumberOfVirtualNetworkLinksWithRegistration": 100,
  "name": "mydnslabel6ad2bc.private.mysql.database.azure.com",
  "numberOfRecordSets": 1,
  "numberOfVirtualNetworkLinks": 0,
  "numberOfVirtualNetworkLinksWithRegistration": 0,
  "provisioningState": "Succeeded",
  "resourceGroup": "myresourcegroup6ad2bc",
  "tags": null,
  "type": "Microsoft.Network/privateDnsZones"
}
```

# Create Azure MySQL Flexible Server

```bash
echo "Your MySQL user $MY_MYSQL_ADMIN_USERNAME password is: $MY_WP_ADMIN_PW" 
```

```bash
az mysql flexible-server create \
    --admin-password $MY_MYSQL_ADMIN_PW \
    --admin-user $MY_MYSQL_ADMIN_USERNAME \
    --auto-scale-iops Disabled \
    --high-availability Disabled \
    --iops 500 \
    --location $MY_LOCATION \
    --name $MY_MYSQL_DB_NAME \
    --database-name wp001 \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --sku-name Standard_B2s \
    --storage-auto-grow Disabled \
    --storage-size 20 \
    --subnet $MY_MYSQL_SN_NAME \
    --private-dns-zone $MY_DNS_LABEL.private.mysql.database.azure.com \
    --tier Burstable \
    --version 8.0.21 \
    --vnet $MY_VNET_NAME \
    --yes -o JSON
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "databaseName": "wp001",
  "host": "mydb6ad2bc.mysql.database.azure.com",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.DBforMySQL/flexibleServers/mydb6ad2bc",
  "location": "East US",
  "resourceGroup": "myResourceGroup6ad2bc",
  "skuname": "Standard_B2s",
  "subnetId": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Network/virtualNetworks/myVNet6ad2bc/subnets/myMySQLSN6ad2bc",
  "username": "dbadmin6ad2bc",
  "version": "8.0.21"
}
```

## Disable Azure MySQL Flexible Server SSL connection for Wordpress
```bash
az mysql flexible-server parameter set \
    -g $MY_RESOURCE_GROUP_NAME \
    -s $MY_MYSQL_DB_NAME \
    -n require_secure_transport -v "OFF" -o JSON
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "allowedValues": "ON,OFF",
  "currentValue": "OFF",
  "dataType": "Enumeration",
  "defaultValue": "ON",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.DBforMySQL/flexibleServers/mydb6ad2bc/configurations/require_secure_transport",
  "isConfigPendingRestart": "False",
  "isDynamicConfig": "True",
  "isReadOnly": "False",
  "name": "require_secure_transport",
  "resourceGroup": "myResourceGroup6ad2bc",
  "source": "user-override",
  "systemData": null,
  "type": "Microsoft.DBforMySQL/flexibleServers/configurations",
  "value": "OFF"
}
```

# Create VM

```bash
az vm create \
    --name $MY_VM_NAME \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --admin-username $MY_VM_USERNAME \
    --authentication-type ssh \
    --image $MY_VM_IMAGE \
    --location $MY_LOCATION \
    --nic-delete-option Delete \
    --os-disk-caching ReadOnly \
    --os-disk-delete-option Delete \
    --os-disk-size-gb 30 \
    --size $MY_VM_SIZE \
    --generate-ssh-keys \
    --storage-sku Premium_LRS \
    --nics $MY_VM_NIC_NAME \
    --custom-data cloud-init.txt 
```
Results:

<!-- expected_similarity=0.3 -->
```JSON
{
  "fqdns": "mydnslabel6ad2bc.eastus.cloudapp.azure.com",
  "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup6ad2bc/providers/Microsoft.Compute/virtualMachines/myVMName6ad2bc",
  "identity": {
    "principalId": "283f71ef-c1e1-4908-a4a5-fbed5a6fcaa9",
    "tenantId": "bd7153ee-d085-4a28-a928-2f0ef402f076",
    "type": "SystemAssigned",
    "userAssignedIdentities": null
  },
  "location": "eastus",
  "macAddress": "60-45-BD-D8-1D-84",
  "powerState": "VM running",
  "privateIpAddress": "10.19.0.4",
  "resourceGroup": "myResourceGroup6ad2bc",
  "zones": ""
}
```

# VM status
```bash
az vm show \
  --resource-group $MY_RESOURCE_GROUP_NAME \
  --name $MY_VM_NAME \
  --query provisioningState -o tsv
```
Results:

<!-- expected_similarity=0.3 -->
```ASCII
Succeeded
```

Validate that the application is running by visiting the application url:

```bash
curl --max-time 120 "https://$FQDN"
```
Results:

<!-- expected_similarity=0.3 -->
```HTML
<!DOCTYPE html>
<html lang="en-US">
<head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
<meta name='robots' content='max-image-preview:large' />
<title>Azure hosted blog</title>
<link rel="alternate" type="application/rss+xml" title="Azure hosted blog &raquo; Feed" href="https://mydnslabel6ad2bc.eastus.cloudapp.azure.com/?feed=rss2" />
<link rel="alternate" type="application/rss+xml" title="Azure hosted blog &raquo; Comments Feed" href="https://mydnslabel6ad2bc.eastus.cloudapp.azure.com/?feed=comments-rss2" />
```
