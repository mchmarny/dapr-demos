# Dapr pipeline demo 

Dapr supports a wide array of state and pubsub building blocks across multiple OSS, Cloud, and on-prem services. This demo shows how to use a few of these components to build tweet sentiment processing pipeline.

![alt text](./img/overview.png "Pipeline Overview")


## Run in standalone mode

### Setup 

To run these demos locally, you will have first create a secret file (`pipeline/secrets.json`). These will be used by Dapr components at runtime. To get the Twitter API secretes you will need to register your app [here](https://developer.twitter.com/en/apps/create).

```json
{
    "Twitter": {
        "ConsumerKey": "",
        "ConsumerSecret": "",
        "AccessToken": "",
        "AccessSecret": ""
    },
    "Azure": {
        "CognitiveAPIKey": ""
    }
}
```


### Start tweet viewer app 

Navigate to the [tweet-viewer](./tweet-viewer) directory and run:

```shell
cd tweet-viewer
dapr run \
    --app-id tweet-viewer \
    --app-port 8084 \
    --app-protocol http \
    --components-path ./config \
    go run main.go
```

Once the app starts, you should be able to navigate to http://localhost:8084/. There won't be anything there yet, but if you see `connection: open` in the top right corner that means the WebSocket connection to the back-end is established. 


### Start sentiment scoring service 

Navigate to the [sentiment-scorer](./sentiment-scorer) directory and run:

```shell
cd sentiment-scorer
dapr run \
    --app-id sentiment-scorer \
    --app-port 60005 \
    --app-protocol grpc \
    --components-path ./config \
    go run main.go
```

The last line from the above command should be

```shell
✅  You're up and running! Both Dapr and your app logs will appear here.
```

### Start tweet processing service 

Navigate to the [tweet-processor](./tweet-processor) directory and run:

```shell
cd tweet-processor
dapr run \
    --app-id tweet-processor \
    --app-port 60002 \
    --app-protocol grpc \
    --components-path ./config \
    go run main.go
```

The last line from the above command should be

```shell
✅  You're up and running! Both Dapr and your app logs will appear here.
```


### Start tweet provider

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

### View sentiment scored tweets in the UI 

Navigate once more to http://localhost:8084/ and provided there were tweets matching your query you should now see tweets displayed in the UI. 

![](./img/ui.png)


## Run on Kubernetes 

> This section is WIP. Use the individual Makefiles located in each app to deploy to Kubernetes for now. 

### tweet-processor

Deploy `tweet-processor` and wait for it to be ready

```shell
kubectl apply -f tweet-processor/k8s/result-pubsub.yaml
kubectl apply -f tweet-processor/k8s/source-pubsub.yaml
kubectl apply -f tweet-processor/k8s/deployment.yaml
kubectl rollout status deployment/tweet-processor
```

> Note, throughout this demo, if you have changed an existing component, make sure to reload the deployment and wait until the new version is ready

```shell
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

Check the Dapr logs to make sure the components were registered correctly 

```shell
kubectl logs -l app=tweet-processor -c daprd --tail 200
```

### sentiment-scorer

Create a secret for Azure Cognitive Services

```shell
kubectl create secret generic sentiment-secret --from-literal=token="your-azure-cognitive-service-token"
```

Deploy `sentiment-scorer` and wait for it to be ready 

```shell
kubectl apply -f sentiment-scorer/deployment.yaml
kubectl rollout status deployment/sentiment-scorer
```

Check the logs to make sure Dapr was started correctly 

```shell
kubectl logs -l app=sentiment-scorer -c daprd --tail 200
```

To test the service, you can first export the API token

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

And then invoke the service manually

```shell
curl -i -d '{ "text": "dapr is the best" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/invoke/sentiment-scorer/method/sentiment"
```

Response should look something like this 

```json 
{ "sentiment":"positive", "confidence":1 }
```


### tweet-viewer

Deploy `tweet-viewer` along with its component

```shell
kubectl apply -f tweet-viewer/k8s/source-pubsub.yaml
kubectl apply -f tweet-viewer/k8s/deployment.yaml
kubectl rollout status deployment/tweet-viewer
```

Patch ingress to add the viewer rule

```shell
kubectl get ing/ingress-rules -o json \
  | jq '.spec.rules += [{"host":"tweets.cloudylabs.dev","http":{"paths":[{"backend": {"serviceName":"tweet-viewer","servicePort":80},"path":"/"}]}}]' \
  | kubectl apply -f -
```

Check that the ingress was updated 

```shell
kubectl get ingress
```

Should include `viewer.`

```shell
NAME            HOSTS                                      ADDRESS   PORTS     AGE
ingress-rules   api.cloudylabs.dev,tweets.cloudylabs.dev   x.x.x.x   80, 443   9h
```

Check in browser: https://tweets.cloudylabs.dev

### tweet-provider

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

Check Dapr to make sure components were registered correctly 

```shell
kubectl logs -l app=tweet-provider -c daprd --tail 200
```

## View

Navigate back to: https://tweets.cloudylabs.dev

If everything went well, you should see some tweets appear. 

> Note, this demo shows only tweets meeting your query posted since the viewer was started. If you chosen an unpopular search term you may have to be patient

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
