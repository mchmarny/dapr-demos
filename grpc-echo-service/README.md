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

## Deploy

DEploy and wait for the pod to be ready 

```shell
kubectl apply -f deployment.yaml
watch kubectl get pods
```

Follow logs

```shell
kubectl logs -l demo=echo -c service -f
```

In a separate terminal session export API token

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

And invoke the service

```shell
curl -d '{ "message": "ping" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/invoke/grpc-echo-service/method/echo"
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](./LICENSE)
