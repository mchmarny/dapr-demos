CLUSTER_NAME ?=demo6
DOMAIN		 ?=chmarny.dev
API_TOKEN	 ?=not-a-real-secret
NODE_COUNT   ?=3
NODE_TYPE    ?=Standard_D4_v2

.PHONY: help cluster certs test dapr metricpass forwards unforward cleanup
all: help

clusterup: ## Create k8s cluster in AKS (make cluster CLUSTER_NAME=demo)
	$(info Assumes your default az resource group and locaiton are defined)
	$(info   az account set --subscription <id or name>)
	$(info   az configure --defaults location=<preferred location> group=<preferred resource group>)
	az aks create \
		--name $(CLUSTER_NAME) \
		--node-count $(NODE_COUNT) \
		--node-vm-size $(NODE_TYPE) \
		--enable-addons monitoring \
		--generate-ssh-keys
	az aks get-credentials --name $(CLUSTER_NAME)
	
certs: ## Create wildcard certificates using letsencrypt
	$(info Generating TLS sertificates...)
	sudo certbot certonly --manual --preferred-challenges dns -d "*.$(DOMAIN)"
	sudo cp "/etc/letsencrypt/live/$(DOMAIN)/fullchain.pem" ./cert-ca.pem	
	sudo cp "/etc/letsencrypt/live/$(DOMAIN)/privkey.pem" ./cert-pk.pem
	sudo chmod 644 *.pem

dapr: ## Install Dapr into current kubeconfig selected context
	# Updating Help repos...
	helm repo add dapr https://daprio.azurecr.io/helm/v1/repo
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add stable https://kubernetes-charts.storage.googleapis.com
	helm repo add elastic https://helm.elastic.co
	helm repo update
	# Installing Dapr...
	kubectl create ns dapr-system
	helm install dapr dapr/dapr -n dapr-system --set global.logAsJson=true
	# Installing Observabiliity...
	kubectl create namespace dapr-monitoring
	kubectl apply -f ./fluentd-config.yaml -f ./fluentd.yaml
	kubectl apply -f ./zipkin-config.yaml -f ./zipkin.yaml
	helm install elasticsearch elastic/elasticsearch -n dapr-monitoring
	helm install dapr-prom stable/prometheus -n dapr-monitoring
	helm install grafana stable/grafana -n dapr-monitoring
	helm install kibana elastic/kibana -n dapr-monitoring
	# Deploying ingress...
	kubectl create secret generic dapr-api-token --from-literal=token="$(API_TOKEN)"
	helm install nginx stable/nginx-ingress -f ./annotations.yaml
	kubectl create secret tls tls-secret --key cert-pk.pem --cert cert-ca.pem
	kubectl apply -f ./ingress.yaml
	# Wating for monitoring to come up....
	kubectl get pods -n dapr-monitoring -w
	kubectl get service --selector=app=nginx-ingress,component=controller \
		-o jsonpath='{.items[*].status.loadBalancer.ingress[0].ip}'
	# Update DNS with an A * record for this IP

test: ## Execute Dapr API health call
	curl -v \
     -H "Content-type: application/json" \
     -H "dapr-api-token: $(API_TOKEN)" \
     "https://api.$(DOMAIN)/v1.0/healthz"

metricpass: ## Retrieve grafana admin password
	$(info Grafana admin password:)
	kubectl get secret -n dapr-monitoring grafana \
		-o jsonpath="{.data.admin-password}" | base64 --decode

forwards: ## Forward observability ports
	$(info Forwarding ports: kibana=5601, grafana=8080, zipkin=9411)
	kubectl port-forward svc/kibana-kibana 5601 -n dapr-monitoring &
	kubectl port-forward svc/grafana 8080:80 -n dapr-monitoring &
	kubectl port-forward svc/zipkin 9411:9411 &
	open http://localhost:5601
	open http://localhost:8080
	open http://localhost:9411

unforward: ## Stop previously forwarded ports 
	pkill kubectl -9	

clusterdown: ## Delete previously created AKS cluster (make cleanup CLUSTER_NAME=demo)
	$(info Previously created clusters)
	az aks list -o table
	az aks delete --name $(CLUSTER_NAME) \

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
