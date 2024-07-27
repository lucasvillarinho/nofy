GO_FILES := $(shell find . -name '*.go')

GOFMT := gofumpt
GOIMPORTS := goimports


.PHONY: help
help:
	@echo "Use: make [target]"
	@echo ""
	@echo "Tasks:"
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-40s %s\n", $$1, $$2}'

.phony: pre-commit-install
pre-commit-install: ## Pre-commit install
	@echo "Install pre-commit hooks..."
	@pre-commit install --hook-type commit-msg
	@pre-commit install

.phony: pre-commit-run
pre-commit-run: format ## Pre-commit run
	@echo "Run pre-commit hooks..."
	@pre-commit run --all-files
	@echo "Pre-commit hooks passed successfully"

build:
	@echo "Building..."
	@go build -o bin/ ./...
	@echo "Build completed successfully"


.PHONY: lint
lint: build ## Run lint
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Linter passed successfully"

.PHONY: format
format:
	@echo "Formatting code..."
	@$(GOFMT) -w $(GO_FILES)
	@$(GOIMPORTS) -w $(GO_FILES)
	@echo "Code formatted successfully"