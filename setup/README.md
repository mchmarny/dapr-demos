## Dapr setup

An opinionated Dapr deployment on Kubernetes clusters. I often use it for the deployment of my demos. It includes:

* Latest version of Dapr
* Metrics Monitoring
  * [Prometheus](https://prometheus.io/) for metrics aggregation
  * [Grafana](https://grafana.com/) for metrics visualization with Dapr monitoring dashboards
* Log Management
  * [Fluentd](https://www.fluentd.org/) for log collection and forwarding
  * [Elasticsearch](https://www.elastic.co/) for log aggregation and query execution
  * [Kibana](https://www.elastic.co/products/kibana) for tull0text log query and visualization
* Distributed Tracing
  * [Jaeger](https://www.jaegertracing.io/) for capturing traces, latency and dependency trace viewing
* Ingress and TLS termination
  * [ngnx](https://nginx.org/en/) for ingress controller and TLS to service mapping 
  * [letsencrypt](https://letsencrypt.org/) as certificate provider
  
## Prerequisites

* 1.15+ Kubernates cluster (if you don't have one, see `make cluster` option to set one up on AKS)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Helm 3](https://helm.sh/docs/intro/install/)
* [certbot](https://certbot.eff.org/) to generate wildcard certificates 

## How to use it

This section assumes you already have a Kubernates cluster. If not, see the [Create Cluster](#create-cluster) section below

> Note, Currently this only supports default namespace deployments

Start by updating the variables at the top of the makefile:

* `DOMAIN` - The root of the domain for which you will be creating wildcard certificates
* `API_TOKEN` - Dapr public API token 


To create a wildcard certificate using [Let's Encrypt](https://letsencrypt.org/) run

```shell
make certs
```

To install Dapr and all observability components as well as configure TLS run 

```shell
make dapr
```

To test your deployment and to validate the API Authentication run

```shell
make test
```

To get access to the observability dashboards run

> Note, this action will open Kibana (port 5601), Grafana (port 8080), Zipkin (port 9411) dashboards in your default browser. If these ports are already being used you may have to edit the `forwards` action in `Makefile`

```shell
make forwards
```

### Setup Metrics Dashboards

To create Prometheus data source and import Dapr dashboards run 

```shell
make metricdash
```

To stop forwarding the above ports run `make unforward`

To login to the Grafana UI (http://localhost:8080) you will admin password, you can get it using the `make metricpass` action

### Setup Log Indexes 

To create Kibana index run 

```shell
make logindex
```

To login to Kibana UI (http://localhost:5601) forward the observability ports `make forwards`

## Create Cluster

If you don't already have a cluster, you can create one on AKS. Start by updating these variables at the `Makefile`

* `CLUSTER_NAME` - Used in cluster creation only 
* `NODE_COUNT` - NUmber of nodes in the cluster default pool
* `NODE_TYPE` - VM type used for the nodes in default pool 

This action assumes your default Azure resource group and location are already defined. If not, run

```shell
az account set --subscription <id or name>
az configure --defaults location=<preferred location> group=<preferred resource group>
```

To create a demo cluster on AKS run

```shell
make clusterup
```

To delete the cluster and all of its resources run `make clusterdown`

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.