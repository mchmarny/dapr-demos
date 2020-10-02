RELEASE_VERSION  =v0.11.1
SERVICE_NAME    ?=tweet-viewer
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: tidy test debug build run jsonevent xmlevent binevent image lint clean tag
all: help

tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

test: tidy ## Tests the entire project 
	go test -count=1 -race ./...

debug: tidy ## Runs uncompiled code in Dapr
	dapr run \
        --app-id $(SERVICE_NAME) \
        --app-port 8084 \
        --app-protocol http \
        --dapr-http-port 3500 \
        --components-path ./config \
		--log-level debug \
        go run main.go

post: ## Posts sample tweet 
	curl -v -d @tweet.json \
		 -H "Content-type: application/json" \
		 http://localhost:3500/v1.0/publish/processed-tweets-pubsub/processed-tweets

image: tidy ## Builds and publish docker image 
	docker build --build-arg APP_VERSION=$(RELEASE_VERSION) -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

deploy: ## Deploys prebuild image to k8s using currently selected context 
	kubectl apply -f k8s/source-pubsub.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl rollout restart deployment/tweet-viewer
	kubectl rollout status deployment/tweet-viewer

api-patch: ## Patch API gatewaty to add route map
	kubectl patch ingress ingress-rules --type json -p "$(cat k8s/ingress.json)"


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
