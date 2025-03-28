.ONESHELL:
SHELL := /bin/bash

VERSION ?= $(shell git describe --tags --always --dirty)


.PHONY: run
run:
	source .env
	go run ./cmd/core


.PHONY: build
build:
	@go build \
		-ldflags="-w -s -X main.Version=$(VERSION)" \
		-o core ./cmd/core
