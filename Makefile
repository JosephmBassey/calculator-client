# SHELL := /bin/bash
export
BINARY=client
include develop.env

.PHONY: build
build:
	go build -mod=vendor -v -o $(BINARY) ./cmd


.PHONY: run
run: ## Compile and run locally
	go run ./cmd