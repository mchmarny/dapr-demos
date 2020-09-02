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
kubectl rollout status deployment/cron-binding-demo
```

If you have changed an existing component, make sure to reload the deployment and wait until the new version is ready

```shell
kubectl rollout restart deployment/cron-binding-demo
kubectl rollout status deployment/cron-binding-demo
```

Follow logs to view schedule firing 

```shell
kubectl logs -l app=cron-binding-demo -c daprd -f
```

Depending on the frequency you used there may not be an entry right away but you should see something similar to this

```json
{
  "app_id":"cron-binding-demo",
  "instance":"cron-binding-demo-6c88dbb467-j54br",
  "level":"debug",
  "msg":"next run: 59m59.629538771s",
  "scope":"dapr.contrib",
  "time":"2020-08-31T13:08:34.37049343Z",
  "type":"log",
  "ver":"0.10.0"
}
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)



