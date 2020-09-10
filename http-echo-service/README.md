# grpc-service

For more information about service invocation see the [Dapr docs](https://github.com/dapr/docs/tree/master/concepts/service-invocation)

## Run 

To run this demo in Dapr, run:

```shell
dapr run \
    --app-id grpc-service-demo \
    --app-port 8080 \
    --app-protocol http \
    --dapr-http-port 3500 \
    --components-path ./config \
    go run main.go
```

## Deploy

Deploy and wait for the pod to be ready 

```shell
kubectl apply -f deployment.yaml
kubectl rollout status deployment/http-echo-service
```

If you have changed an existing component, make sure to reload the ingress and wait until the new version is ready

```shell
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

Follow logs

```shell
kubectl logs -l app=http-echo-service -c service -f
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
     "https://api.cloudylabs.dev/v1.0/invoke/http-echo-service/method/echo"
```

The response should include the sent message 

```json
{ 
    "message": "ping" 
}
```

And the logs

```shell
Invocation (ContentType:application/json, Verb:POST, QueryString:map[], Data:{ "message": "ping" })
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
