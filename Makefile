.ONESHELL:
SHELL := /bin/bash

VERSION ?= $(shell git describe --tags --always --dirty)


.PHONY: run
run:
	@source .env
	@go run ./cmd/core


.PHONY: migrate
migrate:
	@source .env
	@tern migrate --migrations migrations


.PHONY: build
build:
	@go build \
		-ldflags="-w -s -X main.Version=$(VERSION)" \
		-o core ./cmd/core


.PHONY: clear-cache
clear-cache:
	@docker exec -it valkey redis-cli FLUSHALL
