# grpc-service

For more information about service invocation see the [Dapr docs](https://github.com/dapr/docs/tree/master/concepts/service-invocation)

> You can replicate this demo on any Kubernetes cluster configured with Dapr. To demo the cross-namespace service invocation with external API gateway you will need "dapr'ized' cluster ingress (ingress with Dapr sidecar). You can setup fully configured Dapr cluster with all these dependencies using included [Dapr cluster setup](../setup#dapr-cluster-setup).

## Run 

To run this demo in Dapr, run:

```shell
dapr run \
    --app-id echo \
    --app-port 50001 \
    --app-protocol grpc \
    --dapr-http-port 3500 \
    --components-path ./config \
    go run main.go
```

## Deploy

To deploy this demo, first setup the `echo` namespace:

```shell
kubectl apply -f deployment/space.yaml
```

Then deploy and wait for the `echo-service` app pod to be ready:

```shell
kubectl apply -f deployment/app.yaml
kubectl rollout status deployment/echo-service -n echo 
```

If you have changed an existing component, make sure to reload the ingress and wait until the new version is ready

```shell
kubectl rollout restart deployment/echo-service -n echo 
kubectl rollout status deployment/echo-service -n echo 
```

Follow logs

```shell
kubectl logs -l app=echo-service -c service -f -n echo
```

In a separate terminal session export API token

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -n nginx -o jsonpath="{.data.token}" | base64 --decode)
```

And invoke the service

```shell
curl -d '{ "message": "ping" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.demo.dapr.team/v1.0/invoke/echo-service.echo/method/echo"
```

> Notice the use of `echo-service.echo` namespace in the service invocation 

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
