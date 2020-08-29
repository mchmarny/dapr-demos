RELEASE_VERSION  =v0.10.1
SERVICE_NAME    ?=queue-outoscaling-producer
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: all
all: help

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy
	go mod vendor

.PHONY: run
run: tidy ## Runs uncompiled code
	go run main.go

.PHONY: image
image: tidy ## Builds and publishes docker image 
	docker build -t "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)" .
	docker push "$(DOCKER_USERNAME)/$(SERVICE_NAME):$(RELEASE_VERSION)"

.PHONY: load
load: ## Execute producer and publish events onto the topic 
	kubectl apply -n data -f config/producer.yaml
	kubectl rollout -n data status deployment/$(SERVICE_NAME)
	kubectl logs -n data -l app=queue-outoscaling-producer -f

.PHONY: scale
scale: ## Scale producer
	kubectl scale -n data --replicas=10 deployment/$(SERVICE_NAME)

.PHONY: stop
stop: ## Stop producer 
	kubectl delete -n data -f config/producer.yaml

.PHONY: lint
lint: ## Lints the entire project 
	golangci-lint run --timeout=3m

.PHONY: tag
tag: ## Creates release tag 
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

.PHONY: kafka-up
kafka-up: ## Creates local Kafka server 
	docker-compose -f ./config/kafka.yaml up -d

.PHONY: kafka-down
kafka-down: ## Creates local Kafka server 
	docker-compose -f ./config/kafka.yaml down

.PHONY: kafka-client
kafka-client: ## Deploys Kafka client pod
	kubectl apply -n data -f config/kafka-client.yaml
	kubectl wait -n data --for=condition=ready pod kafka-client --timeout=120s

.PHONY: kafka-login
kafka-login: kafka-client ## Login to Kafka server 
	kubectl -n data exec -it kafka-client -- /bin/bash

.PHONY: kafka-topic
kafka-topic: kafka-client ## Creates messages topic used for autoscaling demo 
	kubectl -n data exec -it kafka-client -- kafka-topics \
		--zookeeper kafka-cp-zookeeper-headless:2181 \
		--topic messages \
		--create \
		--partitions 5 \
		--replication-factor 1 \
		--if-not-exists
		
.PHONY: clean
clean: ## Cleans up generated files 
	go clean
	rm -fr ./bin
	rm -fr ./vendor

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
