.PHONY: build-eg run-eg clean

BINARY_DIR := bin
EG_BINARY := $(BINARY_DIR)/eg

# -------------------------- Native build targets ------------------------------

RELEASE_BUILD := false
LATEST_TAG := $(shell git describe --tags --abbrev=0)
LATEST_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
MODULE_ROOT :=  $(shell go list -m)
build-eg:
	@echo "Building EG CLI..."
ifeq ($(RELEASE_BUILD), true)
	@CGO_ENABLED=0 go build -ldflags "-X $(MODULE_ROOT)/cmd/eg/commands.VERSION=$(LATEST_TAG) -X $(MODULE_ROOT)/cmd/eg/commands.COMMIT=$(LATEST_COMMIT) -X $(MODULE_ROOT)/cmd/eg/commands.DATE=$(BUILD_DATE)" -o "$(EG_BINARY)" main.go
else
	@CGO_ENABLED=0 go build -ldflags "-X $(MODULE_ROOT)/cmd/eg/commands.VERSION=dev -X $(MODULE_ROOT)/cmd/eg/commands.COMMIT=$(LATEST_COMMIT) -X $(MODULE_ROOT)/cmd/eg/commands.DATE=$(BUILD_DATE)" -o "$(EG_BINARY)" main.go
endif

# ------------------------------ Install targets -------------------------------

install-eg:
	@echo "Installing the EG CLI..."
	@CGO_ENABLED=0 go install -ldflags "-X $(MODULE_ROOT)/cmd/eg/commands.VERSION=dev -X $(MODULE_ROOT)/cmd/eg/commands.COMMIT=$(LATEST_COMMIT) -X $(MODULE_ROOT)/cmd/eg/commands.DATE=$(BUILD_DATE)" main.go

# ------------------------------- Run targets ----------------------------------

run-eg: build-eg
	@echo "Running the EG CLI"
	@"$(EG_BINARY)"

clean:
	@echo "Cleaning up"
	@rm -rf "$(EG_DIR)"

