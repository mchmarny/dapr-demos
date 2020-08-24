# Use of Dapr as a component API Server

This demo users Dapr instance with API token authentication to show the use of Dapr as a API server for any of its 70+ components. This demo uses Sendgrid output binding.

## Setup 

Create a `email-secret`

```shell
kubectl create secret generic email-secret --from-literal=apiKey="your-sandgrid-api-key"
```

Deploy component and ensure the gateway instances are aware of it

```shell
kubectl apply -f binding.yaml
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

## Use

First, export API token

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

Now POST an email to the Dapr API following message using `curl`.


```shell
curl -v -d @./email.json \
     -H "Content-Type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/bindings/send-email"
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](../LICENSE)
