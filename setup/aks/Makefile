CLUSTER_NAME     ?=demo
CLUSTER_VERSION  ?=1.18.8
NODE_COUNT       ?=3
NODE_TYPE        ?=Standard_D4_v2

.PHONY: all
all:
	@echo === ACTIVE CONFIGURATION ===
	@echo "CLUSTER_NAME:     ${CLUSTER_NAME}"
	@echo "CLUSTER_VERSION:  ${CLUSTER_VERSION}"
	@echo "NODE_COUNT:       ${NODE_COUNT}"
	@echo "NODE_TYPE:        ${NODE_TYPE}"
	@echo 
	@echo "Export these as environment variables to change their values"
	@echo 


.PHONY: cluster-list
cluster-list: ## List AKS clusters
	az aks list -o table

.PHONY: version-list
version-list: ## List Kubernetes versions supported in AKS
	az aks get-versions -o json --query 'orchestrators[].orchestratorVersion'

.PHONY: cluster
cluster: ## Create AKS cluster
	az aks create \
		--name $(CLUSTER_NAME) \
		--node-count $(NODE_COUNT) \
		--node-vm-size $(NODE_TYPE) \
		--enable-addons monitoring \
		--kubernetes-version $(CLUSTER_VERSION) \
		--generate-ssh-keys \
		--enable-managed-identity
	az aks get-credentials --name $(CLUSTER_NAME)

.PHONY: cluster-down
cluster-down: cluster-list ## Delete previously created AKS cluster (make clusterdown CLUSTER_NAME=demo)
	az aks delete --name $(CLUSTER_NAME)

.PHONY: ip
ip: ## Create Static IP for the existing AKS cluster
	$(eval CLUSTERRC=$(shell az aks show -n $(CLUSTER_NAME) -o tsv --query nodeResourceGroup))
	az network public-ip create \
		-g $(CLUSTERRC) \
		-n "${CLUSTER_NAME}IP" \
		--sku STANDARD
	az network public-ip show \
		-g $(CLUSTERRC) \
		-n "${CLUSTER_NAME}IP" \
		-o tsv --query ipAddress

.PHONY: node-pool
node-pool: ## Add new AKS node pool
	$(eval NODE_SUFIX=$(shell date +"%d%m"))
	az aks nodepool add \
		--cluster-name $(CLUSTER_NAME) \
		--name nodepool$(NODE_SUFIX) \
		--node-count $(NODE_COUNT) \
		--node-vm-size $(NODE_TYPE) \
		--mode System \
		--no-wait


.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
