# Step CA tutorial
This is a tutorial with A sample implementation of a PKI with a standalone Certificate Authority with step-ca on Azure, leveraging Azure Key Vault, Azure MySQL (soon) and Managed Identities. 

The first step in the tutorial is to create environment variables. The following will declare variable defaults if there are not some already set.
```bash
[[ -z "${AZURE_RG_NAME}" ]] && export AZURE_RG_NAME='pki'
[[ -z "${AZURE_LOCATION}" ]] && export AZURE_LOCATION='westeurope'
[[ -z "${AZURE_DNS_RESOLVER_OUTBOUND_TARGET_DNS}" ]] && export AZURE_DNS_RESOLVER_OUTBOUND_TARGET_DNS="[{\"ipAddress\": \"192.168.0.11\"},{\"ipAddress\": \"192.168.0.13\"}]"
[[ -z "${AZURE_DNS_RESOLVER_OUTBOUND_DOMAIN}" ]] && export AZURE_DNS_RESOLVER_OUTBOUND_DOMAIN='test.com.'
[[ -z "${CA_INIT_DNS}" ]] && export CA_INIT_DNS='azuredns.com'
[[ -z "${CA_SSH_PUBLIC_KEY}" ]] && export CA_SSH_PUBLIC_KEY="$(cat ~/.ssh/id_rsa.pub)"
[[ -z "${DB_ADMIN_PASSWORD}" ]] && export DB_ADMIN_PASSWORD='myPassword12345!'
[[ -z "${CA_INIT_NAME}" ]] && export CA_INIT_NAME='your CA Name'
[[ -z "${CA_INIT_DNS}" ]] && export CA_INIT_DNS='your DNS fqdn'
[[ -z "${CA_INIT_PROVISIONER_JWT}" ]] && export CA_INIT_PROVISIONER_JWT="$(az account show -o tsv --query user.name)"
```

Create Resource Group with the following command
```bash
az group create --name $AZURE_RG_NAME --location $AZURE_LOCATION -o none
```

```bash
az deployment group create --resource-group $AZURE_RG_NAME -o none \
  --template-file infra/base/step-ca-infra.bicep \
  --parameters caVMName="$CA_CAVMNAME" \
  --parameters keyvaultName="$CA_KEYVAULTNAME" \
  --parameters caVMPublicSshKey="$CA_SSH_PUBLIC_KEY" \
  --parameters ca_INIT_PROVISIONER_JWT="$CA_INIT_PROVISIONER_JWT" \
  --parameters ca_INIT_PASSWORD="$CA_INIT_PASSWORD" \
  --parameters ca_INIT_NAME="$CA_INIT_NAME" \
  --parameters ca_INIT_DNS="$CA_INIT_DNS" \
  --parameters dbLoginPassword="$DB_ADMIN_PASSWORD" \
  --parameters dnsResolverOutboundTargetDNS="$AZURE_DNS_RESOLVER_OUTBOUND_TARGET_DNS" \
  --parameters dnsResolverOutboundDNSDomainName="$AZURE_DNS_RESOLVER_OUTBOUND_DOMAIN"
```

## Connect to Private CA

Store necessary parameters
```bash
export AZURE_BASTION=$(az deployment group list -g pki -o tsv --query [0].properties.parameters.bastionName.value)
export CA_ADMIN_NAME=$(az deployment group list -g pki -o tsv --query [0].properties.parameters.caVMAdminUsername.value)
export CA_VM_NAME=$(az deployment group list -g pki -o tsv --query [0].properties.parameters.caVMName.value)
```

Connect Via Azure Bastion
```bash
az network bastion ssh -n $AZURE_BASTION -g $AZURE_RG_NAME \
  --auth-type ssh-key --username $CA_ADMIN_NAME --ssh-key ~/.ssh/id_rsa \
  --target-resource-id $(az vm show -g $AZURE_RG_NAME --name $CA_VM_NAME -o tsv --query id)
```