# grpc-service

> To help you start with this demo as a template for your new project I created [dapr-grpc-service-template](https://github.com/mchmarny/dapr-grpc-service-template)

Dapr service handler demo in `go`. To use run it, first start the service

```shell
dapr run --app-id my-service \
	    --app-port 50001 \
	    --protocol grpc \
	    --port 3500 \
         go run main.go
```

To invoke that service

```shell
curl -d '{ "message": "ping" }' \
     -H "Content-type: application/json" \
     "http://localhost:3500/v1.0/invoke/my-service/method/echo"
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](../LICENSE)
