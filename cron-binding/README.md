# cron-binding

> To help you start with this demo as a template for your new project I created [dapr-http-cron-handler-template](https://github.com/mchmarny/dapr-http-cron-handler-template)

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
