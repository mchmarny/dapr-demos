.PHONY: help, workflowport, fnport, postmeta, postmail, postcancel, postgateway
all: help

workflowport: ## Forwards Dapr Workflow locally 
	kubectl port-forward deployment/dapr-workflows-host 3500:3500

fnport: ## Forwards Dapr Functions locally 
	kubectl port-forward deployment/auditor 3500:3500

postmeta: ## Invokes cancel on the workflow 
	curl -v -d @data/meta.json -H "Content-type: application/json" \
		 http://localhost:3500/v1.0/publish/cancellations

postmail: ## Invokes email on the workflow 
	curl -v -d @data/email.json -H "Content-type: application/json" \
		 http://localhost:3500/v1.0/bindings/email

postcancel: ## Invokes cancel on the workflow 
	curl -v -d @data/cancellation.json \
		 -H "Content-type: application/json" \
		 http://localhost:3500/v1.0/invoke/workflows/method/order-cancel

postgateway: ## Invokes cancel on the gateway 
		curl -v -d @data/cancellation.json \
			-H "Content-type: application/json" \
			-H "dapr-api-token: ${DEMO_TOKEN}" \
			http://api.thingz.io/order

	
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
