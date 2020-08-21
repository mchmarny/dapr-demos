DOMAIN       ?=cloudylabs.dev
CLUSTER_NAME ?=demo
NODE_COUNT   ?=3
NODE_TYPE    ?=Standard_D4_v2

.PHONY: all
all: help

.PHONY: clusterlist
clusterlist: ## List all your AKS clusters in the default resource group
	az aks list -o table

.PHONY: cluster
cluster: ## Create AKS cluster (make cluster CLUSTER_NAME=demo)
	az aks create \
		--name $(CLUSTER_NAME) \
		--node-count $(NODE_COUNT) \
		--node-vm-size $(NODE_TYPE) \
		--enable-addons monitoring \
		--generate-ssh-keys
	az aks get-credentials --name $(CLUSTER_NAME)
	
.PHONY: certs
certs: ## Create wildcard TLS certificates using letsencrypt
	sudo certbot certonly --manual --preferred-challenges dns -d "*.$(DOMAIN)"
	sudo cp "/etc/letsencrypt/live/$(DOMAIN)/fullchain.pem" ./cert-ca.pem	
	sudo cp "/etc/letsencrypt/live/$(DOMAIN)/privkey.pem" ./cert-pk.pem
	sudo chmod 644 *.pem

.PHONY: dapr
dapr: ## Install and configures Dapr along with all its dependencies 
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
	kubectl apply -f ./config/fluentd-config.yaml -f ./config/fluentd.yaml
	kubectl apply -f ./config/zipkin-config.yaml -f ./config/zipkin.yaml
	helm install elasticsearch elastic/elasticsearch -n dapr-monitoring
	helm install dapr-prom stable/prometheus -n dapr-monitoring
	helm install grafana stable/grafana -n dapr-monitoring
	helm install kibana elastic/kibana -n dapr-monitoring
	# Deploying ingress...
	$(eval API_TOKEN=$(shell openssl rand -base64 32))
	kubectl create secret generic dapr-api-token --from-literal=token="$(API_TOKEN)"
	helm install nginx ingress-nginx/ingress-nginx -f ./config/annotations.yaml
	kubectl create secret tls tls-secret --key cert-pk.pem --cert cert-ca.pem
	sed "s/DOMAINNAME/${DOMAIN}/g" config/ingress-template.yaml > config/ingress.yaml
	kubectl apply -f ./config/ingress.yaml

.PHONY: dns
dns: ## Check DNS resolution for cluster IP
	dig api.$(DOMAIN)
	$(eval LB_IP=$(shell kubectl get service --selector=app=nginx-ingress,component=controller -o jsonpath='{.items[*].status.loadBalancer.ingress[0].ip}'))
	@echo === DNS CHECK ===
	@echo 
	@echo "Ensure the A record for 'api.${DOMAIN}' in the ANSWER SECTION resolves to:"
	@echo 
	@echo "  ${LB_IP}"
	@echo 
	@echo If not, update DNS with below entry and re-run this test before moving to the next step  
	@echo
	@echo "  Hostname:   *"
	@echo "  IP address: ${LB_IP}"
	@echo "  TTL:        1m"
	@echo

.PHONY: test
test: ## Test deployment and execute Dapr API health checks
	$(eval API_TOKEN=$(shell kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode))
	curl -v \
		 -H "Content-type: application/json" \
		 -H "dapr-api-token: $(API_TOKEN)" \
		 "https://api.$(DOMAIN)/v1.0/healthz"
	@echo 
	@echo === DNS CHECK ===
	@echo Ensure server certificate has:
	@echo 
	@echo "  subject: CN=*.${DOMAIN}"
	@echo "  subjectAltName: host 'api.${DOMAIN}' matched cert's '*.${DOMAIN}'"
	@echo "  SSL certificate verify ok"
	@echo 
	@echo And that the response status from Dapr health checks was '200'
	@echo 
	@echo "  HTTP/2 200"
	@echo

.PHONY: token
token: ## Print Dapr API token 
	$(eval API_TOKEN=$(shell kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode))
	@echo
	@echo Dapr API token is: 
	@echo
	@echo "  ${API_TOKEN}"
	@echo

.PHONY: pass
pass: ## Print Grafana admin password
	$(eval GPASS=$(shell kubectl get secret -n dapr-monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode))
	@echo
	@echo Grafana admin password is: 
	@echo 
	@echo "  ${GPASS}"
	@echo 

.PHONY: ports
ports: ## Forward observability ports
	kubectl port-forward svc/kibana-kibana 5601 -n dapr-monitoring &
	kubectl port-forward svc/grafana 8888:80 -n dapr-monitoring &
	kubectl port-forward svc/zipkin 9411 &
	@echo Ports forwarded:
	@echo
	@echo   kibana  - http://localhost:5601
	@echo   grafana - http://localhost:8888
	@echo   zipkin  - http://localhost:9411
	@echo
	@echo "To stop forwarding run 'make portstop'"
	@echo
	

.PHONY: observe
observe: ports ## Create Prometheus data source, Elastic log index, import Dapr dashboards
	$(eval GPASS=$(shell kubectl get secret -n dapr-monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode))
	curl -X POST -s -k -u "admin:$(GPASS)" \
		 -H "Content-Type: application/json" \
		 -d '{ "name":"Dapr", "type":"prometheus", "url":"http://dapr-prom-prometheus-server.dapr-monitoring", "access":"proxy", "basicAuth":false }' \
	     http://localhost:8888/api/datasources
	curl -X POST -s -k -u "admin:$(GPASS)" \
		 -H "Content-Type: application/json" \
		 -d @config/system-services-dashboard.json \
	     http://localhost:8888/api/dashboards/db
	curl -X POST -s -k -u "admin:$(GPASS)" \
		 -H "Content-Type: application/json" \
		 -d @config/sidecar-dashboard.json \
	     http://localhost:8888/api/dashboards/db
	curl -X POST -s -k -u "admin:$(GPASS)" \
		 -H "Content-Type: application/json" \
		 -d @config/actor-dashboard.json \
	     http://localhost:8888/api/dashboards/db
	curl -X POST -H "kbn-xsrf: true" \
		 -H "Content-Type: application/json" \
		 -d '{"attributes":{"title":"dapr*","timeFieldName":"@timestamp"}}' \
		 "http://localhost:5601/api/saved_objects/index-pattern/dapr"
	curl -X POST -H "kbn-xsrf: true" \
		 -H "Content-Type: application/json" \
		 -d '{"value":"dapr"}' \
		 "http://localhost:5601/api/kibana/settings/defaultIndex"

.PHONY: redis
redis: ## Install redis cluster
	kubectl create ns data
	helm install redis bitnami/redis -n data 
	$(eval REDIS_PASSWORD=$(shell kubectl get secret -n data redis -o jsonpath="{.data.redis-password}" | base64 --decode))
	kubectl create secret generic redis-secret --from-literal=password="$(REDIS_PASSWORD)"

.PHONY: portstop
portstop: ## Stop previously forwarded observability ports 
	pkill kubectl -9	

.PHONY: cleanup
cleanup: clusterlist ## Delete previously created AKS cluster (make cleanup CLUSTER_NAME=demo)
	az aks delete --name $(CLUSTER_NAME)

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
