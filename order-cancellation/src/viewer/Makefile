APP_NAME         =order-viewer
RELEASE_VERSION  =v0.11.1
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: test mod run dapr lint image tag clean help
all: test

tidy: ## Updates the go modules and vendors 
	go mod tidy
	go mod vendor

test: mod ## Tests the entire project 
	go test -v -count=1 -race ./...
	# go test -v -count=1 -run NameOfSingleTest ./...

run: mod ## Runs the uncompiled code by itsef 
	go run -v main.go

build: mod ## Build app binary locally 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags netgo -ldflags \
    "-w -extldflags '-static' -X main.AppVersion=$(RELEASE_VERSION)" \
    -mod vendor -o ./bin/service .

dapr: mod ## Runs the uncompiled code in Dapr
	dapr run --app-id viewer \
			 --app-port 8083 \
			 --app-protocol http \
			 --dapr-http-port 3500 \
			 --components-path ./config \
			 go run main.go

image: mod ## Builds docker image and publishes it to Dockerhub
	docker build --build-arg APP_VERSION=$(RELEASE_VERSION) \
	-t $(DOCKER_USERNAME)/$(APP_NAME):$(RELEASE_VERSION) .
	docker push $(DOCKER_USERNAME)/$(APP_NAME):$(RELEASE_VERSION)

lint: ## Lints the entire project
	golangci-lint run --timeout=3m

tag: ## Creates release tag 
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

clean: ## Cleans go and generated files in ./dapr/proto/
	go clean
	rm -fr ./vendor

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
