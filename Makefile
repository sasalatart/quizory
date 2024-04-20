# Makefile for Quizory

# Variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARY_NAME = out/quizory

all: help

help :
	@echo "build    : Builds the project."
	@echo "clean    : Removes the outputs from running the build command."
	@echo "migrate  : Runs database migrations."
	@echo "codegen  : Runs codegen tools."
	@echo "test     : Runs tests."
	@echo "infra    : Starts the infrastructure (docker-compose)."

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Clean the project
clean:
	rm -f $(BINARY_NAME)

# Run database migrations
migrate:
	$(GOCMD) run ./cmd/migrate

# Generate code
codegen:
	$(GOCMD) run ./cmd/codegen

# Run tests
test:
	$(GOTEST) -v ./...

infra:
	docker-compose up -d
