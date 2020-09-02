# grpc-service

For more information about service invocation see the [Dapr docs](https://github.com/dapr/docs/tree/master/concepts/service-invocation)

## Run 

To run this demo in Dapr, run:

```shell
API_TOKEN="your-azure-cognitive-service-token" dapr run \
    --app-id sentiment-scorer \
    --app-port 60001 \
    --app-protocol grpc \
    --dapr-http-port 3500 \
    --components-path ./config \
    go run main.go
```

## Deploy

Create a `sentiment-secret`

```shell
kubectl create secret generic sentiment-secret --from-literal=token="your-azure-cognitive-service-token"
```

Deploy and wait for the pod to be ready 

```shell
kubectl apply -f deployment.yaml
kubectl rollout restart deployment/sentiment-scorer
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

Follow logs

```shell
kubectl logs -l demo=sentiment -c service -f
```

In a separate terminal session export API token

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

And invoke the service

```shell
curl -d '{ "text": "dapr is the best" }' \
    -H "Content-type: application/json" \
    -H "dapr-api-token: ${API_TOKEN}" \
    "https://api.cloudylabs.dev/v1.0/invoke/sentiment-scorer/method/sentiment"
```

Response should look something like this 

```json 
{ "sentiment":"positive", "confidence":1 }
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](./LICENSE)
