# Dapr Cluster Setup

An opinionated deployment of Dapr on Kubernetes, configured with:

* Ingress with custom domain and TLS termination
  * [NGINX](https://nginx.org/en/) for ingress controller and TLS to service mapping 
  * [letsencrypt](https://letsencrypt.org/) as certificate provider
* [KEDA](https://keda.sh/) for autoscaling
* Metrics Monitoring
  * [Prometheus](https://prometheus.io/) for metrics aggregation
  * [Grafana](https://grafana.com/) for metrics visualization with Dapr monitoring dashboards
* Log Management
  * [Fluentd](https://www.fluentd.org/) for log collection and forwarding
  * [Elasticsearch](https://www.elastic.co/) for log aggregation and query execution
  * [Kibana](https://www.elastic.co/products/kibana) for full-text log query and visualization
* Distributed Tracing
  * [Jaeger](https://www.jaegertracing.io/) for capturing traces, latency and dependency viewing

> All demos in the [dapr-demo](../) repository are validated on this deployment
  
## Prerequisites

* 1.15+ Kubernates cluster. If needed, you can setup cluster on:
  * [AKS](./aks/)
  * [GKE](./gke/)
  * AKS (coming)
* Tooling on the machine where you will be running this setup:
  * [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to do k8s stuff (`brew install kubectl`)
  * [Helm 3](https://helm.sh/docs/intro/install/) to install Dapr and its dependencies (`brew install helm`)
  * [certbot](https://certbot.eff.org/lets-encrypt/osx-other.html) to generate wildcard cert (`brew install certbot`)
* Domain name and access to the DNS service where you can manage that domain (required for letsencrypt challenge during cert generation and the `A` record creation to pont to the ingress gateway IP for custom domain support)

## Setup 

The following parameters can be used to configure your deployment. Define these as environment variables to set or override the default value:

```shell
DOMAIN            # default: example.com
DAPR_HA           # default: true
DAPR_LOG_AS_JSON  # default: true
```

> Note, make sure the correct "target" cluster is set kubectl context (`kubectl config current-context`). You can lists all registered contexts using: `kubectl config get-contexts`, and if needed, set it using `kubectl config use-context demo`.

## Usage

Start by navigate to the [setup](./setup) directory

> Run `make` by itself to see the active configuration 

To deploy and configure Dapr 

* `make dapr` to install Dapr, KEDA, and the entire observability stack
* `make config` to perform post-install configurations

> Optionally you can use `make upgrade` to in place upgrade Dapr to specific version

To configure external access 

* `make ip` (optional) to create static IP in the cluster resource group
* `make certs` to create TLS certs using letsencrypt
* `make ingress` to configures NGINX ingress, SSL termination, Dapr API auth
* `make dns` to configure your DNS service for custom domain support 
* `make test` to test deployment

To deploy in-cluster data services

* `make redis` to install Redis into the cluster 
* `make mongo` to install Mongo into the cluster 
* `make kafka` to install Kafka into the cluster 

And few cluster operations helpers

* `node-list` to print node resources and their usage
* `make ports` to forward observability dashboards ports 
* `make pass` to print the Grafana password (username: admin)
* `make nodes` to print node resource usage

Then for each namespace you want to deploy Dapr apps to

* `make namespace` to create/configure namespace with service secrets

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


## Help

To find the list of all the commands with their short descriptions run: 

```shell
make help
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)