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
pre-commit-run: ## Pre-commit run
	@echo "Run pre-commit hooks..."
	@pre-commit run --all-files

.PHONY: lint
lint: ## Run lint
	@golangci-lint run

