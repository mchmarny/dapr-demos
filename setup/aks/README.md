# AKS Cluster Setup

The following parameters can be used to configure your deployment. Define these as environment variables to set or override the default value:

```shell
CLUSTER_NAME      # default: demo
CLUSTER_VERSION   # default: 1.18.8
NODE_COUNT        # default: 3
NODE_TYPE         # default: Standard_D4_v2
```

> Note, this assumes your default Azure resource group and location are already defined. If not, run

```shell
az account set --subscription <id or name>
az configure --defaults location=<preferred location> group=<preferred resource group>
```

## Usage

Start by navigating to the [setup/aks](./setup/aks) directory

> Run `make` by itself to see the active configuration 

* `make cluster` to create a cluster on AKS (make cluster CLUSTER_NAME=demo)
* `make ip` (optional) to create static IP in the cluster resource group
* `make node-pool` (optional) to add new AKS node pool
* `make node-list` to print node resource usage
* `make cluster-list` to list your AKS clusters
* `make version-list` to list Kubernetes versions supported on AKS

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