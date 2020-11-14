# WIP: Dapr git-ops demo 

> this demo is still being developed, don't use it!

> cluster setup 

## Setup 

### Deploy

To setup the demo, first create the namespace: 

```shell
kubectl apply -f k8s/ns.yaml
```

Than applying the rest:

```shell
kubectl apply -f k8s/
```

Check on the status: 

```shell
kubectl get pods -n gitops
```

The response should include the `gitops` pod in status `Running` with container ready state `2/2`:

```shell
NAME                      READY   STATUS    RESTARTS   AGE
gitops-5fb4d4d6f9-6m74l   2/2     Running   0          25s
```

Also, check on the ingress: 

```shell
kubectl get ingress -n gitops
```

Should include `gitops` host as well as the cluster IP mapped in your DNS:

```shell
NAME                   HOSTS              ADDRESS    PORTS   AGE
gitops-ingress-rules   gitops.thingz.io   x.x.x.x    80      19s
```

If everything went well, you should be able to navigate now to: 

http://gitops.thingz.io

## Demo 

### Edit it

Start by editing the `staticMessage` variable in [app/main.go](app/main.go) to simulate developer making code changes:

> Make sure to save your changes

```go
const (
	staticMessage = "hello PDX"
)
```

Then increment the version number variable (`APP_VERSION`) in the [app/Makefile](app/Makefile):

```shell
APP_VERSION ?=v0.1.4
```

### Tag it

When ready to make a release, tag it and push the tag to GitHub:

```shell
make tag
```

This will `git tag` it and `git push origin` your version tag to trigger to pipeline

### View it

Navigate to the cluster where the app is deployed to get the current release:

https://gitops.thingz.io/

You can also monitor the GitOps pipeline to see when you are ready to refresh the app in the browser:


## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
