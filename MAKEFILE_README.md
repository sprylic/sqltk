# SQL Builder Integration Test Makefile

This Makefile provides easy configuration and execution of integration tests for the SQL builder library. It supports both MySQL and PostgreSQL databases, with options to use local databases or Docker containers.

**Important**: This Makefile runs the existing integration tests in `mysql_integration_test.go` and `postgres_integration_test.go` using the proper database connection strings (`MYSQL_DSN` and `POSTGRES_DSN`).

## Quick Start

1. **Show available commands:**
   ```bash
   make help
   ```

2. **Run unit tests:**
   ```bash
   make test
   ```

3. **Run integration tests with Docker:**
   ```bash
   make test-all-docker
   ```

4. **Run integration tests with local databases:**
   ```bash
   make test-all-integration
   ```

## Integration Tests

The Makefile runs the following integration tests:

- **`TestMySQLIntegration`** - Tests MySQL DDL operations, CRUD operations, and advanced features
- **`TestPostgresIntegration`** - Tests PostgreSQL DDL operations, CRUD operations, and advanced features

These tests automatically:
- Create temporary test databases
- Run comprehensive SQL builder tests
- Clean up test databases when done

## Configuration

### Environment Variables

You can configure the behavior using environment variables:

```bash
# MySQL configuration
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password

# PostgreSQL configuration
export PG_HOST=localhost
export PG_PORT=5432
export PG_USER=postgres
export PG_PASSWORD=your_password
export PG_DB=postgres

# Test configuration
export TEST_TIMEOUT=30s
export TEST_RACE=false
export TEST_COVERAGE=false
export TEST_VERBOSE=true
```

### Configuration File

Copy `config.env.example` to `config.env` and modify as needed:

```bash
cp config.env.example config.env
# Edit config.env with your settings
```

Then source the configuration:

```bash
source config.env
```

## Available Commands

### Basic Testing

- `make test` - Run unit tests
- `make test-integration` - Run integration tests with configured database
- `make test-mysql` - Run MySQL integration tests
- `make test-postgres` - Run PostgreSQL integration tests
- `make test-all-integration` - Run all integration tests (MySQL + PostgreSQL)
- `make quick-test` - Run quick tests (unit tests only)
- `make full-test` - Run all tests (unit + integration)

### Docker-based Testing

- `make docker-mysql` - Start MySQL database in Docker
- `make docker-postgres` - Start PostgreSQL database in Docker
- `make test-docker-mysql` - Run MySQL tests with Docker database
- `make test-docker-postgres` - Run PostgreSQL tests with Docker database
- `make test-all-docker` - Run all integration tests with Docker databases
- `make docker-stop` - Stop all test database containers

### Coverage and Analysis

- `make coverage` - Generate test coverage report
- `make coverage-integration` - Generate integration test coverage report
- `make bench` - Run benchmarks
- `make lint` - Run linter (requires golangci-lint)
- `make vet` - Run go vet
- `make fmt` - Format code

### Development

- `make deps` - Install dependencies
- `make build` - Build the library
- `make install` - Install the library
- `make examples` - Run all examples
- `make clean` - Clean build artifacts
- `make check` - Run all checks (format, vet, lint, test)
- `make ci` - Run CI pipeline
- `make dev` - Setup development environment
- `make watch` - Watch for changes and run tests (requires fswatch)

### Database Setup

- `make setup-mysql` - Setup MySQL database for testing
- `make setup-postgres` - Setup PostgreSQL database for testing

## Usage Examples

### Running Tests with Local Database

1. **Setup local MySQL:**
   ```bash
   make setup-mysql
   make test-mysql
   ```

2. **Setup local PostgreSQL:**
   ```bash
   make setup-postgres
   make test-postgres
   ```

### Running Tests with Docker

1. **Run all tests with Docker databases:**
   ```bash
   make test-all-docker
   ```

2. **Run only MySQL tests with Docker:**
   ```bash
   make test-docker-mysql
   ```

3. **Run only PostgreSQL tests with Docker:**
   ```bash
   make test-docker-postgres
   ```

### Custom Configuration

1. **Override database settings:**
   ```bash
   DB_HOST=my-db-server DB_PORT=3307 make test-mysql
   ```

2. **Enable race detection:**
   ```bash
   TEST_RACE=true make test
   ```

3. **Generate coverage report:**
   ```bash
   TEST_COVERAGE=true make test
   ```

4. **Run with custom timeout:**
   ```bash
   TEST_TIMEOUT=60s make test-integration
   ```

## Database Connection Strings

The Makefile automatically generates the correct connection strings for your integration tests:

- **MySQL**: `MYSQL_DSN="user:password@tcp(host:port)/"`
- **PostgreSQL**: `POSTGRES_DSN="postgres://user:password@host:port/db?sslmode=disable"`

These are the same environment variables that your existing integration tests expect.

## Prerequisites

### Required Software

- **Go** (1.19 or later)
- **Make** (GNU Make)
- **Docker** (for Docker-based testing)

### Optional Software

- **MySQL Client** (for local MySQL testing)
- **PostgreSQL Client** (for local PostgreSQL testing)
- **golangci-lint** (for linting)
- **fswatch** (for file watching)

### Installation

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install make docker.io mysql-client postgresql-client
```

**macOS:**
```bash
brew install make docker mysql postgresql
```

**Install golangci-lint:**
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Troubleshooting

### Docker Issues

1. **Container already exists:**
   ```bash
   make docker-stop
   make docker-mysql
   ```

2. **Port already in use:**
   ```bash
   DOCKER_MYSQL_PORT=3307 make docker-mysql
   ```

3. **Docker permissions:**
   ```bash
   sudo usermod -aG docker $USER
   # Then log out and back in
   ```

### Database Connection Issues

1. **Check MySQL connection:**
   ```bash
   mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD -e "SELECT 1;"
   ```

2. **Check PostgreSQL connection:**
   ```bash
   PGPASSWORD=$PG_PASSWORD psql -h$PG_HOST -p$PG_PORT -U$PG_USER -d$PG_DB -c "SELECT 1;"
   ```

3. **Test integration test connection:**
   ```bash
   MYSQL_DSN="root:password@tcp(localhost:3306)/" go test -tags=integration -run TestMySQLIntegration
   ```

### Test Issues

1. **Clean and retry:**
   ```bash
   make clean
   make test
   ```

2. **Run with verbose output:**
   ```bash
   TEST_VERBOSE=true make test
   ```

3. **Skip integration tests:**
   ```bash
   go test ./... -tags=!integration
   ```

## CI/CD Integration

The Makefile is designed to work well in CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run tests
  run: |
    make ci
    make test-all-docker
```

```yaml
# Example GitLab CI
test:
  script:
    - make deps
    - make check
    - make test-all-docker
```

## Contributing

When adding new tests:

1. **Unit tests:** Use `make test`
2. **Integration tests:** Use `make test-integration`
3. **Database-specific tests:** Use `make test-mysql` or `make test-postgres`
4. **Add test tags:** Use `//go:build integration` for integration tests

## License

This Makefile is part of the SQL Builder library and follows the same license terms. 