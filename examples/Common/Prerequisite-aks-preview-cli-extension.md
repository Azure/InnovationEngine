### Install or Update AKS Preview CLI Extension

[!INCLUDE [preview features callout](~/reusable-content/ce-skilling/azure/includes/aks/includes/preview/preview-callout.md)]

To install the aks-preview extension, run the following command:

```bash
AKS_PREVIWS_VERSION=${az extension list --query "[?name=='aks-preview'].version" --output tsv}

if [ -z "$AKS_PREVIWS_VERSION" ]; then
    az extension add --name aks-preview
else
    az extension update --name aks-preview
fi

echo "Latest AKS Preview CLI extension installed."
```

<!-- expected_similarity=1.0 -->
```text
Latest AKS Preview CLI extension installed.yo
```
