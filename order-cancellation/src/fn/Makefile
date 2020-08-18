APP_NAME         =auditor
RELEASE_VERSION  =v0.1.7
DOCKER_USERNAME ?=$(DOCKER_USER)

.PHONY: test mod run dapr lint image tag clean help
all: test

restore: ## Restores the project
	dotnet restore

run: ## Runs function locally 
	dapr run --app-id fn \
             --app-port 3001 \
             --port 3501 \
             --components-path config \
             env TopicName=cancellations StateStore="fn-store" func host start

event: ## Publishes event onto the Dapr topic 
	curl -d @./config/meta.json \
		 -H "Content-type: application/json" \
		 http://localhost:3501/v1.0/publish/cancellations

test: restore ## Tests the entire project 
	dotnet test

image: restore ## Builds docker image and publishes it to Dockerhub
	docker build -t $(DOCKER_USERNAME)/$(APP_NAME):$(RELEASE_VERSION) .
	docker push $(DOCKER_USERNAME)/$(APP_NAME):$(RELEASE_VERSION)


help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
