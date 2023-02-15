<!-- 
```variables
export MY_RESOURCE_GROUP_NAME=testResourceGroup
export MY_LOCATION=eastus
export MY_VM_NAME=myVM
export MY_VM_IMAGE=debian
export MY_ADMIN_USERNAME=azureuser
```
-->

```bash
printenv | grep "^MY_"
```

<!--expected_similarity=.7-->
```output
MY_VM_IMAGE=debian
MY_RESOURCE_GROUP_NAME=testResourceGroup
MY_LOCATION=eastus
MY_VM_NAME=myVM
MY_ADMIN_USERNAME=azureuser
```