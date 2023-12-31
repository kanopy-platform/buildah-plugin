GO_MODULE := $(shell git config --get remote.origin.url | grep -o 'github\.com[:/][^.]*' | tr ':' '/')
CMD_NAME := $(shell basename ${GO_MODULE})
GIT_COMMIT := $(shell git rev-parse HEAD)
PLUGIN_TYPE ?= drone
ARCH ?= "amd64"
REGISTRY_NAME ?= registry.example.com
CONTAINER_RUNTIME ?= docker

RUN ?= .*
PKG ?= ./...

.PHONY: test
test: tidy ## Run go tests in local environment
	golangci-lint run --timeout=5m $(PKG)
	go test -cover -short -run=$(RUN) $(PKG)

.PHONY: dev
dev: tidy
	mkdir -p bin
	go build -o bin/${CMD_NAME} ./cmd/

.PHONY: tidy
tidy:
	go mod tidy
	go mod verify

.PHONY: local
local:
	${CONTAINER_RUNTIME} build --build-arg GIT_COMMIT=${GIT_COMMIT} --build-arg PLUGIN_TYPE=${PLUGIN_TYPE} -t ${REGISTRY_NAME}/$(CMD_NAME):latest .

.PHONY: local-push
local-push: local
	${CONTAINER_RUNTIME} push ${REGISTRY_NAME}/$(CMD_NAME):latest

.PHONY: local-run
local-run: local-push ## Build and run the application in a local container
	${CONTAINER_RUNTIME} run ${REGISTRY_NAME}/$(CMD_NAME):latest

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
