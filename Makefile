# Image URL to use all building/pushing image targets
DEV_IMG ?= registry.devops.rivtower.com/cita-cloud/cita-node-proxy:v0.0.1
IMG ?= citacloud/cita-node-proxy

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.23

VERSION=$(shell git describe --tags --match 'v*' --always --dirty)
GIT_COMMIT?=$(shell git rev-parse --short HEAD)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o bin/cita-node-proxy main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: dev-build
dev-build: ## Build dev image with the manager.
	docker build --platform linux/amd64 -t ${DEV_IMG} . --build-arg version=$(GIT_COMMIT)

.PHONY: dev-push
dev-push: ## Push dev image with the manager.
	docker push ${DEV_IMG}

.PHONY: image-latest
image-latest:
	# Build image with latest stable
	docker buildx build -t $(IMG):latest --build-arg version=$(GIT_COMMIT) \
    		--platform linux/amd64,linux/arm64 . --push

.PHONY: image-version
image-version:
	[ -z `git status --porcelain` ] || (git --no-pager diff && exit 255)
	docker buildx build -t $(IMG):$(VERSION) --build-arg version=$(GIT_COMMIT) \
		--platform linux/amd64,linux/arm64 . --push
