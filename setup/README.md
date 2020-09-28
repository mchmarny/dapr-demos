# Dapr Cluster Setup

An opinionated deployment of Dapr on Kubernetes, configured with:

* Ingress with custom domain and TLS termination
  * [ngnx](https://nginx.org/en/) for ingress controller and TLS to service mapping 
  * [letsencrypt](https://letsencrypt.org/) as certificate provider
* KEDA autoscaling
* Metrics Monitoring
  * [Prometheus](https://prometheus.io/) for metrics aggregation
  * [Grafana](https://grafana.com/) for metrics visualization with Dapr monitoring dashboards
* Log Management
  * [Fluentd](https://www.fluentd.org/) for log collection and forwarding
  * [Elasticsearch](https://www.elastic.co/) for log aggregation and query execution
  * [Kibana](https://www.elastic.co/products/kibana) for full-text log query and visualization
* Distributed Tracing
  * [Jaeger](https://www.jaegertracing.io/) for capturing traces, latency and dependency trace viewing

> All demos in the [dapr-demo](../) repository are validated on this deployment
  
## Prerequisites

* 1.15+ Kubernates cluster (see [Create Cluster](#create-cluster-on-aks) section below if you don't already have one)
* CLIs locally on the machine where you will be running this setup:
  * [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to do k8s stuff (`brew install kubectl`)
  * [Helm 3](https://helm.sh/docs/intro/install/) to install Dapr and its dependencies (`brew install helm`)
  * [certbot](https://certbot.eff.org/lets-encrypt/osx-other.html) to generate wildcard cert (`brew install certbot`)
* Domain name and access to the DNS service where you can manage that domain (required for letsencrypt challenge during cert generation and the `A` record creation to pont to the ingress gateway IP for custom domain support)

## Usage

Update [Makefile](./Makefile) variables as necessary, then:

If you need a cluster (otherwise use one selected in your kubectol context)

* `make cluster` to create a cluster on AKS

If you need TLS certificates, otherwise, use your own

* `make certs` to create TLS certs using letsencrypt

To deploy and configure Dapr 

* `make dapr` to install Dapr, KEDA, and the entire observability stack
* `make config` to perform post-install configuration

> Optionally you can use `make daprupgrade` to in place upgrade Dapr to specific version

To configure external access 

* `make ingress` to configures Ngnx ingress, SSL termination, Dapr API auth
* `make dns` to configure your DNS service for custom domain support 
* `make test` to test deployment

To deploy in-cluster data services

* `make redis` to install Redis into the cluster 
* `make mongo` to install Mongo into the cluster 
* `make kafka` to install Kafka into the cluster 

And few cluster operations helpers

* `make ports` to forward observability dashboards ports
* `make namespace` to create namespace and configure service secrets 

## Accessing observability dashboards 

To get access to the Kibana, Grafana, Zipkin dashboards run:

```shell
make ports
```

This will forward the necessary ports so you can access the dashboards using: 

* kibana - http://localhost:5601
* grafana - http://localhost:8888
* zipkin - http://localhost:9411

To stop port forwarding run 

```shell
make portstop
```

## Create Cluster on AKS

If you don't already have a Kubernates cluster, you can create one in AKS by following these steps

* Update [Makefile](./Makefile) to set:
  * `CLUSTER_NAME` - Name of the cluster you want to create 
  * `NODE_COUNT` - NUmber of nodes in the cluster default pool
  * `NODE_TYPE` - VM type used for the nodes in default pool 
* Run `make cluster`

> Note, this assumes your default Azure resource group and location are already defined. If not, run

```shell
az account set --subscription <id or name>
az configure --defaults location=<preferred location> group=<preferred resource group>
```

## Cleanup

To lists previously created clusters run 

```shell
make clusterdown
```

To delete any of the previously created clusters run 

> yes, there will be a prompt to confirm before deleting

```shell
make clusterdown CLUSTER_NAME=name
```

## Help

To find the list of all the commands with their short descriptions run: 

```shell
make help
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)