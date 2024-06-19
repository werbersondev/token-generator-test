# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
  include .env
  export $(shell sed 's/=.*//' .env)
endif

.PHONY: download
download:
	@echo "==> Downloading go.mod dependencies"
	@go mod download

.PHONY: install-tools
install-tools: download
	@echo "==> Installing moq"
	@go install github.com/matryer/moq@latest

.PHONY: setup
setup: install-tools
	@go mod tidy

.PHONY: test
test:
	@echo "==> Running tests"
	@go test -race -failfast ./...

.PHONY: test-coverage
test-coverage:
	@echo "==> Running tests with coverage"
	@go test -race -failfast -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: generate
generate:
	@echo "==> Running go generate"
	@go generate ./...

.PHONY: setup/local-dep
setup/local-dep:
	@docker compose down
	@docker compose up -d
	@sh scripts/seed_sonar_project.sh

.PHONY: run/http
run/http:
	@go run cmd/httpservice/main.go

.PHONY: run/worker
run/worker:
	@if [ -z "$(SONAR_AUTH_TOKEN)" ]; then \
		echo "Error: SONAR_AUTH_TOKEN is not set"; \
		exit 1; \
	fi
	@SONAR_AUTH_TOKEN=$(SONAR_AUTH_TOKEN) go run cmd/worker/main.go