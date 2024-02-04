# Copyright (c) 2022 AlertAvert.com.  All rights reserved.
# Created by M. Massenzio, 2022-03-14


GOOS ?= $(shell uname -s | tr "[:upper:]" "[:lower:]")
GOARCH ?= amd64
GOMOD := $(shell go list -m)

# Versioning
# The `version` is a static value, set in settings.yaml, and ONLY used to tag the release.
tag := $(shell cat settings.yaml | yq -r .version)

# The `release` is generated automatically from the most recent tag (version) and the
# current commit hash, and is used to tag the Docker image.
release := $(shell git describe --tags --always --dirty="-dev")
prog := majordomo
bin := build/bin/$(prog)-$(release)_$(GOOS)-$(GOARCH)

image := alertavert/$(prog)
compose := docker/compose.yaml
dockerfile := docker/Dockerfile

# Source files & Test files definitions
#
# Edit only the packages list, when adding new functionality,
# the rest is deduced automatically.
#
pkgs := $(shell find pkg -mindepth 1 -type d)
all_go := $(shell for d in $(pkgs); do find $$d -name "*.go"; done)
test_srcs := $(shell for d in $(pkgs); do find $$d -name "*_test.go"; done)
srcs := $(filter-out $(test_srcs),$(all_go))

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories.

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: clean
img=$(shell docker images -q --filter=reference=$(image))
clean: ## Cleans up the binary, container image and other data
	@rm -rf build/*
	@[ ! -z $(img) ] && docker rmi $(img) || true
	@rm -rf certs

version: ## Displays the current version tag (release)
	@echo $(release)

fmt: ## Formats the Go source code using 'go fmt'
	@go fmt $(pkgs) ./cmd

##@ Development
$(bin): cmd/main.go $(srcs)
	@mkdir -p $(shell dirname $(bin))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-X main.Release=$(release)" \
		-o $(bin) cmd/main.go

.PHONY: build
build: $(bin) ## Builds the server binary and the React UI
	@cd ui/app && npm run build
	@rm -rf build/ui && mv ui/app/build build/ui

.PHONY: test
test: $(srcs) $(test_srcs)  ## Runs all tests
	@mkdir -p build/reports
	ginkgo -p -keepGoing -cover -coverprofile=coverage.out -outputdir=build/reports $(pkgs)

.PHONY: watch
watch: $(srcs) $(test_srcs)  ## Runs all tests every time a source or test file changes
	ginkgo watch -p $(pkgs)

build/reports/coverage.out: test ## Runs all tests and generates the coverage report

.PHONY: coverage
coverage: build/reports/coverage.out ## Shows the coverage report in the browser
	@go tool cover -html=build/reports/coverage.out

.PHONY: all
all: build test ## Builds the binary and runs all tests

.PHONY: run
run: $(bin) ## Runs the server binary
	$(bin) -debug -port 5005

##@ Container Management
# Convenience targets to run locally containers and
# setup the test environments.

tag: ## Tags the current release
	@echo "Tagging release $(tag)"
	@git tag -a $(tag) -m "Release $(tag)"
	@git push origin $(tag)

container: ## Builds the container image
	docker build -f $(dockerfile) -t $(image):$(release) .

.PHONY: run-container
run-container: container ## Runs the container locally
	docker run --rm -it -p 8080:8080 $(image):$(release)

##@ UI Development

.PHONY: ui-setup
ui-setup: ## Installs the UI dependencies
	@cd ui/app && npm install

.PHONY: ui-dev
ui-dev: ## Runs the UI in development mode
	@cd ui/app && npm start

##@ TLS Support
#
# This section is WIP and subject to change

config_dir := ssl-config
ca-csr := $(config_dir)/ca-csr.json
ca-config := $(config_dir)/ca-config.json
server-csr := $(config_dir)/localhost-csr.json

.PHONY: gencert
gencert: $(ca-csr) $(config) $(server-csr) ## Generates all certificates in the certs directory (requires cfssl, see https://github.com/cloudflare/cfssl#installation)
	cfssl gencert \
		-initca $(ca-csr) | cfssljson -bare ca


	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=$(ca-config) \
		-profile=server \
		$(server-csr)  | cfssljson -bare server
	@mkdir -p certs
	@mv *.pem certs/
	@rm *.csr
	@chmod a+r certs/*
	@echo "Certificates generated in $(shell pwd)/certs"

.PHONY: clean-cert
clean-cert:
	@rm -rf certs
