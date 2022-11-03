NAME=webinar-scale-up-app
BINARY_NAME={NAME}.out
# Enables support for tools such as https://github.com/rakyll/gotest
TEST_COMMAND ?= go test ./...
# Tags specific for building
GOTAGS ?=
# List all our actual files, excluding vendor
GOPKGS ?= $(shell go list $(FILES) | grep -v /vendor/)
GOFILES ?= $(shell find . -name '*.go' | grep -v /vendor/)
GIT_COMMIT ?= $(shell git rev-parse HEAD)

## dev
dev: ## Run the app in dev mode
	@go run cmd/app/main.go -dev

## Build:
build-app: ## Build your project and put the output binary in out/bin/
	mkdir -p out/bin
	go build -o out/bin/$(BINARY_NAME) cmd/app/main.go

## Docker-Build:
docker-build: ## Docker-build and tag the image with the latest git commit hash
	docker build -t $(NAME):$(GIT_COMMIT) .

## Docker-Run:
docker-run: ## Run the docker image with the latest git commit hash
	export IMAGE_TAG=$(GIT_COMMIT); \
	export IMAGE_NAME=$(NAME); \
	docker-compose -f services/Docker-compose.yaml up -d

docker-stop:
	docker-compose -f services/Docker-compose.yaml down

## Run:
run-app: ## Run your project
	./out/bin/$(BINARY_NAME)

## Clean
clean: ## Remove build related file
	@rm -fr ./bin
	@rm -fr ./out
	@rm -f ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml
.PHONY: clean

mod-download: ## Downloads the Go module.
	@echo "==> Downloading Go module"
	@go mod download
.PHONY: mod-download

tidy: ## Cleans the Go module.
	@echo "==> Tidying module"
	@go mod tidy
.PHONY: tidy

fmt: ## Properly formats Go files and orders dependencies.
	@echo "==> Running gofmt"
	@gofmt -s -w ${GOFILES}
.PHONY: fmt

vet: ## Identifies common errors.
	@echo "==> Running go vet"
	@go vet ./...
.PHONY: vet

test: ## Runs the test suite with VCR mocks enabled.
	@echo "==> Testing ${NAME}"
	@$(TEST_COMMAND) -timeout=30s -parallel=20 -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test

test-race: ## Runs the test suite with the -race flag to identify race conditions, if they exist.
	@echo "==> Testing ${NAME} (race)"
	@$(TEST_COMMAND) -timeout=60s -race -tags="${GOTAGS}" ${GOPKGS} ${TESTARGS}
.PHONY: test-race

help: ## Prints this help menu.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
