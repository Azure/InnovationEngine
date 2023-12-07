---
title: "Tutorial: Secure a web server with TLS/SSL certificates"
description: In this tutorial, you learn how to use the Azure CLI to secure a Linux virtual machine that runs the NGINX web server with SSL certificates stored in Azure Key Vault.
author: mattmcinnes
ms.service: virtual-machines
ms.collection: linux
ms.topic: tutorial
ms.date: 04/09/2023
ms.author: mattmcinnes
ms.custom: innovation-engine, mvc, devx-track-azurecli, GGAL-freshness822, devx-track-linux
#Customer intent: As an IT administrator or developer, I want to learn how to secure a web server with TLS/SSL certificates so that I can protect my customer data on web applications that I build and run.
ms.permissions: Microsoft.Network/networkInterfaces/write, Microsoft.Network/virtualNetworks/subnets/read, Microsoft.Network/virtualNetworks/subnets/join/action, Microsoft.Compute/virtualMachines/write, Microsoft.Network/virtualNetworks/read, Microsoft.KeyVault/vaults/listVersions/actionMicrosoft.Compute/virtualMachines/read Microsoft.KeyVault/vaults/read, Microsoft.Resources/subscriptions/resourcegroups/read, Microsoft.Resources/subscriptions/resourcegroups/write, Microsoft.KeyVault/vaults/write, Microsoft.KeyVault/certificates/create, Microsoft.Network/networkSecurityGroups/write
---

# Tutorial: Use TLS/SSL certificates to secure a web server


**Applies to:** :heavy_check_mark: Linux VMs 

To secure web servers, a Transport Layer Security (TLS), previously known as Secure Sockets Layer (SSL), certificate can be used to encrypt web traffic. These TLS/SSL certificates can be stored in Azure Key Vault, and allow secure deployments of certificates to Linux virtual machines (VMs) in Azure. In this tutorial you learn how to:

> [!div class="checklist"]
> * Create an Azure Key Vault
> * Generate or upload a certificate to the Key Vault
> * Create a VM and install the NGINX web server
> * Inject the certificate into the VM and configure NGINX with a TLS binding

This tutorial uses the CLI within the [Azure Cloud Shell](../../cloud-shell/overview.md), which is constantly updated to the latest version. To open the Cloud Shell, select **Try it** from the top of any code block.

If you choose to install and use the CLI locally, this tutorial requires that you're running the Azure CLI version 2.0.30 or later. Run `az --version` to find the version. If you need to install or upgrade, see [Install Azure CLI]( /cli/azure/install-azure-cli).



## Define Environment Variables

The First step in this tutorial is to define environment variables.

```bash
export RANDOM_ID="$(openssl rand -hex 3)"
export MyName1=myResourceGroupSecureWeb$RANDOM_ID
export MyName2=$keyvault_name$RANDOM_ID
export MyName3=mycert$RANDOM_ID
export MyName4=myVM$RANDOM_ID
export MyLocation=eastus
export MyResourceGroup=myResourceGroupSecureWeb$RANDOM_ID
export MyEnabledForDeployment=```
export MyVaultName=$keyvault_name$RANDOM_ID
export MyPolicy="$(az
export MyQuery="[?attributes.enabled].id"
export MyOutput=tsv)
export MySecrets1="$secret"
export MySecrets2="$vm_secret"
export MyKeyvault=$keyvault_name)
export MyImage=Ubuntu2204
export MyAdminUsername=azureuser$RANDOM_ID
export MyGenerateSshKeys=\
export MyCustomData=cloud-init-web-server.txt
export MyPort=443
```
## Overview
Azure Key Vault safeguards cryptographic keys and secrets, such as certificates or passwords. Key Vault helps streamline the certificate management process and enables you to maintain control of keys that access those certificates. You can create a self-signed certificate inside Key Vault, or upload an existing, trusted certificate that you already own.

Rather than using a custom VM image that includes certificates baked-in, you inject certificates into a running VM. This process ensures that the most up-to-date certificates are installed on a web server during deployment. If you renew or replace a certificate, you don't also have to create a new custom VM image. The latest certificates are automatically injected as you create more VMs. During the whole process, the certificates never leave the Azure platform or are exposed in a script, command-line history, or template.


## Create an Azure Key Vault
Before you can create a Key Vault and certificates, create a resource group with [az group create](/cli/azure/group). The following example creates a resource group named *myResourceGroupSecureWeb* in the *eastus* location:

```azurecli-interactive 
az group create --name $MyName1 --location $MyLocation
```

Results:

<!-- expected_similarity=0.3 -->
```json
{
    "id": "/subscriptions/{subscriptionId}/resourceGroups/myResourceGroupSecureWeb",
    "location": "eastus",
    "managedBy": null,
    "name": "myResourceGroupSecureWeb",
    "properties": {
        "provisioningState": "Succeeded"
    },
    "tags": null,
    "type": "Microsoft.Resources/resourceGroups"
}
```

Next, create a Key Vault with [az keyvault create](/cli/azure/keyvault) and enable it for use when you deploy a VM. Each Key Vault requires a unique name, and should be all lowercase. Replace *\<mykeyvault>* in the following example with your own unique Key Vault name:

```azurecli-interactive 
keyvault_name=<mykeyvault>
az keyvault create \
    --resource-group $MyName1 \
    --name $keyvault_name \
    --enabled-for-deployment
```

Results:

<!-- expected_similarity=0.3 -->
```json
{
    "id": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroupSecureWeb/providers/Microsoft.KeyVault/vaults/mykeyvault",
    "location": "centralus",
    "name": "mykeyvault",
    "properties": {
        "accessPolicies": [],
        "createMode": null,
        "enabledForDeployment": true,
        "enabledForDiskEncryption": null,
        "enabledForTemplateDeployment": null,
        "enableSoftDelete": null,
        "enableRbacAuthorization": null,
        "networkAcls": null,
        "privateEndpointConnections": null,
        "provisioningState": "InProgress",
        "sku": {
            "family": null,
            "name": "standard"
        },
        "tenantId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
        "vaultUri": "https://mykeyvault.vault.azure.net/"
    },
    "type": "Microsoft.KeyVault/vaults"
}
```

## Generate a certificate and store in Key Vault
For production use, you should import a valid certificate signed by trusted provider with [az keyvault certificate import](/cli/azure/keyvault/certificate). For this tutorial, the following example shows how you can generate a self-signed certificate with [az keyvault certificate create](/cli/azure/keyvault/certificate) that uses the default certificate policy:

```azurecli-interactive 
az keyvault certificate create \
    --vault-name $keyvault_name \
    --name mycert \
    --policy "$(az keyvault certificate get-default-policy)"
```

### Prepare a certificate for use with a VM
To use the certificate during the VM create process, obtain the ID of your certificate with [az keyvault secret list-versions](/cli/azure/keyvault/secret). Convert the certificate with [az vm secret format](/cli/azure/vm/secret#az-vm-secret-format). The following example assigns the output of these commands to variables for ease of use in the next steps:

```azurecli-interactive 
secret=$(az keyvault secret list-versions \
          --vault-name $keyvault_name \
          --name $MyName3 \
          --query "[?attributes.enabled].id" --output tsv)
vm_secret=$(az vm secret format --secrets "$secret" -g $MyName1 --keyvault $keyvault_name)
```

Results:

<!-- expected_similarity=0.3 -->
```json
{
  "result": "$(az keyvault secret list-versions --vault-name $keyvault_name --name mycert --query '[?attributes.enabled].id' --output tsv)",
  "vm_secret": "$(az vm secret format --secrets "$secret" -g myResourceGroupSecureWeb --keyvault $keyvault_name)"
}
```

### Create a cloud-init config to secure NGINX
[Cloud-init](https://cloudinit.readthedocs.io) is a widely used approach to customize a Linux VM as it boots for the first time. You can use cloud-init to install packages and write files, or to configure users and security. As cloud-init runs during the initial boot process, there are no extra steps or required agents to apply your configuration.

When you create a VM, certificates and keys are stored in the protected */var/lib/waagent/* directory. To automate adding the certificate to the VM and configuring the web server, use cloud-init. In this example, you install and configure the NGINX web server. You can use the same process to install and configure Apache. 

Create a file named *cloud-init-web-server.txt* and paste the following configuration:

```yaml
#cloud-config
package_upgrade: true
packages:
  - nginx
write_files:
  - owner: www-data:www-data
  - path: /etc/nginx/sites-available/default
    content: |
      server {
        listen 443 ssl;
        ssl_certificate /etc/nginx/ssl/mycert.cert;
        ssl_certificate_key /etc/nginx/ssl/mycert.prv;
      }
runcmd:
  - secretsname=$(find /var/lib/waagent/ -name "*.prv" | cut -c -57)
  - mkdir /etc/nginx/ssl
  - cp $secretsname.crt /etc/nginx/ssl/mycert.cert
  - cp $secretsname.prv /etc/nginx/ssl/mycert.prv
  - service nginx restart
```

### Create a secure VM
Now create a VM with [az vm create](/cli/azure/vm). The certificate data is injected from Key Vault with the `--secrets` parameter. You pass in the cloud-init config with the `--custom-data` parameter:

```azurecli-interactive 
az vm create \
    --resource-group $MyName1 \
    --name $MyName4 \
    --image $MyImage \
    --admin-username $MyAdminUsername \
    --generate-ssh-keys \
    --custom-data $MyCustomData \
    --secrets "$vm_secret"
```

Results:

<!-- expected_similarity=0.3 -->
```json
{
    "status": "InteractiveSessionStarted",
    "message": "Started interactive session with VM creation command"
}
```

It takes a few minutes for the VM to be created, the packages to install, and the app to start. When the VM has been created, take note of the `publicIpAddress` displayed by the Azure CLI. This address is used to access your site in a web browser.

To allow secure web traffic to reach your VM, open port 443 from the Internet with [az vm open-port](/cli/azure/vm):

```azurecli-interactive 
az vm open-port \
    --resource-group myResourceGroupSecureWeb \
    --name myVM \
    --port 443
```


### Test the secure web app
Now you can open a web browser and enter *https:\/\/\<publicIpAddress>* in the address bar. Provide your own public IP address from the VM create process. Accept the security warning if you used a self-signed certificate:

![Accept web browser security warning](./media/tutorial-secure-web-server/browser-warning.png)

Your secured NGINX site is then displayed as in the following example:

![View running secure NGINX site](./media/tutorial-secure-web-server/secured-nginx.png)


## Next steps

In this tutorial, you secured an NGINX web server with a TLS/SSL certificate stored in Azure Key Vault. You learned how to:

> [!div class="checklist"]
> * Create an Azure Key Vault
> * Generate or upload a certificate to the Key Vault
> * Create a VM and install the NGINX web server
> * Inject the certificate into the VM and configure NGINX with a TLS binding

Follow this link to see pre-built virtual machine script samples.

> [!div class="nextstepaction"]
> [Linux virtual machine script samples](https://github.com/Azure-Samples/azure-cli-samples/tree/master/virtual-machine)

<details>
<summary><h2>FAQs</h2></summary>

#### Q. What is the command-specific breakdown of permissions needed to implement this doc? 

A. _Format: Commands as they appears in the doc | list of unique permissions needed to run each of those commands_


  - ```azurecli-interactive az group create --name $MyName1 --location $MyLocation ```

      - Microsoft.Resources/subscriptions/resourcegroups/read
      - Microsoft.Resources/subscriptions/resourcegroups/write
  - ```azurecli-interactive keyvault_name=<mykeyvault> az keyvault create \ --resource-group $MyName1 \ --name $keyvault_name \ --enabled-for-deployment ```

      - Microsoft.KeyVault/vaults/write
  - ```azurecli-interactive az keyvault certificate create \ --vault-name $keyvault_name \ --name mycert \ --policy "$(az keyvault certificate get-default-policy)" ```

      - Microsoft.KeyVault/certificates/create
  - ```azurecli-interactive secret=$(az keyvault secret list-versions \ --vault-name $keyvault_name \ --name $MyName3 \ --query "[?attributes.enabled].id" --output tsv) vm_secret=$(az vm secret format --secrets "$secret" -g $MyName1 --keyvault $keyvault_name) ```

      - Microsoft.KeyVault/vaults/listVersions/actionMicrosoft.Compute/virtualMachines/read Microsoft.KeyVault/vaults/read
  - ```azurecli-interactive az vm create \ --resource-group $MyName1 \ --name $MyName4 \ --image $MyImage \ --admin-username $MyAdminUsername \ --generate-ssh-keys \ --custom-data $MyCustomData \ --secrets "$vm_secret" ```

      - Microsoft.Network/networkInterfaces/write
      - Microsoft.Network/virtualNetworks/subnets/join/action
      - Microsoft.Compute/virtualMachines/write
      - Microsoft.Network/virtualNetworks/read
      - Microsoft.Network/virtualNetworks/subnets/read
  - ```azurecli-interactive az vm open-port \ --resource-group myResourceGroupSecureWeb \ --name myVM \ --port 443 ```

      - Microsoft.Network/virtualNetworks/subnets/join/action
      - Microsoft.Network/networkSecurityGroups/write
      - Microsoft.Compute/virtualMachines/write

#### Q. What is the purpose of using TLS/SSL certificates to secure a web server? 

A. TLS/SSL certificates are used to encrypt web traffic and protect customer data on web applications. They ensure secure deployments of certificates to Linux virtual machines in Azure, allowing for the most up-to-date certificates to be installed on web servers during deployment.


#### Q. How do I create an Azure Key Vault? 

A. To create an Azure Key Vault, you need to create a resource group first using the 'az group create' command. Then, use the 'az keyvault create' command to create the Key Vault, providing a unique name for the Key Vault and enabling it for deployment. You can find detailed steps and examples in the Azure CLI documentation: [Create an Azure Key Vault](https://docs.microsoft.com/azure/key-vault/quick-create-cli).


#### Q. What are the options for generating or uploading a certificate to the Key Vault? 

A. You have two options for generating or uploading a certificate to Azure Key Vault. The first option is to import a valid certificate signed by a trusted provider using the 'az keyvault certificate import' command. The second option is to generate a self-signed certificate using the 'az keyvault certificate create' command. Both options are explained in detail in the Azure CLI documentation: [Import a certificate into Azure Key Vault](https://docs.microsoft.com/azure/key-vault/certificates/import-certificate-cli) and [Create a self-signed certificate in Azure Key Vault](https://docs.microsoft.com/azure/key-vault/certificates/create-cert-cli).


#### Q. How do I prepare a certificate for use with a VM? 

A. To use a certificate during the VM create process, you need to obtain the ID of the certificate using the 'az keyvault secret list-versions' command. Then, convert the certificate using the 'az vm secret format' command. Both commands can be combined to assign the output to variables for ease of use in the next steps. You can find detailed steps and examples in the Azure CLI documentation: [Prepare a certificate for use with a VM](https://docs.microsoft.com/azure/virtual-machines/scripts/virtual-machines-linux-cli-sample-create-vm-disk-encrypt?key=vhd){lease add proper link}


#### Q. What is cloud-init and how can I use it to configure a Linux VM? 

A. Cloud-init is a widely used approach to customize a Linux VM as it boots for the first time. It allows you to install packages, write files, configure users, and perform other tasks during the initial boot process. You can use cloud-init to automate adding a certificate to a VM and configuring the web server. The tutorial provides an example of creating a cloud-init configuration file to install and configure the NGINX web server. You can find more information about cloud-init in the [Cloud-init documentation](https://cloudinit.readthedocs.io).


#### Q. How do I create a secure VM with the injected certificate and configured NGINX web server? 

A. To create a secure VM with the injected certificate and configured NGINX web server, you need to use the 'az vm create' command. Pass in the necessary parameters such as resource group, VM name, image, admin username, SSH keys, custom data (cloud-init config), and the secrets from Key Vault using the '--secrets' parameter. The tutorial provides a detailed example of creating a secure VM with the injected certificate and configured NGINX web server. You can find more information about creating VMs with the Azure CLI in the [Azure CLI documentation](https://docs.microsoft.com/azure/virtual-machines/scripts/virtual-machines-linux-cli-sample-create-vm-disk-encrypt).


#### Q. How do I open port 443 to allow secure web traffic to reach my VM? 

A. To open port 443 and allow secure web traffic to reach your VM, you can use the 'az vm open-port' command. Specify the resource group, VM name, and the port to be opened (in this case, port 443). The tutorial provides an example of opening port 443 using the Azure CLI. You can find more information about opening ports for VMs in the [Azure CLI documentation](https://docs.microsoft.com/azure/virtual-machines/scripts/virtual-machines-linux-cli-sample-create-availability-set).


#### Q. How do I test the secure web app running on my VM? 

A. To test the secure web app running on your VM, you can open a web browser and enter 'https://<publicIpAddress>' in the address bar. Replace '<publicIpAddress>' with the public IP address of your VM, which you can find in the Azure CLI output after creating the VM. Accept any security warnings if you used a self-signed certificate. The tutorial provides screenshots of the web browser and the secured NGINX site. You can find more information about accessing VMs in a web browser in the [Azure Virtual Machines documentation](https://docs.microsoft.com/azure/virtual-machines/windows/connect-ui-overview).


#### Q. What are the next steps after securing the NGINX web server with a TLS/SSL certificate? 

A. After securing the NGINX web server with a TLS/SSL certificate, you can explore pre-built virtual machine script samples for further customization and configuration. The tutorial provides a link to the [Azure CLI Samples GitHub repository](https://github.com/Azure-Samples/azure-cli-samples/tree/master/virtual-machine), where you can find a variety of script samples for Linux virtual machines.


#### Q. Which Azure CLI version is required to follow this tutorial? 

A. This tutorial requires that you're running the Azure CLI version 2.0.30 or later. You can check your current CLI version by running 'az --version'. If you need to install or upgrade the Azure CLI, you can find instructions in the [Azure CLI documentation](https://docs.microsoft.com/azure/virtual-machines/scripts/virtual-machines-linux-cli-sample-create-availability-set#install-the-azure-cli).

</details>