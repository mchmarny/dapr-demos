# GKE Cluster Setup

The following parameters can be used to configure your deployment. Define these as environment variables to set or override the default value:

```shell
CLUSTER_NAME     # default: demo
CLUSTER_ZONE     # default: us-west1-a
CLUSTER_VERSION  # default: 1.17.13-gke.600
NODE_COUNT       # default: 2
NODE_COUNT_MIN   # default: 1
NODE_COUNT_MAX   # default: 5
NODE_TYPE        # default: n2-standard-4
```

```shell
gcloud config set project <your project ID>
gcloud config set compute/region <your preferred region>
gcloud config set compute/zone <your preferred zone>
```

## Usage

Start by navigating to the [setup/gke](./setup/gke) directory

> Run `make` by itself to see the active configuration 

* `make cluster` to create a cluster on GKE (make cluster CLUSTER_NAME=demo)
* `make cluster-list` to list your AKS clusters
* `make version-list` to list Kubernetes versions supported on AKS

Once cluster is created, you can follow [these instructions](../) to configure Dapr.

## Cleanup

To lists previously created clusters run 

```shell
make cluster-list
```

To delete any of the previously created clusters run 

> yes, there will be a prompt to confirm before deleting

```shell
make cluster-down CLUSTER_NAME=name
```

## Help

To find the list of all the commands with their short descriptions run: 

```shell
make help
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../../LICENSE)