CLUSTER_NAME     ?=demo
CLUSTER_ZONE     ?=us-west1-a
CLUSTER_VERSION  ?=1.17.13-gke.600
NODE_COUNT       ?=2
NODE_COUNT_MIN   ?=1
NODE_COUNT_MAX   ?=5
NODE_TYPE        ?=n2-standard-4

.PHONY: all
all:
	@echo === ACTIVE CONFIGURATION ===
	@echo "CLUSTER_NAME:     ${CLUSTER_NAME}"
	@echo "CLUSTER_ZONE:     ${CLUSTER_ZONE}"
	@echo "CLUSTER_VERSION:  ${CLUSTER_VERSION}"
	@echo "NODE_COUNT:       ${NODE_COUNT}"
	@echo "NODE_COUNT_MIN:   ${NODE_COUNT_MIN}"
	@echo "NODE_COUNT_MAX:   ${NODE_COUNT_MAX}"
	@echo "NODE_TYPE:        ${NODE_TYPE}"
	@echo 
	@echo "Export these as environment variables to change their values"
	@echo 


.PHONY: setup
setup: ## Enables GKE APIs
	gcloud services enable \
		cloudapis.googleapis.com \
		container.googleapis.com \
		containerregistry.googleapis.com

.PHONY: cluster-list
cluster-list: ## List GKE clusters
	gcloud container clusters list

.PHONY: version-list
version-list: ## List Kubernetes versions supported in GKE
	gcloud container get-server-config

.PHONY: cluster
cluster: setup ## Create GKE cluster
	# Create cluster
	gcloud container clusters create $(CLUSTER_NAME) \
	  --addons=HorizontalPodAutoscaling,HttpLoadBalancing \
	  --zone $(CLUSTER_ZONE) \
	  --node-locations $(CLUSTER_ZONE) \
	  --cluster-version $(CLUSTER_VERSION) \
	  --machine-type $(NODE_TYPE) \
	  --num-nodes $(NODE_COUNT) \
	  --min-nodes $(NODE_COUNT_MIN) \
	  --max-nodes $(NODE_COUNT_MAX) \
	  --enable-ip-alias \
	  --enable-stackdriver-kubernetes \
	  --enable-autoscaling \
	  --enable-autorepair \
	  --scopes cloud-platform
	# Get cluster credentials 
	gcloud container clusters get-credentials $(CLUSTER_NAME) \
	  --zone $(CLUSTER_ZONE)
	# Create cluster binding for current user
	kubectl create clusterrolebinding cluster-admin-binding \
	  --clusterrole=cluster-admin \
	  --user=$(gcloud config get-value core/account)

.PHONY: cluster-down
cluster-down: cluster-list ## Delete previously created GKE cluster (make clusterdown CLUSTER_NAME=demo)
	gcloud container clusters delete $(CLUSTER_NAME) \
	  --zone $(CLUSTER_ZONE)

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
