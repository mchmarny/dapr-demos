RELEASE_VERSION  =v0.10.3
SERVICE_NAME    ?=sentiment-scorer
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: tidy debug invoke image deploy call lint clean tag
all: help

tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

run: tidy ## Runs uncompiled code in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-port 60005 \
		--app-protocol grpc \
		--dapr-http-port 3505 \
        --components-path ./config \
		--log-level debug \
        go run main.go

invoke: ## Invokes service through Dapr API 
	curl -d '{ "text": "dapr is the best" }' \
     -H "Content-type: application/json" \
     "http://localhost:3505/v1.0/invoke/$(SERVICE_NAME)/method/sentiment"
	curl -d '{ "text": "spinach is the worst" }' \
     -H "Content-type: application/json" \
     "http://localhost:3505/v1.0/invoke/$(SERVICE_NAME)/method/sentiment"
	curl -d '{ "text": "something is green" }' \
     -H "Content-type: application/json" \
     "http://localhost:3505/v1.0/invoke/$(SERVICE_NAME)/method/sentiment"

image: tidy ## Builds and publish docker image 
	docker build -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

deploy: ## Deploys prebuild image to k8s using currently selected context 
	kubectl apply -f deployment.yaml
	kubectl rollout restart deployment/sentiment-scorer
	kubectl rollout restart deployment/nginx-ingress-nginx-controller
	kubectl rollout status deployment/nginx-ingress-nginx-controller

api-invoke: ## Invokes service through Dapr API 
	$(eval API_TOKEN=$(shell kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode))
	curl -d '{ "text": "dapr is the best" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: $(API_TOKEN)" \
     "https://api.cloudylabs.dev/v1.0/invoke/$(SERVICE_NAME)/method/sentiment"

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
