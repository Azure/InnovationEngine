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

deploy-api-container: build-api-container
	@echo "Deploying the Innovation Engine API container"
	@docker run -d -p 8080:8080 innovation-engine-api
