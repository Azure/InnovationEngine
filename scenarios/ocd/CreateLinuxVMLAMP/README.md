# Variable declaration

```bash
export UNIQUE_POSTFIX="$(($RANDOM % 254 + 1))"
export MY_RESOURCE_GROUP_NAME="myResourceGroup$UNIQUE_POSTFIX"
export MY_LOCATION="eastus"
export MY_VM_NAME="myVMName$UNIQUE_POSTFIX"
export MY_VM_USERNAME="azureadmin"
export MY_VM_SIZE='Standard_DS2_v2'
export MY_VM_IMAGE='Canonical:0001-com-ubuntu-minimal-jammy:minimal-22_04-lts-gen2:latest'
export MY_PUBLIC_IP_NAME="myPublicIP$UNIQUE_POSTFIX"
export MY_DNS_LABEL="mydnslabel$UNIQUE_POSTFIX"
export MY_NSG_NAME="myNSGName$UNIQUE_POSTFIX"
export MY_NSG_SSH_RULE="Allow-Access$UNIQUE_POSTFIX"
export MY_VM_NIC_NAME="myVMNicName$UNIQUE_POSTFIX"
export MY_VNET_NAME="myVNet$UNIQUE_POSTFIX"
export MY_VNET_PREFIX="10.$UNIQUE_POSTFIX.0.0/22"
export MY_SN_NAME="mySN$UNIQUE_POSTFIX"
export MY_SN_PREFIX="10.$UNIQUE_POSTFIX.0.0/24"
export MY_MYSQL_DB_NAME="myDB$UNIQUE_POSTFIX"
export MY_MYSQL_ADMIN_USERNAME="dbadmin$UNIQUE_POSTFIX"
export MY_MYSQL_ADMIN_PW="etregdgdfggg$UNIQUE_POSTFIX"
export MY_MYSQL_SN_NAME="myMySQLSN$UNIQUE_POSTFIX"
export MY_WP_ADMIN_PW="$(openssl rand -base64 32)"
export MY_WP_ADMIN_USER="wpcliadmin"
export FQDN="${MY_DNS_LABEL}.${MY_LOCATION}.cloudapp.azure.com"
```

# Create RG
```bash
az group create \
    --name $MY_RESOURCE_GROUP_NAME \
    --location $MY_LOCATION
```

Results:
```expected_similarity=0.3
{
  "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup104",
  "location": "eastus",
  "managedBy": null,
  "name": "myResourceGroup104",
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
```expected_similarity=0.3
{
  "newVNet": {
    "addressSpace": {
      "addressPrefixes": [
        "10.104.0.0/16"
      ]
    },
    "enableDdosProtection": false,
    "etag": "W/\"7859ca24-18f2-4569-8ecd-4bfdd84e355d\"",
    "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup104/providers/Microsoft.Network/virtualNetworks/myVNet104",
    "location": "eastus",
    "name": "myVNet104",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup104",
    "resourceGuid": "b4255e64-8b6a-4b98-b9dd-727c8b9a1ab7",
    "subnets": [
      {
        "addressPrefix": "10.104.0.0/22",
        "delegations": [],
        "etag": "W/\"7859ca24-18f2-4569-8ecd-4bfdd84e355d\"",
        "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup104/providers/Microsoft.Network/virtualNetworks/myVNet104/subnets/mySN104",
        "name": "mySN104",
        "privateEndpointNetworkPolicies": "Disabled",
        "privateLinkServiceNetworkPolicies": "Enabled",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup104",
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
```expected_similarity=0.3
{
  "publicIp": {
    "ddosSettings": {
      "protectionMode": "VirtualNetworkInherited"
    },
    "dnsSettings": {
      "domainNameLabel": "mydnslabel104",
      "fqdn": "mydnslabel104.eastus.cloudapp.azure.com"
    },
    "etag": "W/\"a48ca844-4aa0-4bd5-b0af-83262773ee30\"",
    "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup104/providers/Microsoft.Network/publicIPAddresses/myPublicIP104",
    "idleTimeoutInMinutes": 4,
    "ipAddress": "52.152.193.7",
    "ipTags": [],
    "location": "eastus",
    "name": "myPublicIP104",
    "provisioningState": "Succeeded",
    "publicIPAddressVersion": "IPv4",
    "publicIPAllocationMethod": "Static",
    "resourceGroup": "myResourceGroup104",
    "resourceGuid": "887cbd99-b430-405a-9ed9-2590336720dd",
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
```expected_similarity=0.3
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
        "etag": "W/\"f30a1bb7-d798-472c-9783-5da67d766ef5\"",
        "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup104/providers/Microsoft.Network/networkSecurityGroups/protect-vms/defaultSecurityRules/AllowVnetInBound",
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
    "etag": "W/\"f30a1bb7-d798-472c-9783-5da67d766ef5\"",
    "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup104/providers/Microsoft.Network/networkSecurityGroups/protect-vms",
    "location": "eastus",
    "name": "protect-vms",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup104",
    "resourceGuid": "c08deada-a6d7-4876-b3cf-777f05e33bcd",
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
```expected_similarity=0.3
{
  "access": "Allow",
  "destinationAddressPrefix": "*",
  "destinationAddressPrefixes": [],
  "destinationPortRange": "22",
  "destinationPortRanges": [],
  "direction": "Inbound",
  "etag": "W/\"f5b7f774-8dea-43c3-bf1a-daf1c6241e7b\"",
  "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup164/providers/Microsoft.Network/networkSecurityGroups/myNSGName164/securityRules/Allow-SSH164",
  "name": "Allow-SSH164",
  "priority": 100,
  "protocol": "Tcp",
  "provisioningState": "Succeeded",
  "resourceGroup": "myResourceGroup164",
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
```expected_similarity=0.3
{
  "NewNIC": {
    "auxiliaryMode": "None",
    "auxiliarySku": "None",
    "disableTcpStateTracking": false,
    "dnsSettings": {
      "appliedDnsServers": [],
      "dnsServers": [],
      "internalDomainNameSuffix": "3ftdfd3rkcnuxaci1mkks2kkoe.bx.internal.cloudapp.net"
    },
    "enableAcceleratedNetworking": false,
    "enableIPForwarding": false,
    "etag": "W/\"c7a45d71-cc10-4139-be7c-8f932e7c0ffb\"",
    "hostedWorkloads": [],
    "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup164/providers/Microsoft.Network/networkInterfaces/myVMNicName164",
    "ipConfigurations": [
      {
        "etag": "W/\"c7a45d71-cc10-4139-be7c-8f932e7c0ffb\"",
        "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup164/providers/Microsoft.Network/networkInterfaces/myVMNicName164/ipConfigurations/ipconfig1",
        "name": "ipconfig1",
        "primary": true,
        "privateIPAddress": "10.164.0.4",
        "privateIPAddressVersion": "IPv4",
        "privateIPAllocationMethod": "Dynamic",
        "provisioningState": "Succeeded",
        "resourceGroup": "myResourceGroup164",
        "subnet": {
          "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup164/providers/Microsoft.Network/virtualNetworks/myVNet164/subnets/mySN164",
          "resourceGroup": "myResourceGroup164"
        },
        "type": "Microsoft.Network/networkInterfaces/ipConfigurations"
      }
    ],
    "location": "eastus",
    "name": "myVMNicName164",
    "networkSecurityGroup": {
      "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup164/providers/Microsoft.Network/networkSecurityGroups/myNSGName164",
      "resourceGroup": "myResourceGroup164"
    },
    "nicType": "Standard",
    "provisioningState": "Succeeded",
    "resourceGroup": "myResourceGroup164",
    "resourceGuid": "a6e28735-1d4b-4ef5-87e4-f85ef0e41e30",
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
```expected_similarity=0.3
{
  "etag": "W/\"aedcac3a-9726-4cb0-8378-cf6a37bde1a0\"",
  "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup164/providers/Microsoft.Network/networkInterfaces/myVMNicName164/ipConfigurations/ipconfig1",
  "name": "ipconfig1",
  "primary": true,
  "privateIPAddress": "10.164.0.4",
  "privateIPAddressVersion": "IPv4",
  "privateIPAllocationMethod": "Dynamic",
  "provisioningState": "Succeeded",
  "publicIPAddress": {
    "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup164/providers/Microsoft.Network/publicIPAddresses/myPublicIP164",
    "resourceGroup": "myResourceGroup164"
  },
  "resourceGroup": "myResourceGroup164",
  "subnet": {
    "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup164/providers/Microsoft.Network/virtualNetworks/myVNet164/subnets/mySN164",
    "resourceGroup": "myResourceGroup164"
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
                try_files $uri $uri/ /index.php?$args;
            }
            location ~ \.php$ {
                include fastcgi_params;
                fastcgi_intercept_errors on;
                fastcgi_pass php;
                fastcgi_param  SCRIPT_FILENAME $document_root$fastcgi_script_name;
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
  - systemctl restart php8.1-fpm
  - systemctl restart nginx
  - certbot run -n --nginx --agree-tos -d $FQDN -m bla@bla.com --redirect
  - ln -s /etc/nginx/sites-available/$FQDN.conf /etc/nginx/sites-enabled/
  - systemctl restart nginx
  - curl --url https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar --output /tmp/wp-cli.phar
  - mv /tmp/wp-cli.phar /usr/local/bin/wp
  - chmod +x /usr/local/bin/wp
  - wp cli update
  - mkdir -m 0755 -p /var/www/$FQDN
  - chown -R azureadmin:www-data /var/www/$FQDN
  - sudo -u azureadmin -i -- wp core download --path=/var/www/$FQDN
  - sudo -u azureadmin -i -- wp core config --dbhost=mydb56.mysql.database.azure.com --dbname=wp001 --dbuser=$MY_MYSQL_ADMIN_USERNAME --dbpass=$MY_MYSQL_ADMIN_PW --path=/var/www/$FQDN
  - sudo -u azureadmin -i -- wp db create --path=/var/www/$FQDN
  - sudo -u azureadmin -i -- wp core install --url=$FQDN --title="Azure hosted blog" --admin_user=$MY_WP_ADMIN_USER --admin_password=$MY_WP_ADMIN_PW --admin_email=example@example.org --path=/var/www/$FQDN 
  - sudo -u azureadmin -i -- wp plugin update --all --path=/var/www/$FQDN
EOF
```

# create private dns zone
```bash
az network private-dns zone create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_DNS_LABEL.private.mysql.database.azure.com
```

Results:
```expected_similarity=0.3
{
  "etag": "e5e6b1b5-af49-4a34-9b55-bc62c2c4579f",
  "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myresourcegroup36/providers/Microsoft.Network/privateDnsZones/mydnslabel36.private.mysql.database.azure.com",
  "location": "global",
  "maxNumberOfRecordSets": 25000,
  "maxNumberOfVirtualNetworkLinks": 1000,
  "maxNumberOfVirtualNetworkLinksWithRegistration": 100,
  "name": "mydnslabel36.private.mysql.database.azure.com",
  "numberOfRecordSets": 1,
  "numberOfVirtualNetworkLinks": 0,
  "numberOfVirtualNetworkLinksWithRegistration": 0,
  "provisioningState": "Succeeded",
  "resourceGroup": "myresourcegroup36",
  "tags": null,
  "type": "Microsoft.Network/privateDnsZones"
}
```

# Create Azure MySQL Flexible Server

```bash
az mysql flexible-server create \
    --admin-password $MY_MYSQL_ADMIN_PW \
    --admin-user $MY_MYSQL_ADMIN_USERNAME \
    --auto-scale-iops Disabled \
    --high-availability Disabled \
    --iops 500 \
    --location $MY_LOCATION \
    --name $MY_MYSQL_DB_NAME \
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
```expected_similarity=0.3
{
  "connectionString": "mysql wp001 --host mydb56.mysql.database.azure.com --user dbadmin56 --password=etregdgdfggg56",
  "databaseName": "wp001",
  "host": "mydb56.mysql.database.azure.com",
  "id": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup56/providers/Microsoft.DBforMySQL/flexibleServers/mydb56",
  "location": "East US",
  "password": "etregdgdfggg56",
  "resourceGroup": "myResourceGroup56",
  "skuname": "Standard_B2s",
  "subnetId": "/subscriptions/7f9b0964-9093-4e26-b299-451fea2d435d/resourceGroups/myResourceGroup56/providers/Microsoft.Network/virtualNetworks/myVNet56/subnets/myMySQLSN56",
  "username": "dbadmin56",
  "version": "8.0.21"
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



## Then, replace the examples below with the info for your WordPress database:
wp core config create --dbhost=host.db --dbname=prefix_db --dbuser=username --dbpass=password

wp db create
## Change permissions for wp-config.php
chmod 600 wp-config.php

## Configure wp-config.php
wp core install --url=yourwebsite.com --title="Your Blog Title" --admin_name=wordpress_admin --admin_password=4Long&Strong1 --admin_email=you@example.com
wp plugin update --all

## Enable file uploads
cd wp-content
mkdir uploads
chgrp web uploads/
chmod 775 uploads/