RELEASE_VERSION  =v0.10.2
SERVICE_NAME    ?=tweet-processor
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: tidy test debug build run jsonevent xmlevent binevent image lint clean tag
all: help

tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

debug: tidy ## Runs uncompiled code in Dapr
	dapr run \
        --app-id $(SERVICE_NAME) \
        --app-port 60002 \
        --app-protocol grpc \
        --dapr-http-port 3500 \
        --components-path ./config \
		--log-level debug \
        go run main.go

image: tidy ## Builds and publish docker image 
	docker build -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

deploy: ## Deploys prebuild image to k8s using currently selected context 
	kubectl apply -f k8s/source-pubsub.yaml
	kubectl apply -f k8s/result-pubsub.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl rollout restart deployment/tweet-processor
	kubectl rollout status deployment/tweet-processor

event: ## Publishes sample JSON message to Dapr pubsub API 
	$(eval API_TOKEN=$(shell kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode))
	curl -d '{ "from": "John", "to": "Lary", "message": "hi" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: $(API_TOKEN)" \
     "https://api.cloudylabs.dev/v1.0/publish/grpc-events/messages"

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
