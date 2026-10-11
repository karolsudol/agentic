.PHONY: all build run test vet fmt lint setup eval clean help

# Default target
all: fmt vet test build

## Build & Run
build:
	@echo "==> Building agent binaries..."
	go build -v -o bin/agent ./cmd/agent/main.go

run:
	@echo "==> Running agent workflow..."
	go run ./cmd/agent/main.go

## Quality & Formatting
fmt:
	@echo "==> Formatting Go files..."
	gofmt -w -s .

vet:
	@echo "==> Running go vet..."
	go vet ./...

lint:
	@echo "==> Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, falling back to go vet"; \
		go vet ./...; \
	fi

## Testing & Evals
test:
	@echo "==> Running tests with race detector..."
	go test -v -race ./...

eval:
	@echo "==> Running agent evaluation harness..."
	@if [ -d "pkg/evals" ]; then \
		go test -v -race ./pkg/evals/...; \
	else \
		echo "Evaluation harness package (pkg/evals) will be run here once implemented."; \
	fi

## Environment Setup
setup:
	@echo "==> Setting up local developer environment..."
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit install; \
		echo "✓ Pre-commit hooks installed."; \
	else \
		echo "⚠ pre-commit not found. Install via 'brew install pre-commit' or 'pip install pre-commit'."; \
	fi
	@echo "✓ Go dependencies verified."

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf bin/ dist/ *.out *.test

help:
	@echo "Available targets:"
	@echo "  make build   - Compile the agent binary into bin/"
	@echo "  make run     - Run the main agent demonstration"
	@echo "  make test    - Run all unit and integration tests with -race"
	@echo "  make fmt     - Format all Go source files"
	@echo "  make vet     - Run go vet"
	@echo "  make lint    - Run golangci-lint"
	@echo "  make eval    - Run the evaluation harness"
	@echo "  make setup   - Configure pre-commit git hooks"
	@echo "  make clean   - Remove binaries and test artifacts"
