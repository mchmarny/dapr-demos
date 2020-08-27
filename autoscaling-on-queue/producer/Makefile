RELEASE_VERSION  =v0.10.5
SERVICE_NAME    ?=autoscaling-producer
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

.PHONY: exec
exec: ## Execute producer 
	kubectl apply -n data -f conf/producer.yaml
	kubectl wait -n data --for=condition=ready pod autoscaling-producer --timeout=120s
	kubectl logs autoscaling-producer -n data -f

.PHONY: stop
stop: ## Stod producer 
	kubectl delete -n data -f conf/producer.yaml

.PHONY: lint
lint: ## Lints the entire project 
	golangci-lint run --timeout=3m

.PHONY: tag
tag: ## Creates release tag 
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

.PHONY: kafka-up
kafka-up: ## Creates local Kafka server 
	docker-compose -f ./conf/local-kafka.yaml up -d

.PHONY: kafka-down
kafka-down: ## Creates local Kafka server 
	docker-compose -f ./conf/local-kafka.yaml down

.PHONY: kafka-port
kafka-port: ## Forwards cluster Kafka port locally
	kubectl port-forward svc/kafka -n data 9092 &
	@echo use 'pkill kubectl -9' to stop 

.PHONY: kafka-portless
kafka-portless: ## Stops forwarding cluster Kafka port locally
	@echo use 'pkill kubectl -9' to stop 

.PHONY: kafka-client
kafka-client: ## Creates Kafka client 
	kubectl apply -n data -f conf/kafka-client.yaml
	kubectl wait -n data --for=condition=ready pod kafka-client --timeout=120s


.PHONY: kafka-sub
kafka-sub: kafka-client ## Subscribes to the Kafka topic 
	kubectl -n data exec -ti kafka-client -- ./bin/kafka-console-consumer.sh \
		--topic prime-requests \
		--from-beginning \
		--bootstrap-server kafka:9092

.PHONY: kafka-topic-count
kafka-topic-count: kafka-client ## Counts number of messages on the topic
	kubectl -n data exec -ti kafka-client -- ./bin/kafka-run-class.sh kafka.tools.GetOffsetShell \
		--broker-list kafka:9092 \
		--topic prime-requests \
		--time -1

.PHONY: kafka-topic-retention
kafka-topic-retentiont: kafka-client ## Counts number of messages on the topic to 1 min 
	kubectl -n data exec -ti kafka-client -- ./bin/kafka-configs.sh \
		--zookeeper kafka-zookeeper:2181 \
		--alter \
		--entity-type topics \
		--entity-name prime-requests \
		--add-config retention.ms=60000

.PHONY: kafka-topic
kafka-topic: kafka-client ## Creates primes topic used for autoscaling demo 
	kubectl -n data exec -ti kafka-client -- ./bin/kafka-topics.sh \
		--zookeeper kafka-zookeeper:2181 \
		--topic prime-requests \
		--create \
		--partitions 1 \
		--replication-factor 1 \
		--config retention.ms=120000

.PHONY: kafka-topic-describe
kafka-topic-describe: kafka-client ## Describes primes topic used for autoscaling demo 
	kubectl -n data exec -ti kafka-client -- ./bin/kafka-topics.sh \
		--zookeeper kafka-zookeeper:2181 \
		--describe \
		--topic prime-requests
		
.PHONY: clean
clean: ## Cleans up generated files 
	go clean
	rm -fr ./bin
	rm -fr ./vendor

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
