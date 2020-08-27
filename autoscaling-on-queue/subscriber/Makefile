RELEASE_VERSION  =v0.10.3
SERVICE_NAME    ?=autoscaling-on-queue
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: all
all: help

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

.PHONY: run
run: tidy ## Runs uncompiled code in Dapr
	dapr run \
      --app-id $(SERVICE_NAME) \
      --app-port 8089 \
      --app-protocol http \
      --dapr-http-port 3500 \
      --components-path ./config \
      --log-level debug \
      go run main.go

.PHONY: local-event
local-event: ## Publishes sample messages in JSON formats to Dapr pubsub API 
	curl -v -d '{ "id": "$(shell uuidgen)", "max": 999999, "Time": $(shell date +%s) }' \
     -H "Content-type: application/json" \
     "http://localhost:3500/v1.0/publish/autoscaling-kafka-queue/primes"

.PHONY: image
image: tidy ## Builds and publishes docker image 
	docker build -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

.PHONY: deploy
deploy: ## Deploys prebuild image to k8s using currently selected context 
	kubectl apply -f k8s/kafka.yaml
	kubectl apply -f k8s/keda.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl rollout restart deployment/$(SERVICE_NAME)
	kubectl rollout restart deployment/nginx-ingress-nginx-controller
	kubectl rollout status deployment/nginx-ingress-nginx-controller

.PHONY: restart
restart: ## Restarts k8s deployment
	kubectl rollout restart deployment/$(SERVICE_NAME)

.PHONY: api-event
api-event: ## Publishes sample JSON message to Dapr pubsub API 
	$(eval API_TOKEN=$(shell kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode))
	curl -v -d '{ "id": "$(shell uuidgen)", "max": 100, "Time": $(shell date +%s) }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: $(API_TOKEN)" \
     "https://api.cloudylabs.dev/v1.0/publish/autoscaling-kafka-queue/primes"


.PHONY: kafka-port
kafka-port: ## Forwards cluster Kafka port locally
	kubectl port-forward svc/kafka -n data 9092 &
	@echo use 'pkill kubectl -9' to stop 

.PHONY: kafka-portless
kafka-portless: ## Stops forwarding cluster Kafka port locally
	@echo use 'pkill kubectl -9' to stop 

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
