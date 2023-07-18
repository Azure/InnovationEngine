---
title: 'Quickstart: Use the Azure CLI to create a Linux VM'
---

# Prerequisites

Innovation Engine can process prerequisites for documents. This code section tests that the pre requisites functionality works in Innovation Engine.
It will run the following real prerequisites along with a look for and fail to run a fake prerequisite.

You must have completed [Fuzzy Matching Test](testScripts/fuzzyMatchTest.md) and you must have completed [Comment Test](testScripts/CommentTest.md)

You also need to have completed [This is a fake file](testScripts/fakefile.md)

And there are going to be additional \ and ( to throw off the algorithm... 

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

# Test non-executable code blocks
If a code block does not have an executable tag it will simply render the codeblock as text

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

# Testing regular comments 

Innovation engine is able to handle comments and actual do fancy things with special comments.

There are comments you can't see here.
<!--This is a test comment in markdown -->


<!--This is a multi line comment in markdown


 in markdown -->

# Testing Declaring Environment Variables from Comments
Innovation Engine can declare environment variables via hidden inline comments. This feature is useful for running documents E2E as part of CI/CD

<!!--
```variables
export MY_VARIABLE=willBeChanged
```
 -->
<!--
Here is an example of that
```variables
export MY_VARIABLE=myVariable
```
-->

```azurecli-interactive
echo $MY_VARIABLE
```


# Test Running an Azure Command
```azurecli-interactive
az group exists --name MyResourceGroup
```

# Next Steps

These are the next steps... at some point we need to do something here