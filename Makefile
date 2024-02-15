SHELL := /bin/bash

BIN_NAME ?= awi

.PHONY: all
all: build

# help credits: https://github.com/kubernetes-sigs/kubebuilder/blob/master/Makefile

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	golangci-lint run --fix

.PHONY: lint
lint: ## Run golangci-lint linter
	golangci-lint run

.PHONY: build
build: ## Build binary
	CGO_ENABLED=0 go build -ldflags="-w -s" -o ${BIN_NAME}

.PHONY: install
install: ## Install binary in $GOPATH/bin
ifndef GOPATH
	@echo "Error: GOPATH is not set"
	exit 1
endif
	CGO_ENABLED=0 go build -ldflags="-w -s" -o ${GOPATH}/bin/${BIN_NAME}

.PHONY: unit-test
unit-test: ## Run unit tests
	go test ./...

.PHONY: race-unit-test
race-unit-test: ## Run unit tests with race flag to detect races
	go test -race ./...

.PHONY: test-fmt
test-fmt: ## Test fmt and imports formatting
	test -z $$(goimports -w -l cmd pkg)

.PHONY: test
test: test-fmt lint race-unit-test ## Run fmt, linters and unit tests
