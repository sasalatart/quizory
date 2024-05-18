GOCMD = go
JSCMD = pnpm
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARIES_DIR = out
CLIENT_DIR = client

all: help

help:
	@echo "docker    	: Starts the local Docker infra."
	@echo "install  	: Installs the dependencies (Go & JS)."
	@echo "lint  	  	: Runs linters (golangci-lint & eslint)."
	@echo "migrate  	: Runs database migrations."
	@echo "codegen  	: Runs codegen tools."
	@echo "build    	: Builds the project."
	@echo "clean    	: Removes the outputs generated by the build command."
	@echo "test     	: Runs tests."
	@echo "dev      	: Runs the local Docker infra, client and API server in dev mode."

docker:
	docker-compose up

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
	$(GOCMD) run ./cmd/codegen

build-ai-gen:
	$(GOBUILD) -o $(BINARIES_DIR)/aigen -v ./cmd/aigen

build-api:
	$(GOBUILD) -o $(BINARIES_DIR)/api -v ./cmd/api

build-client:
	cd $(CLIENT_DIR) && $(JSCMD) build

build: build-ai-gen build-api build-client

clean:
	rm -rf $(BINARIES_DIR) && rm -rf $(CLIENT_DIR)/dist

test:
	$(GOTEST) ./...

client-dev:
	cd client && $(JSCMD) dev

aigen-dev:
	GO_ENV=dev $(GOCMD) run ./cmd/aigen

api-dev:
	$(GOCMD) run ./cmd/api

dev:
	$(MAKE) docker & $(MAKE) client-dev & $(MAKE) api-dev & wait
