# Tweet Provider Demo 

## Setup 

To run these demos locally, you will have first create a secret file (`pipeline/secrets.json`). These will be used by Dapr components at runtime. To get the Twitter API secretes you will need to register your app [here](https://developer.twitter.com/en/apps/create).

```json
{
    "Twitter": {
        "ConsumerKey": "",
        "ConsumerSecret": "",
        "AccessToken": "",
        "AccessSecret": ""
    }
}
```

## Run it

### Standalone Mode

Navigate to the [tweet-provider](./tweet-provider) directory and run:

```shell
cd tweet-provider
dapr run \
    --app-id tweet-provider \
    --app-port 8080 \
    --app-protocol http \
    --components-path ./config \
    go run main.go
```

The last line from the above command should be

```shell
✅  You're up and running! Both Dapr and your app logs will appear here.
```

Your tweets should appear in the logs now


### Kubernetes 


Create secret for `tweet-provider` to connect to Twitter API 

```shell
kubectl create secret generic twitter-secret \
  --from-literal=consumerKey="" \
  --from-literal=consumerSecret="" \
  --from-literal=accessToken="" \
  --from-literal=accessSecret=""
```

Deploy the `tweet-provider` service and its components

```shell
kubectl apply -f tweet-provider/k8s/state.yaml
kubectl apply -f tweet-provider/k8s/pubsub.yaml
kubectl apply -f tweet-provider/k8s/twitter.yaml
kubectl apply -f tweet-provider/k8s/deployment.yaml
kubectl rollout status deployment/tweet-provider
```

If you have changed an existing component, make sure to reload the deployment and wait until the new version is ready

```shell
kubectl rollout restart deployment/tweet-provider
kubectl rollout status deployment/tweet-provider
```

Check Dapr to make sure components were registered correctly 

```shell
kubectl logs -l app=tweet-provider -c daprd --tail 200
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
