# http-event-subscriber

> To help you start with this demo as a template for your new project I created [dapr-http-event-subscriber-template](https://github.com/mchmarny/dapr-event-subscriber-template)


Dapr HTTP event subscriber services demo in `go`. To use run it, first start the service

```shell
dapr run --app-id event-subscriber \
         --app-port 8080 \
         --protocol http \
         --port 3500 \
         --components-path ./config \
         go run main.go
```

Then send an event to that service 

```shell
curl -d '{ "from": "John", "to": "Lary", "message": "hi" }' \
     -H "Content-type: application/json" \
     "http://localhost:3500/v1.0/publish/events"
```

You can use the provided `makefile` to help with executing these commands 

```shell
$ make help
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](../LICENSE)
