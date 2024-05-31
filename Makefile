# Copyright (c) AlertAvert.com.  All rights reserved.
# Created by M. Massenzio, 2022-03-14

GOOS ?= $(shell uname -s | tr "[:upper:]" "[:lower:]")
GOARCH ?= $(shell uname -m)
ifeq ($(GOARCH),x86_64)
	GOARCH=amd64
endif
GOMOD := $(shell go list -m)

# Versioning
# The `version` is a static value, set in settings.yaml, and ONLY used to tag the release.
VERSION ?= $(shell cat settings.yaml | yq -r .version)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
RELEASE := v$(VERSION)-g$(GIT_COMMIT)

prog := majordomo
bin := $(prog)-$(RELEASE)_$(GOOS)-$(GOARCH)

image := alertavert/$(prog)
compose := docker/compose.yaml
dockerfile := Dockerfile

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
	@find . -name "*.out" -exec rm {} \;
	@[ -n "$(img)" ] && docker rmi $(img) || true
	@rm -rf certs

version: ## Displays the current version tag (release)
	@echo v$(VERSION)

fmt: ## Formats the Go source code using 'go fmt'
	@go fmt $(pkgs) ./cmd

##@ Development
.PHONY: build
build: cmd/main.go $(srcs)
	@mkdir -p build/bin
	@echo "Building rel. $(RELEASE); OS/Arch: $(GOOS)/$(GOARCH) - Pkg: $(GOMOD)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-X main.Release=$(RELEASE)" \
		-o build/bin/$(bin) cmd/main.go
	@echo "Majordomo Server $(shell basename $(bin)) built"


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

PORT ?= 5005
.PHONY: dev
dev: build ## Runs the server binary in development mode
	build/bin/$(bin) -debug -port $(PORT)

##@ Container Management
# Convenience targets to run locally containers and
# setup the test environments.

tag: ## Tags the current release
	@echo "Tagging version v$(VERSION) at commit $(GIT_COMMIT)"
	@git tag -a v$(VERSION) -m "Release $(RELEASE)"
	@git push origin --tags

container: build/bin/$(bin) ## Builds the container image
	docker build -f $(dockerfile) \
		--build-arg="VERSION=$(VERSION)" \
		-t $(image):$(RELEASE) .

.PHONY: run-container
run-container:  ## Runs the container locally
	docker run --rm -it -p 5000:5000 \
		-v $${HOME}/.majordomo:/etc/majordomo \
		$(image):$(RELEASE)

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
