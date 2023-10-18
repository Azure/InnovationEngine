# Welcome to the innovation Engine Tutorial
## *TODO ADD MORE DETAIL TO IMPROVE TUTORIAL*

# Running simple bash commands

Innovation engine can execute bash commands. For example

```bash
echo "Hello World"
```

# Test Code block with expected output

```azurecli-interactive
echo "Hello \
world"
```

It also can test the output to make sure everything ran as planned.
<!--expected_similarity=0.8-->
```
Hello world
```

# Test Code block matches expected regex

```bash
echo "Foo Bar"
```

It also can test the output to make sure everything ran as planned.
<!--expected_similarity="Foo \w+"-->
```
Foo Bar
```

# Executable vs non-executable code blocks
Innovation engine supports code blocks which are both executable and non-executable. A code block is executable if the label/tag after the bash scripts is one of the supported executable tags. Those tags are: bash, terraform, azurecli-interactive, and azurecli.

If a code block has a non supported tag like YAML or HTML it will simply render the code block as text and continue parsing the document. 

For example:

```YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: azure-vote-back
spec:
  replicas: 1
  selector:
    matchLabels:
      app: azure-vote-back
  template:
    metadata:
      labels:
        app: azure-vote-back
    spec:
      nodeSelector:
        "kubernetes.io/os": linux
      containers:
      - name: azure-vote-back
        image: mcr.microsoft.com/oss/bitnami/redis:6.0.8
        env:
        - name: ALLOW_EMPTY_PASSWORD
          value: "yes"
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 250m
            memory: 256Mi
        ports:
        - containerPort: 6379
          name: redis
---
apiVersion: v1
kind: Service
metadata:
  name: azure-vote-back
spec:
  ports:
  - port: 6379
  selector:
    app: azure-vote-back
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: azure-vote-front
spec:
  replicas: 1
  selector:
    matchLabels:
      app: azure-vote-front
  template:
    metadata:
      labels:
        app: azure-vote-front
    spec:
      nodeSelector:
        "kubernetes.io/os": linux
      containers:
      - name: azure-vote-front
        image: mcr.microsoft.com/azuredocs/azure-vote-front:v1
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 250m
            memory: 256Mi
        ports:
        - containerPort: 80
        env:
        - name: REDIS
          value: "azure-vote-back"
---
apiVersion: v1
kind: Service
metadata:
  name: azure-vote-front
spec:
  type: LoadBalancer
  ports:
  - port: 80
  selector:
    app: azure-vote-front

```


# Next Steps

These are the next steps... at some point we need to do something here