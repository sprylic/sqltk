# Makefile for SQL Builder Integration Tests
# This Makefile provides easy configuration and execution of integration tests

# Load configuration from config.env if it exists
ifneq (,$(wildcard config.env))
    include config.env
    export
endif

# Default configuration
DB_HOST ?= localhost
DB_PORT ?= 3306
DB_USER ?= root
DB_PASSWORD ?= password
DB_NAME ?= test_sqltk
DB_TYPE ?= mysql

# PostgreSQL configuration
PG_HOST ?= localhost
PG_PORT ?= 5432
PG_USER ?= postgres
PG_PASSWORD ?= postgres
PG_DB ?= postgres

# Docker configuration
DOCKER_MYSQL_IMAGE ?= mysql:8.0
DOCKER_POSTGRES_IMAGE ?= postgres:15
DOCKER_MYSQL_PORT ?= 3306
DOCKER_POSTGRES_PORT ?= 5432

# Test configuration
TEST_TIMEOUT ?= 30s
TEST_RACE ?= false
TEST_COVERAGE ?= false
TEST_VERBOSE ?= true

# Go configuration
GO ?= go
GOFLAGS ?= -mod=mod

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)SQL Builder Integration Test Makefile$(RESET)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Environment variables:$(RESET)"
	@echo "  DB_HOST          Database host (default: localhost)"
	@echo "  DB_PORT          Database port (default: 3306)"
	@echo "  DB_USER          Database user (default: root)"
	@echo "  DB_PASSWORD      Database password (default: password)"
	@echo "  DB_NAME          Database name (default: test_sqltk)"
	@echo "  DB_TYPE          Database type: mysql or postgres (default: mysql)"
	@echo "  PG_HOST          PostgreSQL host (default: localhost)"
	@echo "  PG_PORT          PostgreSQL port (default: 5432)"
	@echo "  PG_USER          PostgreSQL user (default: postgres)"
	@echo "  PG_PASSWORD      PostgreSQL password (default: postgres)"
	@echo "  PG_DB            PostgreSQL default database (default: postgres)"
	@echo "  TEST_TIMEOUT     Test timeout (default: 30s)"
	@echo "  TEST_RACE        Enable race detection (default: false)"
	@echo "  TEST_COVERAGE    Enable coverage reporting (default: false)"
	@echo "  TEST_VERBOSE     Verbose test output (default: true)"
	@echo ""
	@echo "$(YELLOW)Configuration:$(RESET)"
	@if [ -f config.env ]; then \
		echo "  config.env found and loaded"; \
		echo "  Current settings:"; \
		echo "    DB_HOST=$(DB_HOST)"; \
		echo "    DB_PORT=$(DB_PORT)"; \
		echo "    DB_USER=$(DB_USER)"; \
		echo "    PG_HOST=$(PG_HOST)"; \
		echo "    PG_PORT=$(PG_PORT)"; \
		echo "    PG_USER=$(PG_USER)"; \
	else \
		echo "  config.env not found - using defaults"; \
		echo "  Copy config.env.example to config.env to customize"; \
	fi

.PHONY: deps
deps: ## Install dependencies
	@echo "$(BLUE)Installing dependencies...$(RESET)"
	$(GO) mod download
	$(GO) mod tidy

.PHONY: test
test: ## Run unit tests
	@echo "$(BLUE)Running unit tests...$(RESET)"
	$(GO) test $(GOFLAGS) ./... $(if $(filter true,$(TEST_VERBOSE)),-v) $(if $(filter true,$(TEST_RACE)),-race) $(if $(filter true,$(TEST_COVERAGE)),-cover)

.PHONY: test-integration
test-integration: ## Run integration tests with database
	@echo "$(BLUE)Running integration tests...$(RESET)"
	@if [ "$(DB_TYPE)" = "mysql" ]; then \
		echo "$(YELLOW)Using MySQL database$(RESET)"; \
		MYSQL_DSN="$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/" $(GO) test $(GOFLAGS) -timeout $(TEST_TIMEOUT) $(if $(filter true,$(TEST_VERBOSE)),-v) $(if $(filter true,$(TEST_RACE)),-race) $(if $(filter true,$(TEST_COVERAGE)),-cover) -tags=integration ./... -run "TestMySQLIntegration"; \
	elif [ "$(DB_TYPE)" = "postgres" ]; then \
		echo "$(YELLOW)Using PostgreSQL database$(RESET)"; \
		POSTGRES_DSN="postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=disable" $(GO) test $(GOFLAGS) -timeout $(TEST_TIMEOUT) $(if $(filter true,$(TEST_VERBOSE)),-v) $(if $(filter true,$(TEST_RACE)),-race) $(if $(filter true,$(TEST_COVERAGE)),-cover) -tags=integration ./... -run "TestPostgresIntegration"; \
	else \
		echo "$(RED)Unsupported database type: $(DB_TYPE)$(RESET)"; \
		exit 1; \
	fi

.PHONY: test-mysql
test-mysql: ## Run MySQL integration tests
	@echo "$(BLUE)Running MySQL integration tests...$(RESET)"
	MYSQL_DSN="$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/" $(GO) test $(GOFLAGS) -timeout $(TEST_TIMEOUT) $(if $(filter true,$(TEST_VERBOSE)),-v) $(if $(filter true,$(TEST_RACE)),-race) $(if $(filter true,$(TEST_COVERAGE)),-cover) -tags=integration ./... -run "TestMySQLIntegration"

.PHONY: test-postgres
test-postgres: ## Run PostgreSQL integration tests
	@echo "$(BLUE)Running PostgreSQL integration tests...$(RESET)"
	POSTGRES_DSN="postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=disable" $(GO) test $(GOFLAGS) -timeout $(TEST_TIMEOUT) $(if $(filter true,$(TEST_VERBOSE)),-v) $(if $(filter true,$(TEST_RACE)),-race) $(if $(filter true,$(TEST_COVERAGE)),-cover) -tags=integration ./... -run "TestPostgresIntegration"

.PHONY: docker-mysql
docker-mysql: ## Start MySQL database in Docker
	@echo "$(BLUE)Starting MySQL database in Docker...$(RESET)"
	@if docker ps -q -f name=mysql-sqltk-test | grep -q .; then \
		echo "$(YELLOW)MySQL container already running$(RESET)"; \
	else \
		docker run -d --name mysql-sqltk-test \
			-e MYSQL_ROOT_PASSWORD=$(DB_PASSWORD) \
			-e MYSQL_DATABASE=$(DB_NAME) \
			-p $(DOCKER_MYSQL_PORT):3306 \
			$(DOCKER_MYSQL_IMAGE); \
		echo "$(GREEN)MySQL container started$(RESET)"; \
		echo "$(YELLOW)Waiting for MySQL to be ready...$(RESET)"; \
		sleep 10; \
	fi

.PHONY: docker-postgres
docker-postgres: ## Start PostgreSQL database in Docker
	@echo "$(BLUE)Starting PostgreSQL database in Docker...$(RESET)"
	@if docker ps -q -f name=postgres-sqltk-test | grep -q .; then \
		echo "$(YELLOW)PostgreSQL container already running$(RESET)"; \
	else \
		docker run -d --name postgres-sqltk-test \
			-e POSTGRES_PASSWORD=$(PG_PASSWORD) \
			-e POSTGRES_DB=$(PG_DB) \
			-e POSTGRES_USER=$(PG_USER) \
			-p $(DOCKER_POSTGRES_PORT):5432 \
			$(DOCKER_POSTGRES_IMAGE); \
		echo "$(GREEN)PostgreSQL container started$(RESET)"; \
		echo "$(YELLOW)Waiting for PostgreSQL to be ready...$(RESET)"; \
		sleep 5; \
	fi

.PHONY: docker-stop
docker-stop: ## Stop all test database containers
	@echo "$(BLUE)Stopping test database containers...$(RESET)"
	@docker stop mysql-sqltk-test postgres-sqltk-test 2>/dev/null || true
	@docker rm mysql-sqltk-test postgres-sqltk-test 2>/dev/null || true
	@echo "$(GREEN)Test containers stopped and removed$(RESET)"

.PHONY: test-docker-mysql
test-docker-mysql: docker-mysql ## Run MySQL tests with Docker database
	@echo "$(BLUE)Running MySQL integration tests with Docker...$(RESET)"
	@DB_HOST=localhost DB_PORT=$(DOCKER_MYSQL_PORT) DB_USER=root DB_PASSWORD=$(DB_PASSWORD) $(MAKE) test-mysql

.PHONY: test-docker-postgres
test-docker-postgres: docker-postgres ## Run PostgreSQL tests with Docker database
	@echo "$(BLUE)Running PostgreSQL integration tests with Docker...$(RESET)"
	@PG_HOST=localhost PG_PORT=$(DOCKER_POSTGRES_PORT) PG_USER=$(PG_USER) PG_PASSWORD=$(PG_PASSWORD) PG_DB=$(PG_DB) $(MAKE) test-postgres

.PHONY: test-all-docker
test-all-docker: ## Run all integration tests with Docker databases
	@echo "$(BLUE)Running all integration tests with Docker...$(RESET)"
	@$(MAKE) test-docker-mysql
	@$(MAKE) test-docker-postgres
	@$(MAKE) docker-stop

.PHONY: test-all-integration
test-all-integration: ## Run all integration tests (MySQL + PostgreSQL)
	@echo "$(BLUE)Running all integration tests...$(RESET)"
	@$(MAKE) test-mysql
	@$(MAKE) test-postgres

.PHONY: coverage
coverage: ## Generate test coverage report
	@echo "$(BLUE)Generating coverage report...$(RESET)"
	@$(GO) test $(GOFLAGS) -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(RESET)"

.PHONY: coverage-integration
coverage-integration: ## Generate integration test coverage report
	@echo "$(BLUE)Generating integration test coverage report...$(RESET)"
	@TEST_COVERAGE=true $(MAKE) test-integration
	@$(GO) tool cover -html=coverage.out -o coverage-integration.html
	@echo "$(GREEN)Integration coverage report generated: coverage-integration.html$(RESET)"

.PHONY: bench
bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	$(GO) test $(GOFLAGS) -bench=. -benchmem ./...

.PHONY: lint
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not found, skipping linting$(RESET)"; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(RESET)"
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	$(GO) vet ./...

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -f coverage.out coverage.html coverage-integration.html
	@$(GO) clean -cache -testcache

.PHONY: examples
examples: ## Run all examples
	@echo "$(BLUE)Running examples...$(RESET)"
	@for dir in examples/*/; do \
		if [ -f "$$dir"*_example.go ]; then \
			echo "$(YELLOW)Running example in $$dir$(RESET)"; \
			cd "$$dir" && $(GO) run *_example.go && cd ../..; \
		fi; \
	done

.PHONY: build
build: ## Build the library
	@echo "$(BLUE)Building library...$(RESET)"
	$(GO) build $(GOFLAGS) ./...

.PHONY: install
install: ## Install the library
	@echo "$(BLUE)Installing library...$(RESET)"
	$(GO) install $(GOFLAGS) ./...

.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "$(GREEN)All checks passed!$(RESET)"

.PHONY: ci
ci: deps check coverage ## Run CI pipeline
	@echo "$(GREEN)CI pipeline completed!$(RESET)"

# Database setup helpers
.PHONY: setup-mysql
setup-mysql: ## Setup MySQL database for testing
	@echo "$(BLUE)Setting up MySQL database...$(RESET)"
	@if command -v mysql >/dev/null 2>&1; then \
		mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) $(if $(DB_PASSWORD),-p$(DB_PASSWORD)) -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME);"; \
		echo "$(GREEN)MySQL database setup complete$(RESET)"; \
	else \
		echo "$(RED)MySQL client not found. Please install mysql-client or use Docker.$(RESET)"; \
		exit 1; \
	fi

.PHONY: setup-postgres
setup-postgres: ## Setup PostgreSQL database for testing
	@echo "$(BLUE)Setting up PostgreSQL database...$(RESET)"
	@if command -v psql >/dev/null 2>&1; then \
		PGPASSWORD=$(PG_PASSWORD) psql -h$(PG_HOST) -p$(PG_PORT) -U$(PG_USER) -d$(PG_DB) -c "SELECT 1;" >/dev/null 2>&1; \
		echo "$(GREEN)PostgreSQL database setup complete$(RESET)"; \
	else \
		echo "$(RED)PostgreSQL client not found. Please install postgresql-client or use Docker.$(RESET)"; \
		exit 1; \
	fi

# Quick test shortcuts
.PHONY: quick-test
quick-test: ## Quick test run (unit tests only)
	@echo "$(BLUE)Running quick test...$(RESET)"
	$(GO) test $(GOFLAGS) -short ./...

.PHONY: full-test
full-test: test test-integration ## Run all tests (unit + integration)
	@echo "$(GREEN)All tests completed!$(RESET)"

# Development helpers
.PHONY: watch
watch: ## Watch for changes and run tests (requires fswatch)
	@echo "$(BLUE)Watching for changes...$(RESET)"
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . | xargs -n1 -I{} $(MAKE) quick-test; \
	else \
		echo "$(YELLOW)fswatch not found. Install it for file watching.$(RESET)"; \
		echo "$(YELLOW)On macOS: brew install fswatch$(RESET)"; \
		echo "$(YELLOW)On Ubuntu: sudo apt-get install fswatch$(RESET)"; \
	fi

.PHONY: dev
dev: deps ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(RESET)"
	@echo "$(GREEN)Development environment ready!$(RESET)"
	@echo "$(YELLOW)Run 'make help' to see available commands$(RESET)" 