
.DEFAULT_GOAL := build
APP_DIR := $(shell git rev-parse --show-toplevel)/app
VERSION := $(shell git log -1 --format=%h -- $(APP_DIR))
ifneq ($(REGISTRY),)
REGISTRY := $(REGISTRY)/
endif

.PHONY: build
build:
	mkdir -p build
	go build -o build/emoji-app-$$(go env GOOS) -ldflags "-X main.Version=$(VERSION)" .

.PHONY: build-linux
build-linux:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/emoji-app-$$(go env GOOS) -ldflags "-X main.Version=$(VERSION)" .

.PHONY: docker-build
docker-build: build
	docker build -t $(REGISTRY)emojis-app .

