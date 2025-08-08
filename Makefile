COVERAGE_OUTPUT=coverage.output
COVERAGE_HTML=coverage.html

# Load environment variables from .env file
-include .env
export

## @ Help
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make [target]\033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_/-]+:.*?##/ { printf "\033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

## @ Tools
.PHONY: tools/install
tools/install:  ## Install gofumpt, gocritic, gotestfmt, swaggo, goose and mockery
	@go install mvdan.cc/gofumpt@latest
	@go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest

## @ Linter
.PHONY: lint format
lint:  ## Run golangci-lint
	@gocritic check ./...

format:  ## Format code
	@gofumpt -e -l -w .

## @ Tests
.PHONY: test test/unit test/coverage test/coverage-browser mocks/generate mocks/clean

mocks/generate: mocks/clean  ## Generate mock files
	@find . -type f -name '*.go' ! -path '*/mocks/*' -exec dirname {} \; | sort -u | \
	xargs -I{} sh -c '\
	  mockery --all \
	         --dir="{}" \
	         --output="{}/mocks" \
	         --outpkg=mocks \
	         --quiet \
	'

mocks/clean:  ## Clean mock files
	@rm -rf */mocks/*

test:  ## Run all tests
	@$(GO) test -covermode=atomic -count=1 -v ./... -race -coverprofile $(COVERAGE_OUTPUT) -json -ldflags=-extldflags=-Wl | gotestfmt

test/unit:  ## Run only unit tests
	@$(GO) test -covermode=atomic -count=1 -v ./... -short -race -coverprofile $(COVERAGE_OUTPUT) -json -ldflags=-extldflags=-Wl | gotestfmt

test/coverage: ## Run tests, make coverage report and display it into browser
	@go test $(shell go list ./... | grep -v '/mock\|/fixture\|/docs\|/ioc') -covermode=atomic -count=1 -race -coverprofile $(COVERAGE_OUTPUT) -ldflags=-extldflags=-Wl
	@go tool cover -html=$(COVERAGE_OUTPUT) -o $(COVERAGE_HTML)

test/coverage-browser: test/coverage ## Open coverage report in browser
	@wslview $(COVERAGE_HTML)

## @ Application
.PHONY: swagger run
SWAGGER_DEPS := $(shell find internal/http/routes -type f -name '*.go') \
                $(shell find internal/app/*/dto -type f -name '*.go') \
				internal/http/server.go

docs/docs.go: $(SWAGGER_DEPS) ## Generate swagger docs if dependencies changed
	@echo "📝 Generating Swagger docs..."
	@swag init --generalInfo internal/http/server.go --output ./docs

swagger: docs/docs.go  ## Generate swagger docs

run: swagger  ## Run backend http server
	@go run cmd/server/main.go

## @ Clean
.PHONY: clean clean/coverage-cache
clean/coverage-cache:
	@rm -rf $(COVERAGE_OUTPUT)
	@rm -rf $(COVERAGE_HTML)

clean: clean_coverage_cache ## Remove cache files

## @ Migrations
.PHONY: migrate/up migrate/down
migrate/up:  ## Run database migrations up
	@goose -dir .migrations postgres "user=$(DATABASE_USER) password=$(DATABASE_PASSWORD) host=$(DATABASE_HOST) port=$(DATABASE_PORT) dbname=$(DATABASE_NAME) sslmode=disable" up

migrate/down:  ## Run database migrations down
	@goose -dir .migrations postgres "user=$(DATABASE_USER) password=$(DATABASE_PASSWORD) host=$(DATABASE_HOST) port=$(DATABASE_PORT) dbname=$(DATABASE_NAME) sslmode=disable" down


## @ Dependencies
.PHONY: dependencies/up
dependencies/up:
	@docker compose up -d zipkin
	@docker compose up -d postgres
