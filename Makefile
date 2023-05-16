.PHONY: build-ie build-api build-all run-ie run-api clean

BINARY_DIR := bin
IE_BINARY := $(BINARY_DIR)/ie
API_BINARY := $(BINARY_DIR)/api

# -------------------------- Native build targets ------------------------------

build-ie:
	@echo "Building the Innovation Engine CLI..."
	@go build -o "$(IE_BINARY)" cmd/ie/ie.go

build-api:
	@echo "Building the Innovation Engine API..."
	@go build -o "$(API_BINARY)" cmd/api/main.go

build-runner: build-ie build-api
	@echo "Building the Innovation Engine Runner..."
	@go build -o "$(BINARY_DIR)/runner" cmd/runner/main.go

build-all: build-ie build-api build-runner

# ------------------------------- Run targets ----------------------------------

run-ie: build-ie
	@echo "Running the Innovation Engine CLI"
	@"$(IE_BINARY)"

run-api: build-api
	@echo "Running the Innovation Engine API"
	@"$(API_BINARY)"

clean:
	@echo "Cleaning up"
	@rm -rf "$(BINARY_DIR)"

# ----------------------------- Docker targets ---------------------------------

# Builds the API container. 
build-api-container:
	@echo "Building the Innovation Engine API container"
	@docker build -t innovation-engine-api -f infra/api/Dockerfile .


# ----------------------------- Kubernetes targets -----------------------------

# Applies the ingress controller to the cluster and waits for it to be ready.
k8s-apply-ingress-controller:
	@echo "Deploying the ingress controller to your local cluster..."
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.1/deploy/static/provider/cloud/deploy.yaml
	@kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=120s

# Deploys the API deployment, service, and ingress specifications to the 
# cluster, allowing the API to be accessed via the ingress controller.
k8s-apply-api: build-api-container
	@echo "Deploying the Innovation Engine API container to your local cluster..."
	@kubectl apply -f infra/api/deployment.yaml
	@kubectl apply -f infra/api/service.yaml
	@kubectl apply -f infra/api/ingress.yaml

k8s-initialize-cluster: k8s-apply-ingress-controller k8s-apply-api
	@echo "Set up Kubernetes cluster for local development."