IMAGE_NAME  ?=hardened-app2
IMAGE_TAG   ?=v0.1.4
IMAGE_OWNER ?=$(shell git config --get user.username)
APP_ID      ?=app2
APP_PORT    ?=8082

all: help

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

.PHONY: test
test: tidy ## Tests the entire project 
	go test -count=1 -race ./...

.PHONY: run
run: tidy ## Runs uncompiled code in Dapr
	dapr run \
      --app-id $(APP_ID) \
      --app-port $(APP_PORT) \
      --app-protocol http \
      --components-path ../../config \
      --log-level debug \
      go run main.go

.PHONY: image
image: tidy ## Builds and publishes image 
	docker build -t "ghcr.io/$(IMAGE_OWNER)/$(IMAGE_NAME):$(IMAGE_TAG)" .
	docker push "ghcr.io/$(IMAGE_OWNER)/$(IMAGE_NAME):$(IMAGE_TAG)"

.PHONY: lint
lint: ## Lints the entire project 
	golangci-lint run --timeout=3m

.PHONY: clean
clean: ## Cleans up generated files 
	go clean
	rm -fr ./bin
	rm -fr ./vendor

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
