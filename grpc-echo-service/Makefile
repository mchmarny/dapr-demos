RELEASE_VERSION  =v1.0.1
SERVICE_NAME    ?=echo
IMAGE_NAME      ?=grpc-echo-service
IMAGE_TAG       ?=v0.1.4
IMAGE_OWNER     ?=$(shell git config --get user.username)

.PHONY: all
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
		--app-id $(SERVICE_NAME) \
		--app-port 60002 \
		--app-protocol grpc \
		--dapr-http-port 3500 \
        go run main.go

.PHONY: invoke
invoke: ## Invokes service through Dapr API 
	curl -d '{ "message": "ping" }' \
     -H "Content-type: application/json" \
     "http://localhost:3500/v1.0/invoke/$(SERVICE_NAME)/method/echo"

.PHONY: image
image: tidy ## Builds and publish image 
	docker build -t "ghcr.io/$(IMAGE_OWNER)/$(IMAGE_NAME):$(IMAGE_TAG)" .
	docker push "ghcr.io/$(IMAGE_OWNER)/$(IMAGE_NAME):$(IMAGE_TAG)"

.PHONY: call
call: ## Invokes service through Dapr API 
	$(eval API_TOKEN=$(shell kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode))
	curl -d '{ "message": "ping" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: $(API_TOKEN)" \
     "https://api.cloudylabs.dev/v1.0/invoke/$(SERVICE_NAME)/method/echo"

.PHONY: lint
lint: ## Lints the entire project 
	golangci-lint run --timeout=3m

.PHONY: tag
tag: ## Creates release tag 
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

.PHONY: clean
clean: ## Cleans up generated files 
	go clean
	rm -fr ./bin
	rm -fr ./vendor

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
