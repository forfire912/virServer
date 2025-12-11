.PHONY: all build test clean run docker-build docker-up docker-down swagger

# Binary name
BINARY_NAME=virserver
VERSION?=0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go

test:
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -rf artifacts/
	rm -rf snapshots/
	rm -f coverage.txt

run: build
	./bin/$(BINARY_NAME)

deps:
	$(GOMOD) download
	$(GOMOD) tidy

swagger:
	swag init -g cmd/server/main.go -o docs

docker-build:
	docker build -t virserver:$(VERSION) .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

lint:
	golangci-lint run ./...

help:
	@echo "Available targets:"
	@echo "  all         - Run tests and build"
	@echo "  build       - Build the binary"
	@echo "  test        - Run tests with coverage"
	@echo "  clean       - Clean build artifacts"
	@echo "  run         - Build and run the server"
	@echo "  deps        - Download dependencies"
	@echo "  swagger     - Generate Swagger documentation"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up   - Start services with docker-compose"
	@echo "  docker-down - Stop services"
	@echo "  lint        - Run linter"
