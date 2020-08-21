RELEASE_VERSION  =v0.10.0
SERVICE_NAME    ?=grpc-echo-service
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: tidy test debug build run call image lint clean tag
all: help

tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

test: tidy ## Tests the entire project 
	go test -count=1 -race ./...

debug: tidy ## Runs uncompiled code in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-port 50002 \
		--app-protocol grpc \
		--dapr-http-port 3500 \
        --components-path ./config \
        go run main.go

build: tidy ## Builds local release binary
	CGO_ENABLED=0 go build -a -tags netgo -mod vendor -o bin/$(SERVICE_NAME) .

run: build ## Builds binary and runs it in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-port 50002 \
		--app-protocol grpc \
		--dapr-http-port 3500 \
        --components-path ./config \
        bin/$(SERVICE_NAME) 

call: ## Invokes service through Dapr API 
	curl -d '{ "message": "ping" }' \
     -H "Content-type: application/json" \
     "http://localhost:3500/v1.0/invoke/$(SERVICE_NAME)/method/echo"

image: tidy ## Builds and publish docker image 
	docker build -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

deploy: ## Deploys prebuild image to k8s using currently selected context 
	kubectl apply -f deployment.yaml
	kubectl get pods


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
