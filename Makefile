GOCMD = go
JSCMD = pnpm
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOTOOL = $(GOCMD) tool
BINARIES_DIR = out
CLIENT_DIR = client

all: help

help:
	@echo "install        : Installs the dependencies (Go & JS)."
	@echo "lint           : Runs linters (golangci-lint & eslint)."
	@echo "migrate        : Runs database migrations."
	@echo "codegen        : Runs codegen tools."
	@echo "build          : Builds the project."
	@echo "clean          : Removes the outputs generated by the build command."
	@echo "test           : Runs tests."
	@echo "coverage       : Runs tests with coverage report."
	@echo "dev            : Runs the local Docker infra, client, and API server in dev mode."
	@echo "docker-image   : Builds the backend Docker image."

install: install-go install-client

install-client:
	cd $(CLIENT_DIR) && $(JSCMD) install

install-go:
	$(GOCMD) install github.com/volatiletech/sqlboiler/v4@latest && \
	$(GOCMD) install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest && \
	$(GOCMD) mod tidy

lint: lint-go lint-client

lint-client:
	cd $(CLIENT_DIR) && $(JSCMD) lint

lint-go:
	golangci-lint run

migrate:
	$(GOCMD) run ./cmd/migrate

codegen:
	PSQL_HOST=localhost $(GOCMD) run ./cmd/codegen && \
	cd $(CLIENT_DIR) && $(JSCMD) codegen

build: install build-go build-client

build-client:
	cd $(CLIENT_DIR) && $(JSCMD) build

build-go:
	$(GOBUILD) -o $(BINARIES_DIR)/api -v ./cmd/api

clean:
	rm -rf $(BINARIES_DIR) && rm -rf $(CLIENT_DIR)/dist && rm -rf $(CLIENT_DIR)/src/generated/api/apis

test:
	$(GOTEST) -race -shuffle=on ./...

coverage:
	$(GOTEST) -coverprofile coverage.out ./... && \
	$(GOTOOL) cover -html coverage.out -o coverage.html && open coverage.html

dev:
	@sh -c '\
		$(MAKE) docker-dev & DOCKER_PID=$$!; \
		$(MAKE) client-dev & CLIENT_PID=$$!; \
		wait $$DOCKER_PID $$CLIENT_PID'

docker-dev:
	docker compose -f infra/docker/docker-compose.dev.yml up

docker-dev-down:
	docker compose -f infra/docker/docker-compose.dev.yml down

client-dev:
	cd $(CLIENT_DIR) && $(JSCMD) dev

docker-image:
	docker build -t sasalatart/quizory-api -f ./infra/docker/Dockerfile .
