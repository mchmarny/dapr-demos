CLUSTER_NAME ?=demo
DOMAIN		 ?=cloudylabs.dev
API_TOKEN	 ?=not-a-real-secret
NODE_COUNT   ?=3
NODE_TYPE    ?=Standard_D4_v2

.PHONY: help prep dapr observ trace ingress api test forward
all: help

cluster: ## Creates k8s cluster on Azure using defaults (make cluster CLUSTER_NAME=demo)
	$(info Assumes your default az resource group and locaiton are defined)
	az aks create \
		--name $(CLUSTER_NAME) \
		--node-count $(NODE_COUNT) \
		--node-vm-size $(NODE_TYPE) \
		--enable-addons monitoring \
		--generate-ssh-keys
	az aks get-credentials --name $(CLUSTER_NAME)
	
dapr: ## Registers required Help repos and isntalls the latest release of Dapr (make dapr API_TOKEN=not-a-real-secret)
	helm repo add dapr https://daprio.azurecr.io/helm/v1/repo
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add stable https://kubernetes-charts.storage.googleapis.com
	helm repo add elastic https://helm.elastic.co
	helm repo update
	kubectl create ns dapr-system
	helm install dapr dapr/dapr -n dapr-system --set global.logAsJson=true
	kubectl apply -f ./fluentd-config.yaml -f ./fluentd.yaml
	kubectl apply -f ./zipkin-config.yaml -f ./zipkin.yaml
	kubectl create namespace dapr-monitoring
	helm install elasticsearch elastic/elasticsearch -n dapr-monitoring
	helm install dapr-prom stable/prometheus -n dapr-monitoring
	helm install grafana stable/grafana -n dapr-monitoring
	helm install kibana elastic/kibana -n dapr-monitoring
	kubectl create secret generic dapr-api-token --from-literal=token="$(API_TOKEN)"
	helm install nginx stable/nginx-ingress -f ./annotations.yaml
	$(info Generating TLS sertificates...)
	sudo certbot certonly --manual --preferred-challenges dns -d "*.$(DOMAIN)"
	sudo cp "/etc/letsencrypt/live/$(DOMAIN)/fullchain.pem" ./cert-ca.pem	
	sudo cp "/etc/letsencrypt/live/$(DOMAIN)/privkey.pem" ./cert-pk.pem
	sudo chmod 644 *.pem
	kubectl create secret tls tls-secret --key cert-pk.pem --cert cert-ca.pem
	$(info Now update the ingress.yaml file with the domain for which you generated certificates and run make ingress)

ingress: ## Configures Dapr ingress 
	$(info Assumes you already updated the ingress.yaml file with your domain)
	kubectl apply -f ./ingress.yaml
	kubectl get service --selector=app=nginx-ingress,component=controller \
		-o jsonpath='{.items[*].status.loadBalancer.ingress[0].ip}'
	$(info Now update your DNS with a * record for this IP)

test: ## Executes config API call on Dapr (make test DOMAIN=api.cloudylabs.dev API_TOKEN=not-a-real-secret)
	curl -v \
     -H "Content-type: application/json" \
     -H "dapr-api-token: $(API_TOKEN)" \
     "https://$(DOMAIN)/v1.0/healthz"


metricpass: ## Retrives grafana admin password (make metricpass)
	$(info Grafana admin password:)
	kubectl get secret -n dapr-monitoring grafana \
		-o jsonpath="{.data.admin-password}" | base64 --decode

forward: ## Forwards observability ports (make forward)
	kubectl port-forward svc/kibana-kibana 5601 -n dapr-monitoring &
	kubectl port-forward svc/grafana 8080:80 -n dapr-monitoring &
	kubectl port-forward svc/zipkin 9411:9411 &
	$(info Forwarded ports: kibana=5601, grafana=8080, zipkin=9411)

unforward: ## Stops previously forwarded ports 
	pkill kubectl -9	

cleanup: ## Deletes previously created cluster (make cleanup CLUSTER_NAME=demo)
	$(info Previously created clusters)
	az aks list -o table
	az aks delete --name $(CLUSTER_NAME) \

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
