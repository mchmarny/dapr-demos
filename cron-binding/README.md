# cron-binding

> To help you start with this demo as a template for your new project I created [dapr-http-cron-handler-template](https://github.com/mchmarny/dapr-http-cron-handler-template)

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
dapr run --app-id cron-demo \
	    --protocol http \
	    --app-port 8080 \
	    --components-path ./config \
	    go run *.go
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](../LICENSE)



