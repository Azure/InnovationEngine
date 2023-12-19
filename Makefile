.PHONY: build-ie build-api build-all run-ie run-api clean test-all test all

BINARY_DIR := bin
IE_BINARY := $(BINARY_DIR)/ie
API_BINARY := $(BINARY_DIR)/api

# -------------------------- Native build targets ------------------------------

build-ie:
	@echo "Building the Innovation Engine CLI..."
	@CGO_ENABLED=0 go build -o "$(IE_BINARY)" cmd/ie/ie.go

build-api:
	@echo "Building the Innovation Engine API..."
	@CGO_ENABLED=0 go build -o "$(API_BINARY)" cmd/api/main.go

build-runner: build-ie build-api
	@echo "Building the Innovation Engine Runner..."
	@CGO_ENABLED=0 go build -o "$(BINARY_DIR)/runner" cmd/runner/main.go

build-all: build-ie build-api build-runner

# ------------------------------ Install targets -------------------------------

install-ie:
	@echo "Installing the Innovation Engine CLI..."
	@CGO_ENABLED=0 go install cmd/ie/ie.go

# ------------------------------ Test targets ----------------------------------

test-all:
	@echo "Running all tests..."
	@go clean -testcache
	@go test -v ./...

SUBSCRIPTION ?= 00000000-0000-0000-0000-000000000000
SCENARIO ?= ./README.md
test-scenario:
	@echo "Running scenario $(SCENARIO)"
	# $(IE_BINARY) test $(SCENARIO) --subscription $(SUBSCRIPTION)

test-scenarios:
	@echo "Testing out the scenarios"
	for dir in ./scenarios/ocd/*/; do \
		 $(MAKE) test-scenario SCENARIO="$${dir}README.md" SUBCRIPTION="$(SUBSCRIPTION)"; \
	done

test-upstream-scenarios:
	@echo "Pulling the upstream scenarios"
	@git config --global --add safe.directory /home/runner/work/InnovationEngine/InnovationEngine
	@git submodule update --init --recursive
	@echo "Testing out the upstream scenarios"
	for dir in ./upstream-scenarios/scenarios/*/; do \
		if ! [ -f $${dir}README.md ]; then \
			continue; \
		fi; \
		 $(MAKE) test-scenario SCENARIO="$${dir}README.md" SUBCRIPTION="$(SUBSCRIPTION)"; \
	done

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

API_IMAGE_TAG ?= latest

# Builds the API container. 
build-api-container:
	@echo "Building the Innovation Engine API container"
	@docker build -t innovation-engine-api:$(API_IMAGE_TAG) -f infra/api/Dockerfile .


# ----------------------------- Kubernetes targets -----------------------------

# Applies the ingress controller to the cluster and waits for it to be ready.
k8s-deploy-ingress-controller:
	@echo "Deploying the ingress controller to your local cluster..."
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.1/deploy/static/provider/cloud/deploy.yaml
	@kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=120s

# Deploys the API deployment, service, and ingress specifications to the 
# cluster, allowing the API to be accessed via the ingress controller.
k8s-deploy-api: build-api-container
	@echo "Deploying the Innovation Engine API container to your local cluster..."
	@kubectl apply -f infra/api/deployment.yaml
	@kubectl apply -f infra/api/service.yaml
	@kubectl apply -f infra/api/ingress.yaml

k8s-initialize-cluster: k8s-deploy-ingress-controller k8s-deploy-api
	@echo "Set up Kubernetes cluster for local development."

k8s-delete-ingress-controller:
	@echo "Deleting the ingress controller from your local cluster..."
	@kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.1/deploy/static/provider/cloud/deploy.yaml

k8s-delete-api:
	@echo "Deleting the Innovation Engine API container from your local cluster..."
	@kubectl delete -f infra/api/deployment.yaml
	@kubectl delete -f infra/api/service.yaml
	@kubectl delete -f infra/api/ingress.yaml

k8s-refresh-api: k8s-delete-api k8s-deploy-api
	@echo "Refreshed the Innovation Engine API container in your local cluster..."

k8s-delete-cluster: k8s-delete-api k8s-delete-ingress-controller
	@echo "Deleted Kubernetes cluster for local development."

k8s-refresh-cluster: k8s-delete-cluster k8s-initialize-cluster
	@echo "Refreshed Kubernetes cluster for local development."
