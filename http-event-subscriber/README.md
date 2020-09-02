# dapr-event-subscriber-template

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
dapr run \
    --app-id grpc-event-subscriber-demo \
    --app-port 50001 \
    --app-protocol http \
    --dapr-http-port 3500 \
    --components-path ./config \
    go run main.go
```


## Deploy

Deploy and wait for the pod to be ready 

```shell
kubectl apply -f k8s/component.yaml
kubectl apply -f k8s/deployment.yaml
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

If you have changed an existing component, make sure to reload the deployment and wait until the new version is ready

```shell
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

Follow logs

```shell
kubectl logs -l app=http-event-subscriber -c service -f
```

In a separate terminal session export API token

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

And invoke the service

```shell
curl -d '{ "from": "John", "to": "Lary", "message": "hi" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/publish/http-events/messages"
```

In the logs, you should see now an entry similar to this. Feel free to edit the message and try again.

```shell
event - PubsubName:http-events, Topic:messages, ID:6b6cc665-684d-456c-8880-56e20cdf0519, Data: map[from:John message:hi to:Lary]
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
