.PHONY: build-ie build-all run-ie clean test-all test all

BINARY_DIR := bin
IE_BINARY := $(BINARY_DIR)/ie
API_BINARY := $(BINARY_DIR)/api

# -------------------------- Native build targets ------------------------------

RELEASE_BUILD := false
LATEST_TAG := $(shell git describe --tags --abbrev=0)
LATEST_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
MODULE_ROOT :=  $(shell go list -m)
build-ie:
	@echo "Building the Innovation Engine CLI..."
ifeq ($(RELEASE_BUILD), true)
	@CGO_ENABLED=0 go build -ldflags "-X $(MODULE_ROOT)/cmd/ie/commands.VERSION=$(LATEST_TAG) -X $(MODULE_ROOT)/cmd/ie/commands.COMMIT=$(LATEST_COMMIT) -X $(MODULE_ROOT)/cmd/ie/commands.DATE=$(BUILD_DATE)" -o "$(IE_BINARY)" cmd/ie/ie.go
else
	@CGO_ENABLED=0 go build -ldflags "-X $(MODULE_ROOT)/cmd/ie/commands.VERSION=dev -X $(MODULE_ROOT)/cmd/ie/commands.COMMIT=$(LATEST_COMMIT) -X $(MODULE_ROOT)/cmd/ie/commands.DATE=$(BUILD_DATE)" -o "$(IE_BINARY)" cmd/ie/ie.go
endif


build-all: build-ie

# ------------------------------ Install targets -------------------------------

install-ie:
	@echo "Installing the Innovation Engine CLI..."
	@CGO_ENABLED=0 go install -ldflags "-X $(MODULE_ROOT)/cmd/ie/commands.VERSION=dev -X $(MODULE_ROOT)/cmd/ie/commands.COMMIT=$(LATEST_COMMIT) -X $(MODULE_ROOT)/cmd/ie/commands.DATE=$(BUILD_DATE)" cmd/ie/ie.go

# ------------------------------ Test targets ----------------------------------

WITH_COVERAGE := false
test-all:
	@go clean -testcache
ifeq ($(WITH_COVERAGE), true)
	@echo "Running all tests with coverage..."
	@go test -v -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
else
	@echo "Running all tests..."
	@go test -v ./...
endif


SUBSCRIPTION ?= 00000000-0000-0000-0000-000000000000
SCENARIO ?= ./README.md
WORKING_DIRECTORY ?= $(PWD)
ENVIRONMENT ?= local
test-scenario:
	@echo "Running scenario $(SCENARIO)"
ifeq ($(SUBSCRIPTION), 00000000-0000-0000-0000-000000000000)
	$(IE_BINARY) test $(SCENARIO) --working-directory $(WORKING_DIRECTORY) --environment $(ENVIRONMENT)
else
	$(IE_BINARY) test $(SCENARIO) --subscription $(SUBSCRIPTION) --working-directory $(WORKING_DIRECTORY) --environment $(ENVIRONMENT)
endif

test-scenarios:
	@echo "Testing out the scenarios"
	for dir in ./scenarios/ocd/*/; do \
		($(MAKE) test-scenario SCENARIO="$${dir}README.md" SUBCRIPTION="$(SUBSCRIPTION)") || exit $$?; \
	done

test-local-scenarios:
	@echo "Testing out the local scenarios"
	for file in ./scenarios/testing/*.md; do \
		($(MAKE) test-scenario SCENARIO="$${file}") || exit $$?; \
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
		if echo "$${dir}" | grep -q "CreateContainerAppDeploymentFromSource"; then \
			continue; \
		fi; \
		($(MAKE) test-scenario SCENARIO="$${dir}README.md" SUBCRIPTION="$(SUBSCRIPTION)" WORKING_DIRECTORY="$${dir}" ENVIRONMENT="$(ENVIRONMENT)") || exit $$?; \
	done

test-docs:
	@echo "Testing all documents in the docs folder"
	for file in ./docs/*.md; do \
		($(MAKE) test-scenario SCENARIO="$${file}") || exit $$?; \
	done

# ------------------------------- Run targets ----------------------------------

run-ie: build-ie
	@echo "Running the Innovation Engine CLI"
	@"$(IE_BINARY)"

clean:
	@echo "Cleaning up"
	@rm -rf "$(BINARY_DIR)"

