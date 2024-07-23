.DEFAULT_GOAL := default

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)


.PHONY: all
all: help
	@:

IMAGE ?= quay.io/kmamgain/cdevent:latest

export DOCKER_CLI_EXPERIMENTAL=enabled

.PHONY: build # Build the container image
build:
	podman build --platform linux/amd64,linux/arm64 -t quay.io/kmamgain/cdevent:latest  .

.PHONY: publish # Push the image to the remote registry
publish:
	podman push $(IMAGE)

.PHONY: build-go
build-go:
	go build .

.PHONY: test
test: build
test: TESTFLAGS += -race -v
test: test-unit

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)