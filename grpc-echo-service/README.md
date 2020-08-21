# grpc-service

For more information about service invocation see the [Dapr docs](https://github.com/dapr/docs/tree/master/concepts/service-invocation)

## Run 

To run this demo in Dapr, run:

```shell
dapr run \
    --app-id grpc-service-demo \
    --app-port 50001 \
    --app-protocol grpc \
    --dapr-http-port 3500 \
    --components-path ./config \
    go run main.go
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](./LICENSE)
