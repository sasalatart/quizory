GOCMD = go
JSCMD = pnpm
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARIES_DIR = out
CLIENT_DIR = client
GENERATE_FLAG =

all: help

help:
	@echo "install        : Installs the dependencies (Go & JS)."
	@echo "lint           : Runs linters (golangci-lint & eslint)."
	@echo "migrate        : Runs database migrations."
	@echo "codegen        : Runs codegen tools."
	@echo "build          : Builds the project."
	@echo "clean          : Removes the outputs generated by the build command."
	@echo "test           : Runs tests."
	@echo "dev            : Runs the local Docker infra, client and API server in dev mode."
	@echo "docker-image   : Builds the backend Docker image."

install-client:
	cd $(CLIENT_DIR) && $(JSCMD) install

install-go:
	$(GOCMD) install github.com/volatiletech/sqlboiler/v4@latest && \
	$(GOCMD) install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest && \
	$(GOCMD) mod tidy

install: install-go install-client

lint-client:
	cd $(CLIENT_DIR) && $(JSCMD) lint

lint-go:
	golangci-lint run

lint: lint-go lint-client

migrate:
	$(GOCMD) run ./cmd/migrate

codegen:
	$(GOCMD) run ./cmd/codegen && \
	cd $(CLIENT_DIR) && $(JSCMD) codegen

build-client:
	cd $(CLIENT_DIR) && $(JSCMD) build

build-go:
	$(GOBUILD) -o $(BINARIES_DIR)/api -v ./cmd/api

build: install build-go build-client

clean:
	rm -rf $(BINARIES_DIR) && rm -rf $(CLIENT_DIR)/dist

test:
	$(GOTEST) ./...

docker-infra-dev:
	docker-compose -f infra/docker/docker-compose.dev.yml up

client-dev:
	cd client && $(JSCMD) dev

api-dev:
	$(GOCMD) run ./cmd/api $(GENERATE_FLAG)

dev:
	$(MAKE) docker-infra-dev & $(MAKE) client-dev & $(MAKE) api-dev & wait

docker-image:
	docker build -t sasalatart/quizory-api -f ./infra/docker/Dockerfile .
