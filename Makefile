GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
BINARY_NAME = out/quizory

all: help

help:
	@echo "docker    	: Starts the local Docker infra."
	@echo "migrate  	: Runs database migrations."
	@echo "codegen  	: Runs codegen tools."
	@echo "build    	: Builds the project."
	@echo "clean    	: Removes the outputs generated by the build command."
	@echo "test     	: Runs tests."
	@echo "aigen   		: Runs the AI questions generator."

docker:
	docker-compose up

migrate:
	$(GOCMD) run ./cmd/migrate

codegen:
	$(GOCMD) run ./cmd/codegen

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

aigen:
	GO_ENV=dev $(GOCMD) run ./cmd/aigen
