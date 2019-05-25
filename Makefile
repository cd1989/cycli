# Set shell to bash
SHELL := /bin/bash

# Current version of the project.x`
VERSION ?= v0.0.1

# This repo's root import path (under GOPATH).
ROOT := github.com/cd1989/cycli

# A list of all packages.
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /test)

# Git commit sha.
COMMIT := $(shell git rev-parse --short HEAD)

# Golang standard bin directory.
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

UNAME := $(shell uname)

# All targets.
.PHONY: lint test build

lint: $(GOMETALINTER)
	golangci-lint run --disable=gosimple --deadline=300s ./pkg/... ./cmd/...

build: build-local

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

test:
	go test $(PKGS)

build-local:
	CGO_ENABLED=0 GOARCH=amd64 go build -i -v -o ./bin/cycli -ldflags "-s -w" ./

.PHONY: clean
clean:
	-rm -vrf ${OUTPUT_DIR}
clean-generated:
	-rm -rf ./pkg/k8s/informers
	-rm -rf ./pkg/k8s/clientset
	-rm -rf ./pkg/k8s/listers