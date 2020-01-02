
.DEFAULT_GOAL := build
APP_DIR := $(shell git rev-parse --show-toplevel)/app
VERSION := $(shell git log -1 --format=%h -- $(APP_DIR))

.PHONY: build
build:
	mkdir -p build
	go build -o build/app -ldflags "-X main.Version=$(VERSION)" .

.PHONY: docker-build
docker-build: build
	docker build 

