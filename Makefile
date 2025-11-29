.PHONY: build run test lint vet check-all docker-up docker-down

APP_NAME=base_app
APP_CMD_PATH=./cmd/app

# Go commands
GO_BUILD=go build -o $(APP_NAME) $(APP_CMD_PATH)
GO_RUN=go run $(APP_CMD_PATH)
GO_TEST=go test -v ./...
GO_LINT=golangci-lint run
GO_VET=go vet ./...

# Docker commands
DOCKER_COMPOSE=docker-compose -f docker-compose.yml

# Ogen command
OGEN=$(HOME)/go/bin/ogen

.PHONY: build
build:
	@echo "Building the application..."
	@$(GO_BUILD)

.PHONY: run
run:
	@echo "Running the application..."
	@$(GO_RUN)

.PHONY: run-local
run-local:
	@echo "Running the application with local config..."
	@$(GO_RUN) -config=configs/config.local.yaml

.PHONY: run-prepare
run-prepare:
	@echo "Running database migrations..."
	@$(GO_RUN) -mode=Prepare

.PHONY: gen-sqlc
gen-sqlc:
	@echo "Generating SQLC code..."
	@sqlc -f ./contracts/pgsql/sqlc/config.yaml generate

.PHONY: generate
generate:
	@echo "Generating code from OpenAPI spec..."
	@rm -f internal/handler/http/v1/*_gen.go
	@$(OGEN) --target internal/handler/http/v1 --package v1 --clean contracts/http/v1/base_app.yaml

.PHONY: test
test:
	@echo "Running tests..."
	@$(GO_TEST)

.PHONY: lint
lint:
	@echo "Running linter..."
	@$(GO_LINT)

.PHONY: vet
vet:
	@echo "Running go vet..."
	@$(GO_VET)

.PHONY: check-all
check-all: test lint vet gen-sqlc
	@echo "All checks passed!"

.PHONY: migrate-version
migrate-version: install-migrate ## Show current migration version
	@echo "Current migration version:"
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

.PHONY: migrate-status
migrate-status: migrate-version ## Alias for migrate-version

.PHONY: docker-up
docker-up:
	@echo "Starting Docker infrastructure..."
	@$(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down:
	@echo "Stopping Docker infrastructure..."
	@$(DOCKER_COMPOSE) down
