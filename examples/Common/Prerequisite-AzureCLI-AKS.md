### Check kubelogin and install if not exists

```bash
if ! command -v kubelogin &> /dev/null; then
  echo "kubelogin could not be found. Installing kubelogin..."
  az aks install-cli
fi
```
