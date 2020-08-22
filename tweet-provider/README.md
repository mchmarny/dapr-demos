# tweet-provider

Example of how to use Twitter search input binding publishing to a pub/sub topic. 

## Prerequisites

### dapr

To run this demo locally, you will have to have install [dapr](https://github.com). The instructions for how to do that can be found [here](https://github.com/dapr/docs/blob/master/getting-started/environment-setup.md).

### Twitter

To configure the dapr input component to query Twitter API you will also need the consumer key and secret. You can get these by registering a Twitter application [here](https://developer.twitter.com/en/apps/create).

## Setup

Assuming you have all the prerequisites mentioned above you can demo this dapr pipeline in following steps. First, insert your Twitter API secrets into the [config/twitter.yaml](config/twitter.yaml) file.

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: tweets
spec:
  type: bindings.twitter
  metadata:
  - name: consumerKey
    value: ""
  - name: consumerSecret
    value: ""
  - name: accessToken
    value: ""
  - name: accessSecret
    value: ""
  - name: query
    value: "serverless"
```

> The `query` is the twitter search query for which you want to receive tweets

Then make sure the is a valid pubsub component in [config/pubsub.yaml](config/pubsub.yaml). You shouldn't have to make any change to use to Redis PubSub installed with Dapr.

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: tweeter-pubsub
spec:
  type: pubsub.redis
  metadata:
  - name: redisHost
    value: localhost:6379
  - name: redisPassword
    value: ""
```


## Run

Once the Twitter API consumer and access details are set, you are ready to run:

```shell
dapr run \
    --app-id tweet-provider \
    --app-port 8080 \
    --app-protocol http \
    --dapr-http-port 3500 \
    --components-path ./config \
    go run main.go
```

## Deploy

Create secrets 

```shell
kubectl create secret generic twitter-secret \
  --from-literal=consumerKey=$TW_CONSUMER_KEY \
  --from-literal=consumerSecret=$TW_CONSUMER_SECRET \
  --from-literal=accessToken=$TW_ACCESS_TOEKN \
  --from-literal=accessSecret=$TW_ACCESS_SECRET
```

Deploy the `tweet-provider` service and its components

```shell
kubectl apply -f k8s/pubsub.yaml
kubectl apply -f k8s/twitter.yaml
kubectl apply -f k8s/deployment.yaml
```

If you have changed an existing component, make sure to reload the deployment and wait until the new version is ready

```shell
kubectl rollout restart deployment/tweet-provider
kubectl rollout status deployment/tweet-provider
```

Follow logs to see tweets for the "term" used in query

```shell
kubectl logs -l demo=twitter -c service -f
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [MIT](../LICENSE)




