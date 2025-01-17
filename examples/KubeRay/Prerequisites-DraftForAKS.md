### Draft for Azure Kubernetes Service (AKS)

[Draft](https://github.com/Azure/draft) is an open-source project that streamlines Kubernetes development by taking a non-containerized application and generating the Dockerfiles, Kubernetes manifests, Helm charts, Kustomize configurations, and other artifacts associated with a containerized application. 

```bash
if ! command -v draft &> /dev/null
then
  echo "Draft could not be found, installing..."
  curl -fsSL https://raw.githubusercontent.com/Azure/draft/main/scripts/install.sh | bash
fi

echo "Draft is installed."
```

<!-- expected_similarity=".*installed" -->
```text
Draft is installed.
```

For more details on installing Draft see [Azure Kubernetes Service Preview extension](/azure/aks/draft#install-the-aks-preview-azure-cli-extension).