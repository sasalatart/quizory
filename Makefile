GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARIES_DIR = out

all: help

help:
	@echo "docker    	: Starts the local Docker infra."
	@echo "migrate  	: Runs database migrations."
	@echo "codegen  	: Runs codegen tools."
	@echo "build    	: Builds the project."
	@echo "clean    	: Removes the outputs generated by the build command."
	@echo "test     	: Runs tests."
	@echo "aigen-dev	: Runs the AI questions generator in dev mode."

docker:
	docker-compose up

migrate:
	$(GOCMD) run ./cmd/migrate

codegen:
	$(GOCMD) run ./cmd/codegen

build:
	$(GOBUILD) -o $(BINARIES_DIR)/aigen -v ./cmd/aigen

clean:
	rm -f $(BINARIES_DIR)

test:
	$(GOTEST) -v ./...

aigen-dev:
	GO_ENV=dev $(GOCMD) run ./cmd/aigen
