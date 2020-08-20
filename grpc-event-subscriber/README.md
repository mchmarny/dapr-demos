# grpc-event-subscriber

## Components

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: events
spec:
  type: pubsub.redis
  metadata:
  - name: redisHost
    value: localhost:6379
  - name: redisPassword
    value: ""

```

For more information about pub/sub see the [Dapr docs](https://github.com/dapr/docs/tree/master/concepts/publish-subscribe-messaging)

## Run 

To run this demo in Dapr, run:

```shell
dapr run --app-id grpc-event-subscriber-demo \
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
