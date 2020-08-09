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

> Note, Currently this only supports default namespace deployments

To find out all the operations run `make` or `make help`

```shell
clusterup                      Create k8s cluster in AKS (make cluster CLUSTER_NAME=demo)
certs                          Create wildcard certificates using letsencrypt
dapr                           Install Dapr into current kubeconfig selected context
test                           Execute Dapr API health call
metricpass                     Retrieve grafana admin password
forwards                       Forward observability ports
unforward                      Stop previously forwarded ports
clusterdown                    Delete previously created AKS cluster (make cleanup CLUSTER_NAME=demo)
help                           Display available commands
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.