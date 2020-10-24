.PHONY: all
all: help


.PHONY: ping
ping: ## Invokes ping method on app1 service through Dapr API 
	curl -i -d '{ "message": "hi" }' \
        -H "Content-type: application/json" \
        "http://localhost:3500/v1.0/invoke/app1/method/ping"


.PHONY: message
message: ## Publish message directly to the Dapr Pub/Sub API
	curl -i -d '{ "message": "hi" }' \
        -H "Content-type: application/json" \
        "http://localhost:3500/v1.0/publish/pubsub/messages"


.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
