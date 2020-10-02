RELEASE_VERSION  =v0.11.1
SERVICE_NAME    ?=service-format-converter
DOCKER_USERNAME ?=$(DOCKER_USER)
TARGET_FORMAT    =xml
TARGET_SERVICE   =grpc-echo-service
TARGET_METHOD    =echo

.PHONY: all
all: help

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

.PHONY: tidy
test: tidy ## Tests the entire project 
	go test -count=1 -race ./...

.PHONY: run
run: tidy ## Runs uncompiled code in Dapr (make run TARGET_FORMAT=csv)
	TARGET_FORMAT=$(TARGET_FORMAT) \
	TARGET_SERVICE=$(TARGET_SERVICE) \
	TARGET_METHOD=$(TARGET_METHOD) \
	dapr run \
        --app-id $(SERVICE_NAME) \
        --app-port 60012 \
        --app-protocol grpc \
        --components-path ./config \
		--log-level debug \
        go run main.go

.PHONY: image
image: tidy ## Builds and publish docker image 
	docker build -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

.PHONY: deploy
deploy: ## Deploys prebuild image to k8s using currently selected context 
	kubectl apply -f k8s/components.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl rollout restart deployment/$(SERVICE_NAME)

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
