.PHONY: build-ie build-all run-ie clean test-all test all

BINARY_DIR := bin
IE_BINARY := $(BINARY_DIR)/ie
API_BINARY := $(BINARY_DIR)/api

# -------------------------- Native build targets ------------------------------

build-ie:
	@echo "Building the Innovation Engine CLI..."
	@CGO_ENABLED=0 go build -o "$(IE_BINARY)" cmd/ie/ie.go


build-all: build-ie

# ------------------------------ Install targets -------------------------------

install-ie:
	@echo "Installing the Innovation Engine CLI..."
	@CGO_ENABLED=0 go install cmd/ie/ie.go

# ------------------------------ Test targets ----------------------------------

WITH_COVERAGE := false

test-all:
	@go clean -testcache
ifeq ($(WITH_COVERAGE), true)
	@echo "Running all tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
else
	@echo "Running all tests..."
	@go test -v ./...
endif


SUBSCRIPTION ?= 00000000-0000-0000-0000-000000000000
SCENARIO ?= ./README.md
WORKING_DIRECTORY ?= $(PWD)
test-scenario:
	@echo "Running scenario $(SCENARIO)"
	$(IE_BINARY) test $(SCENARIO) --subscription $(SUBSCRIPTION) --working-directory $(WORKING_DIRECTORY)

test-scenarios:
	@echo "Testing out the scenarios"
	for dir in ./scenarios/ocd/*/; do \
		($(MAKE) test-scenario SCENARIO="$${dir}README.md" SUBCRIPTION="$(SUBSCRIPTION)") || exit $$?; \
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
		($(MAKE) test-scenario SCENARIO="$${dir}README.md" SUBCRIPTION="$(SUBSCRIPTION)" WORKING_DIRECTORY="$${dir}") || exit $$?; \
	done

# ------------------------------- Run targets ----------------------------------

run-ie: build-ie
	@echo "Running the Innovation Engine CLI"
	@"$(IE_BINARY)"

clean:
	@echo "Cleaning up"
	@rm -rf "$(BINARY_DIR)"

