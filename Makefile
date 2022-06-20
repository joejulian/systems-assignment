MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MKFILE_PATH))

.PHONY bindir:
bindir:
	mkdir -p bin

CLIENT_DEPS = $(shell find cmd/client pkg/client -name '*.go')
client: go.mod go.sum $(CLIENT_DEPS) bindir
	go build -o bin/client cmd/client/main.go

.PHONY client.docker:
client.docker:
	docker run -ti --rm -u $(shell id -u):$(shell id -g) \
	    -v $(CURRENT_DIR):/app -w /app \
		-e GOCACHE=/tmp golang:1.18 \
		make client

SERVER_DEPS = $(shell find cmd/server pkg/api pkg/flatfile pkg/server -name '*.go')
server: go.mod go.sum $(CLIENT_DEPS) bindir
	go build -o bin/server cmd/server/main.go

.PHONY server.docker:
server.docker:
	docker run -ti --rm -u $(shell id -u):$(shell id -g) \
	    -v $(CURRENT_DIR):/app -w /app \
		-e GOCACHE=/tmp golang:1.18 \
		make server

.PHONY all: client server

.PHONY clean:
clean:
	rm -rf bin

.PHONY test:
test:
	go test -v ./...

.PHONY lint:
lint:
	docker run --rm -v $(CURRENT_DIR):/app -w /app golangci/golangci-lint:v1.46.2 golangci-lint run
