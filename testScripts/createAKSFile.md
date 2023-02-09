
# Welcome to this test document

## **FIXED In latest PR with the addition of the remove leading whitespace function** 

This file currently shows a bug with innovation engine. When the code block is tabbed, it 
does not parse the bash command correctly and breaks the EOF to its own line/command. This
is because in the parser, I simply add 3 ticks to the end of any code block instead of reading
it with the spaces and tabs... So in this case, the EOF is not registering with the 3 ticks
This should create a file name azure-vote.yaml

The ultimate fix to innovation engine would be to remove all leading spaces from new lines so 
that this is not an issue. This issue persists as even if it completes the creation of the YAML
file the YAML file is tabbed and thus

Gives this error "error: error parsing azure-vote-base.yaml: error converting YAML to JSON: yaml: 
line 34: could not find expected ':'"

```azurecli-interactive
export AKS_FILE_NAME=azure-vote-test2.yaml
```

        ```azurecli-interactive
        cat <<EOF > $AKS_FILE_NAME
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

        EOF
        ```

```azurecli-interactive
echo "hello world"
```

```azurecli-interactive
cat $AKS_FILE_NAME
```

<!--expected_similarity=.95-->
```Output
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

# Clean up testing files
```azurecli-interactive
rm $AKS_FILE_NAME
```
