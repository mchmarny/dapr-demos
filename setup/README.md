## WIP: Dapr Setup  

> Note, this is work in progress. Currently only supports default namespace deployments 

An opinionated Dapr deployment on Kubernetes clusters I often use for the deployment of my demos. It includes:

* Latest version of Dapr
* Metrics Monitoring
  * [Prometheus](https://prometheus.io/) for metrics aggregation
  * [Grafana](https://grafana.com/) for metrics visualization with Dapr monitoring dashboards
* Log Management
  * [Fluentd](https://www.fluentd.org/) for log collection and forwarding
  * [Elasticsearch](https://www.elastic.co/) for log aggregation and query execution
  * [Kibana](https://www.elastic.co/products/kibana) for tull0text log query and visualization
* *Distributed Tracing
  * [Jaeger](https://www.jaegertracing.io/) for capturing traces, latency and dependency trace viewing
* *Ingress and TLS termination
  * [ngnx](https://nginx.org/en/) for ingress controller and TLS to service mapping 
  * [certbot](https://certbot.eff.org/) to generate wildcard certificates 

## Prerequisites

* 1.15+ Kubernates cluster (if you don't have one, see `make cluster` option to set one up on AKS)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Helm 3](https://helm.sh/docs/intro/install/)

## How to use it

To find out all the operations run `make` or `make help`

```shell
make
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.