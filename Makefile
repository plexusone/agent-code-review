.PHONY: all build test lint clean install

BINARY_NAME := acr
GO := go
GOFLAGS := -v

all: lint test build

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) ./cmd/acr

test:
	$(GO) test $(GOFLAGS) ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY_NAME)
	$(GO) clean

install:
	$(GO) install $(GOFLAGS) ./cmd/acr

# Run the MCP server
serve:
	$(GO) run ./cmd/acr serve

# Format code
fmt:
	$(GO) fmt ./...
	goimports -w -local github.com/plexusone/agent-code-review .

# Update dependencies
deps:
	$(GO) mod tidy
	$(GO) mod verify
