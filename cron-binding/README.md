# cron-binding

## Binding

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: run
spec:
  type: bindings.cron
  metadata:
  - name: schedule
    value: "@every 3s"
```

For more information about this binding see the [Dapr docs](https://github.com/dapr/docs/blob/master/reference/specs/bindings/cron.md)

## Run 

Dapr cron binding demo in `go`. To use run it, first start the service

```shell
dapr run --app-id cron-binding-demo \
	    --protocol http \
	    --app-port 8080 \
	    --components-path ./config \
	    go run *.go
```


## Deploy

Deploy and wait for the pod to be ready 

```shell
kubectl apply -f k8s/component.yaml
kubectl apply -f k8s/deployment.yaml
```

If you have changed an existing component, make sure to reload the deployment and wait until the new version is ready

```shell
kubectl rollout restart deployment/cron-binding-demo
kubectl rollout status deployment/cron-binding-demo
```

Follow logs to view schedule firing 

```shell
kubectl logs -l demo=cron -c daprd -f
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](../LICENSE)



