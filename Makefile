GO_MODULE := $(shell git config --get remote.origin.url | grep -o 'github\.com[:/][^.]*' | tr ':' '/')
CMD_NAME := $(shell basename ${GO_MODULE})
GIT_COMMIT := $(shell git rev-parse HEAD)
REGISTRY_NAME ?= registry.example.com
CONTAINER_RUNTIME ?= podman
TLS_VERIFY ?= false

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
	${CONTAINER_RUNTIME} build --build-arg GIT_COMMIT=${GIT_COMMIT} -t $(CMD_NAME):latest .

.PHONY: local-run
local-run: local ## Build and run the application in a local docker container
	${CONTAINER_RUNTIME} push ${REGISTRY_NAME}/$(CMD_NAME):latest --tls-verify=${TLS_VERIFY}
	${CONTAINER_RUNTIME} run $(CMD_NAME):latest

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
