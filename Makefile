.PHONY: build-ie build-api build-all run-ie run-api clean

BINARY_DIR := bin
IE_BINARY := $(BINARY_DIR)/ie
API_BINARY := $(BINARY_DIR)/api

build-ie:
	@echo "Building the Innovation Engine CLI..."
	@go build -o "$(IE_BINARY)" cmd/ie/ie.go

build-api:
	@echo "Building the Innovation Engine API..."
	@go build -o "$(API_BINARY)" cmd/api/main.go

build-all: build-ie build-api

run-ie: build-ie
	@echo "Running the Innovation Engine CLI"
	@"$(IE_BINARY)"

run-api: build-api
	@echo "Running the Innovation Engine API"
	@"$(API_BINARY)"

clean:
	@echo "Cleaning up"
	@rm -rf "$(BINARY_DIR)"

build-api-container:
	@echo "Building the Innovation Engine API container"
	@docker build -t innovation-engine-api -f infra/api/Dockerfile .
