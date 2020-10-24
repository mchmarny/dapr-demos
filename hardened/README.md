# hardened demo 

In addition to supporting the Kubernetes logical namespace isolation and Role-based access control (RBAC) authorization, and  in-transit encryption for all sidecar-to-sidecar communication using mutual TLS, Dapr also provides additional granular control-points which can be used to harden your application deployment. 

This demo will overview: 

* Component scoping (which app should be able to access a given component)
* Pub/Sub topic scoping (which app should be able to publish or subscriber to a given topic)
* Which specific secrets should an application be able to access (deny access to others)
* Application-level access control settings with customizable "trustDomain"
* Per-operation access control settings, down to verb level (e.g. only POST on /op1)
* Cross-namespace service invocation with [SPIFFE](https://spiffe.io/) identity verification 

> Note, while you can replicate this demo on any Kubernetes cluster where Dapr is deployed, this demo uses cluster ingress to demo the cross-namespace service invocation. You can setup fully configured Dapr cluster with all the dependencies [here](../setup)

## Setup 

Create a namespace. For purposes of this demo, the namespace will be called `hardened`

```shell
kubectl create ns hardened
```

Create a Redis password

> Define the `REDIS_PASS` environment variable with your secret. You can look it up using `kubectl get svc nginx-ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}'` if necessary

```shell
kubectl create secret generic redis-secret \
    --from-literal=password="${REDIS_PASS}" \
    -n hardened 
```

## Deploy

With the setup completed, deploy the demo

```shell
kubectl apply -f k8s/ -n hardened
```

## Verify 

To ensure the rest of the demo goes smoothly, check that everything was deployed correctly

### Configuration

```shell
kubectl get configurations -n hardened
```

Should include configurations for `app1`, `app2`, and `app3`

```shell
NAME          AGE
app1-config   28s
app2-config   28s
app3-config   28s
```

### Component

```shell
kubectl get components -n hardened
```

Should include `pubsub` and `state` components:

```shell
NAME     AGE
pubsub   45s
state    45s
```

### Deployment

```shell
kubectl get pods -n hardened
```

Should include both `app1`, `app2`, and `app3` pods with the status `Running` and the ready state of `2/2` indicating that the Dapr sidecar has been injected.

```shell
NAME                    READY   STATUS    RESTARTS   AGE
app1-6df587fb45-k46sz   2/2     Running   0          40s
app2-685fd94f69-5vkwl   2/2     Running   0          40s
app3-6d57778cbd-mxn2k   2/2     Running   0          40s
```

## Demo 

If you have not done so already, start by exporting the API token from the Ingress

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

Next, check that the Dapr API has been exposed 

```shell
curl -i \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     https://api.thingz.io/v1.0/healthz
```

The response should look like this 

```shell
HTTP/2 200
date: Sat, 24 Oct 2020 19:39:05 GMT
content-length: 0
strict-transport-security: max-age=15724800; includeSubDomains
```

Now invoke the `ping` method on `app1` in the `hardened` namespace over the Dapr API on the NGNX ingress

```shell
curl -i -d '{ "message": "hello" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     https://api.thingz.io/v1.0/invoke/app1.hardened/method/ping
```

## Restart 

If you update components you may have to restart the deployments

```shell
kubectl rollout restart deployment/app1 -n hardened
kubectl rollout restart deployment/app2 -n hardened
kubectl rollout restart deployment/app3 -n hardened
kubectl rollout status deployment/app1 -n hardened
kubectl rollout status deployment/app2 -n hardened
kubectl rollout status deployment/app3 -n hardened
```

## Cleanup

```shell
kubectl delete -f k8s/ -n hardened
kubectl delete secret redis-secret -n hardened
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
