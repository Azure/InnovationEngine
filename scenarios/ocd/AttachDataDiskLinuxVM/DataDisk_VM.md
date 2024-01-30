az login
az account set --subscription "0bb78609-cc8b-4e7d-be30-eee8cf2dbea4"
az account show

#VARIABLE DECLARATION
export NETWORK_PREFIX="$(($RANDOM % 254 + 1))"
export RANDOM_ID="$(openssl rand -hex 3)"
export MY_RESOURCE_GROUP_NAME="myResourceGroup$RANDOM_ID"
export REGION="eastus"
export MY_VM_NAME="myVMName$RANDOM_ID"
export MY_VM_IMAGE='Ubuntu2204'
export MY_VM_USERNAME="azureuser"
export MY_VM_SIZE='Standard_DS2_v5'
export MY_VNET_NAME="myVNet$RANDOM_ID"
export MY_VNET_PREFIX="10.$NETWORK_PREFIX.0.0/16"
export MY_VM_NIC_NAME="myVMNicName$RANDOM_ID"
export MY_SN_NAME="mySN$RANDOM_ID"
export MY_SN_PREFIX="10.$NETWORK_PREFIX.0.0/24"
export MY_PUBLIC_IP_NAME="myPublicIP$RANDOM_ID"
export MY_DNS_LABEL="mydnslabel$RANDOM_ID"
export MY_NSG_NAME="myNSGName$RANDOM_ID"

#CREATE A RESOURCE GROUP
az group create \
    --name $MY_RESOURCE_GROUP_NAME \
    --location $REGION -o JSON

#SET UP VM NETWORK
az network vnet create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_VNET_NAME \
    --location $REGION \
    --address-prefix $MY_VNET_PREFIX \
    --subnet-name $MY_SN_NAME \
    --subnet-prefix $MY_SN_PREFIX -o JSON

#CREATE STATIC PUBLIC IP 
az network public-ip create \
    --name $MY_PUBLIC_IP_NAME \
    --location $REGION \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --dns-name $MY_DNS_LABEL \
    --sku Standard \
    --allocation-method static \
    --version IPv4 \
    --zone 1 2 3 -o JSON

#CREATE NSG
az network nsg create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_NSG_NAME \
    --location $REGION -o JSON

#CREATE NSG RULES
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

#CREATE NIC
az network nic create \
    --resource-group $MY_RESOURCE_GROUP_NAME \
    --name $MY_VM_NIC_NAME \
    --location $REGION \
    --ip-forwarding false \
    --subnet $MY_SN_NAME \
    --vnet-name $MY_VNET_NAME \
    --network-security-group $MY_NSG_NAME \
    --public-ip-address $MY_PUBLIC_IP_NAME -o JSON

#CREATE LINUX VM
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