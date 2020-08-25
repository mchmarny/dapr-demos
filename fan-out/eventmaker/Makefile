APP_NAME         =eventmaker

all: test

.PHONY: mod
tidy: ## Updates the go modules and vendors all dependancies 
	go mod tidy
	go mod vendor

.PHONY: build
build: mod ## Builds binaries
	CGO_ENABLED=0 go build -o bin/$(APP_NAME) .

.PHONY: events
events: ## Sends mocked events to configured target
	bin/$(APP_NAME) --conn "${EVENT_HUB_CONN}"

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'



