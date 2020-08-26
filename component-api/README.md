# Use of Dapr as a component API Server

This demo users Dapr instance with API token authentication to show the use of Dapr as a API server for any of its 70+ components. To illustrate, this demo will show two use-cases:

* Sending email using Sendgrid output binding
* Querying tweets using Twitter bi-directional binding

## Setup 

### Email Component 

Create a `email-secret`

```shell
kubectl create secret generic email-secret --from-literal=apiKey=""
```

Deploy component and ensure the gateway instances are aware of it

```shell
kubectl apply -f config/email.yaml
```

### Twitter Component

Create a `twitter-secret`

```shell
kubectl create secret generic twitter-secret \
  --from-literal=consumerKey="" \
  --from-literal=consumerSecret="" \
  --from-literal=accessToken="" \
  --from-literal=accessSecret=""
```

Deploy component and ensure the gateway instances are aware of it

```shell
kubectl apply -f config/twitter.yaml
```

### Ingress Gateway

Ensure all the gateway instances are aware of these new components

```shell
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

## Usage

To use any of the components you will need the Dapr API toke: 

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

### Email 

To send email, first edit the [sample email](./sample/email.json) file: 

```json
{
    "operation": "create",
    "metadata": {
        "emailTo": "daprdemo@chmarny.com",
        "subject": "Dapr Demo"
    },
    "data": "<h1>Greetings</h1><p>Hi</p>"
}
```

And POST it to the Dapr API:

```shell
curl -v -d @./sample/email.json \
     -H "Content-Type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/bindings/send-email"
```

### Twitter 

To query the last 100 tweets for particular query, first edit the [sample query](./sample/twitter.json) file:

```json
{
    "operation": "get",
    "metadata": {
        "query": "dapr AND serverless",
        "lang": "en",
        "result": "recent"        
    }
}
```

Metadata parameters:

* `query` - can be any valid Twitter query (supports `AND`, `OR` `BUT NOT`, `FROM`, `TO`, `#`, `@`...)
* `lang` - (optional) is the [ISO 639-1](https://meta.wikimedia.org/wiki/Template:List_of_language_names_ordered_by_code) language code
* `result` - (optional) is one of:
  * `mixed` - include both popular and real time results in the response
  * `recent` - return only the most recent results in the response
  * `popular` - return only the most popular results in the response
* `since_id` - (optional) the not inclusive tweet ID query should start from 

And POST it to the Dapr API:

```shell
curl -v -d @./sample/twitter.json \
     -H "Content-Type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/bindings/query-twitter"
```

And if you have the command-line JSON processor [jq](https://shapeshed.com/jq-json/),  you can format the API results. For example, this will display only the ID, Author, and Text of each tweet:

```
curl -v -d @./sample/twitter.json \
     -H "Content-Type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/bindings/query-twitter" \
     | jq ".[] | .id_str, .user.screen_name, .text"
```

The result

```shell
"1296550502633627648"
"pacodelacruz"
"RT @daprdev: 📣Announcing the release of Dapr v0.10.0!🎉\nTons of new goodies across the board in one of our most packed releases to date!"
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](../LICENSE)
