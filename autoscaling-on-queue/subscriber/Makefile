RELEASE_VERSION  =v0.10.4
SERVICE_NAME    ?=autoscaling-subscriber
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
      --app-port 60022 \
      --app-protocol grpc \
      --components-path ./config \
      --log-level debug \
      go run main.go		

.PHONY: image
image: tidy ## Builds and publishes docker image 
	docker build -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

.PHONY: deploy
deploy: ## Deploys prebuild image to k8s using currently selected context 
	kubectl apply -f k8s/binding.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl rollout restart deployment/$(SERVICE_NAME)
	kubectl rollout status deployment/$(SERVICE_NAME)

.PHONY: autoscale
autoscale: ## Applies Keda to the deployed service 
	kubectl apply -f k8s/keda.yaml

.PHONY: scale
scale: ## Scales the subscriber manually 
	# kubectl autoscale deployment autoscaling-subscriber --min=1 --max=10  
	kubectl scale deployment/$(SERVICE_NAME) --replicas=10

.PHONY: kafka-port
kafka-port: ## Forwards cluster Kafka port locally
	kubectl port-forward svc/kafka -n kafka 9092 &
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
