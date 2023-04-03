SHELL = /bin/bash

# ==============================================================================
# Useful variables

MARIADB_DOCKER_CONTAINER=temporal-tables-mariadb
MYSQL_BIN=docker exec -it $(MARIADB_DOCKER_CONTAINER) mysql -u$(MARIADB_USER) -p$(MARIADB_PASSWORD) -D$(MARIADB_DATABASE)

MARIADB_TEST_DOCKER_CONTAINER=temporal-tables-mariadb-test
TEST_MYSQL_BIN=docker exec -it $(MARIADB_TEST_DOCKER_CONTAINER) mysql -u$(MARIADB_USER) -p$(MARIADB_PASSWORD) -D$(MARIADB_DATABASE)

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==============================================================================
# DB console

.PHONY: mariadb-console
## mariadb-console: launches mariadb local database console
mariadb-console: export MARIADB_DATABASE=hr
mariadb-console: export MARIADB_USER=hr
mariadb-console: export MARIADB_PASSWORD=hr123!
mariadb-console: export MARIADB_ROOT_PASSWORD=123456
mariadb-console: export MARIADB_HOST_NAME=localhost
mariadb-console: export MARIADB_PORT=3310
mariadb-console:
	@ $(MYSQL_BIN)

.PHONY: test-mariadb-console
## test-mariadb-console: launches mariadb test database console
test-mariadb-console: export MARIADB_DATABASE=hr
test-mariadb-console: export MARIADB_USER=hr
test-mariadb-console: export MARIADB_PASSWORD=hr123!
test-mariadb-console: export MARIADB_ROOT_PASSWORD=123456
test-mariadb-console: export MARIADB_HOST_NAME=localhost
test-mariadb-console: export MARIADB_PORT=3311
test-mariadb-console:
	@ $(TEST_MYSQL_BIN)

# ==============================================================================
# DB migrations

.PHONY: create-migration
## create-migration: creates a migration file
create-migration:
	@ if [ -z "$(NAME)" ]; then echo >&2 please set the name of the migration via the variable NAME; exit 2; fi
	@ docker-compose up create-migration

# ==============================================================================
# Tests

.PHONY: test
## test: runs unit tests
test:
	@ go test -cover -v ./... -count=1

.PHONY: coverage
## coverage: run unit tests and generate coverage report in html format
coverage:
	@ go test -coverprofile=coverage.out ./...  && go tool cover -html=coverage.out

.PHONY: test-db-up
## test-db-up: starts test database
test-db-up:
	@ echo "Setting up test MariaDB..."
	@ docker-compose up -d mariadb_test migrate_test
	@ until $(TEST_MYSQL_BIN) -e 'SELECT 1' >/dev/null 2>&1 && exit 0; do \
	  >&2 echo "MariaDB not ready, sleeping for 5 secs..."; \
	  sleep 5 ; \
	done
	@ echo "... MariaDB is up and running!"

.PHONY: int-test
## int-test: runs integration tests
int-test: export MARIADB_DATABASE=hr
int-test: export MARIADB_USER=hr
int-test: export MARIADB_PASSWORD=hr123!
int-test: export MARIADB_ROOT_PASSWORD=123456
int-test: export MARIADB_HOST_NAME=localhost
int-test: export MARIADB_PORT=3311
int-test: test-db-up
	@ go test -v ./test --tags=integration

# ==============================================================================
# Code quality

.PHONY: vet
## vet: runs go vet
vet:
	@ go vet ./...

.PHONY: lint
## lint: runs linter for all packages
lint: 
	@ docker run --rm -v "`pwd`:/workspace:cached" -w "/workspace/." golangci/golangci-lint:latest golangci-lint run

.PHONY: vul-setup
## vul-setup: installs Golang's vulnerability check tool
vul-setup:
	@ if [ -z "$$(which govulncheck)" ]; then echo "Installing Golang's vulnerability detection tool..."; go install golang.org/x/vuln/cmd/govulncheck@latest; fi

.PHONY: vul-check
## vul-check: checks for any known vulnerabilities
vul-check: vul-setup
	@ govulncheck ./...

# ==============================================================================
# Swagger

.PHONY: swagger
## swagger: generates api's documentation
swagger: 
	@ docker run --rm -it -v $(HOME):$(HOME) -w $(PWD) quay.io/goswagger/swagger generate spec -o doc/swagger.json

.PHONY: swagger-ui
## swagger-ui: launches swagger ui
swagger-ui: swagger
	@ docker-compose up swagger-ui

# ==============================================================================
# Execution

.PHONY: run
## run: runs the application
run:
	@ docker-compose up --force-recreate --build api

.PHONY: stop
## stop: stops all containers
stop:
	@ docker-compose down
