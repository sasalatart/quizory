GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARY_NAME = out/quizory

all: help

help:
	@echo "build    : Builds the project."
	@echo "clean    : Removes the outputs from running the build command."
	@echo "migrate  : Runs database migrations."
	@echo "codegen  : Runs codegen tools."
	@echo "test     : Runs tests."
	@echo "infra    : Starts the infrastructure (docker-compose)."

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	rm -f $(BINARY_NAME)

migrate:
	$(GOCMD) run ./cmd/migrate

codegen:
	$(GOCMD) run ./cmd/codegen

test:
	$(GOTEST) -v ./...

infra:
	docker-compose up -d
