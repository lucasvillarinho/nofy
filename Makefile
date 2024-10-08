GO_FILES := $(shell find . -name '*.go')

GOFMT := gofumpt
GOIMPORTS := goimports


.PHONY: help
help:
	@echo "Use: make [target]"
	@echo ""
	@echo "Tasks:"
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-40s %s\n", $$1, $$2}'


.PHONY: pre-commit-run
pre-commit-run: format ## Pre-commit run
	@echo "Run pre-commit hooks..."
	@pre-commit run --all-files
	@echo "Pre-commit hooks passed successfully"

.PHONY: build
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
format: ## Format code
	@echo "Formatting code..."
	@gofumpt -w .
	@goimports -w .
	@golines -m 100 -w .
	@fieldalignment -fix ./...
	@echo "Code formatted successfully"

.PHONY: test
test: ## Run unit test
	go test -v -coverprofile=rawcover.out -json $$(go list ./... | grep -v "github.com/lucasvillarinho/nofy/examples" | grep -v "github.com/lucasvillarinho/nofy/tests/e2e") 2>&1 | tee /tmp/gotest.log | gotestfmt -hide successful-tests,empty-packages

.PHONY: e2e-test
e2e-test: ## Run e2e test
	go test -race $$(go list ./... | grep "github.com/lucasvillarinho/nofy/tests/e2e")