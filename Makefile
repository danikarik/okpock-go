VERSION := $(shell git describe --tags --long 2>/dev/null || git rev-parse --short HEAD)

DB_NAME ?= okpock
TEST_DB_NAME ?= test_okpock
MYSQL := MYSQL_PWD=$(DB_PASS) mysql -u $(DB_USER)

all: help

help: ## Show usage
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

create: ## Create new database
	@echo 'CREATE DATABASE IF NOT EXISTS' $(DB_NAME) | $(MYSQL)

create-test: ## Create test database
	@echo 'DROP DATABASE IF EXISTS' $(TEST_DB_NAME) | $(MYSQL)
	@echo 'CREATE DATABASE IF NOT EXISTS' $(TEST_DB_NAME) | $(MYSQL)

drop: ## Drop database
	@echo 'DROP DATABASE IF EXISTS' $(DB_NAME) | $(MYSQL)

drop-test: ## Drop test database
	@echo 'DROP DATABASE IF EXISTS' $(TEST_DB_NAME) | $(MYSQL)

up: ## Create tables
	@$(MYSQL) $(DB_NAME) < ./migrations/up.sql

down: ## Drop tables
	@$(MYSQL) $(DB_NAME) < ./migrations/down.sql

fresh: down up ## Rebuild database
	@echo 'Done.'

download: ## Download dependencies
	@go mod download

test: ## Run go tests
	@go test -v -count=1 ./pkg/...

release: ## Build release binary
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -X main.Version=$(VERSION)" -o ./bin/application ./cmd/applicationd/...

build: ## Build local binary
	@go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/application ./cmd/applicationd/...

secret: ## Generate secret string and copy to buffer
	@openssl rand -hex 32 | tr -d '\n' | pbcopy

loc: ## Calculate LOC
	@find . -name '*.go' | xargs wc -l
