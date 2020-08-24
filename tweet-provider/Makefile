RELEASE_VERSION  =v0.10.1
SERVICE_NAME    ?=tweet-provider
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: help tidy build run image lint tag clean
all: help

tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

test: tidy ## Tests the entire project 
	go test -count=1 -race ./...

build: tidy ## Builds local release binary
	CGO_ENABLED=0 go build -a -tags netgo -mod vendor -o bin/$(SERVICE_NAME) .

debug: ## Runs uncompiled code it in Dapr in debug mode
	dapr run \
		--app-id $(SERVICE_NAME) \
        --app-port 8080 \
        --app-protocol http \
        --components-path ./config \
        go run main.go

run: build ## Builds binary and runs it in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
	    --app-port 8080 \
        --app-protocol http \
	    --components-path ./config \
        bin/$(SERVICE_NAME) 

image: tidy ## Builds and publish docker image 
	docker build -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

deploy: ## Deploys prebuild image to k8s using currently selected context 
	kubectl apply -f k8s/pubsub.yaml
	kubectl apply -f k8s/twitter.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl rollout restart deployment/tweet-provider
	kubectl rollout status deployment/tweet-provider

lint: ## Lints the entire project 
	golangci-lint run --timeout=3m

tag: ## Creates release tag 
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

clean: ## Cleans up generated files 
	go clean
	rm -fr ./bin
	rm -fr ./vendor

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
