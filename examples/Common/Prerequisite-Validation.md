# VM SKU Availability

This document provides a script to test if the selected VM SKU is available in the selected region.

It is expected that the environment has variables that define the desired VM SKU and region. 

```bash
echo "Region selected is '$LOCATION'"
echo "VM SKU requested is '$VM_SKU'"
```

<!-- expected-similarity=0.8 -->
```text
Region selected is 'LOCATION'
VM SKU requested is 'VM_SKU'
```

# Check SKU Availabiltiy

```bash
az vm list-sizes --location $LOCATION --query "[?name=='$VM_SKU']" --output table
if [ $? -ne 0 ]; then
    echo "The requested VM SKU is not available in the selected region."
else
    echo "The requested VM SKU is available in the selected region."
fi
```

<!-- expected_results=1.0 -->
```text
The requested VM SKU is available in the selected region.
```

## Check VM SKU Availability

While it is a time consuming process testing in advance if a VM SKU is available in the selected region can reduce frustration later.

```bash
REASON_CODE=$(az vm list-skus --location $LOCATION --query "[?name=='$VM_SKU'].restrictions[].reasonCode" --output tsv)

if [ -z "$REASON_CODE" ]; then
    echo "VM SKU ('$VM_SKU') is available."
else
    echo "VM SKU ('$VM_SKU') is not available: $REASON_CODE"
fi
```

If the VM SKU is available this code block will output:

<!-- expected_similarity=".*is available." -->
```text
VM SKU ('VM_SKU') is available
```

With a little imagination it would be possible to create a validation script that updates the validaton script to automatically select an alternative SKU or region if the requested SKU was nopt available. For example `az vm list-skus --location eastus --query "[?capabilities[?name=='vCPUs' && to_number(value)>=2]]" --output table` will return all the SKUs that have 2 vCPUs.