.PHONY: all
all: help

.PHONY: apply
apply: ## Applies all and restarts the apps 
	# apply
	kubectl apply -f k8s/ -n hardened
	# restart
	kubectl rollout restart deployment/app1 -n hardened
	kubectl rollout restart deployment/app2 -n hardened
	kubectl rollout restart deployment/app3 -n hardened
	# status
	kubectl rollout status deployment/app1 -n hardened
	kubectl rollout status deployment/app2 -n hardened
	kubectl rollout status deployment/app3 -n hardened

.PHONY: ping
ping: ## Execute ping on the Dapr API 
	$(eval API_TOKEN=$(shell kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode))
	curl -i -d '{ "message": "hello" }' \
      -H "Content-type: application/json" \
      -H "dapr-api-token: $(API_TOKEN)" \
      https://api.thingz.io/v1.0/invoke/app1.hardened/method/ping

.PHONY: count
count: ## Execute counter on the Dapr API 
	$(eval API_TOKEN=$(shell kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode))
	curl -i -d '{ "on": 1603627556200126373, "count": 2 }' \
      -H "Content-type: application/json" \
      -H "dapr-api-token: $(API_TOKEN)" \
      https://api.thingz.io/v1.0/invoke/app2.hardened/method/counter

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
