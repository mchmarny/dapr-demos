# twitter-binding

Example of how to use Twitter search input binding.

## Prerequisites

### dapr

To run this demo locally, you will have to have install [dapr](https://github.com). The instructions for how to do that can be found [here](https://github.com/dapr/docs/blob/master/getting-started/environment-setup.md).

### Twitter

To configure the dapr input component to query Twitter API you will also need the consumer key and secret. You can get these by registering a Twitter application [here](https://developer.twitter.com/en/apps/create).

## Setup

Assuming you have all the prerequisites mentioned above you can demo this dapr pipeline in following steps. First, insert your Twitter API secrets into the [components/twitter.yaml](components/twitter.yaml) file.

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

The `query` is the twitter search query for which you want to receive tweets.


## Run

Once the Twitter API consumer and access details are set, you are ready to run:

```shell
dapr run go run handler.go main.go \
         --app-id "consumer" \
         --app-port 8080 \
         --protocol http \
         --components-path ./config \
         --port 3500
```

Assuming everything went OK, you should see something like this:

```shell
ℹ️  Updating metadata for app command: handler.go main.go
✅  You're up and running! Both Dapr and your app logs will appear here.
```

Hope you found this demo helpful. 

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [MIT](../LICENSE)




