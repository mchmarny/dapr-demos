.PHONY: tidy run help
all: help

tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy

run: ## Runs uncompiled code it in Dapr
	dapr run \
	     --app-id cron-demo \
	     --protocol http \
	     --app-port 8080 \
	     --components-path ./config \
	     --log-level debug \
	     go run *.go

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
