# hardened demo 

> WIP, this readme is still being updated. All the deployment files are there but the actual hardened Dapr app demo is not fully described yet. 

## Setup 

Create a namespace

```shell
kubectl create ns hardened
```

Create a Redis password 

> Define the `REDIS_PASS` environment variable with your secret

```shell
kubectl create secret generic redis-secret \
    --from-literal=password="${REDIS_PASS}" \
    -n hardened 
```

## Deploy

```shell
kubectl apply -f k8s/ -n hardened
```

## Verify 


### Configuration

```shell
kubectl get configurations -n hardened
```

Should include `app1-config` and `app2-config` configurations:

```shell
NAME          AGE
app1-config   28s
app2-config   28s
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

Should include both `app1` and `app2` pods with the status `Running` and the ready state of `2/2` indicating that the Dapr sidecar has been injected.

```shell
NAME                    READY   STATUS    RESTARTS   AGE
app1-568c8547f4-v5psz   2/2     Running   0          23s
app2-5d46976b99-2hmcw   2/2     Running   0          23s
```

## Demo 

Forwards Dapr port for app1

```shell
kubectl port-forward deployment/app1 3501:3500 -n hardened
```

Follow the logs for app1 

```shell
kubectl logs -l app=app1 -n hardened -c app -f
```

Invoke its `call` method

```shell
curl -i -d '{ "message": "hello" }' \
     -H "Content-type: application/json" \
     "http://localhost:3501/v1.0/invoke/app1/method/call"
```

Invoke its `pub` method

```shell
curl -i -d '{ "message": "hello" }' \
     -H "Content-type: application/json" \
     "http://localhost:3501/v1.0/invoke/app1/method/pub"
```

## Restart 

If you update components you may have to restart the deployments

```shell
kubectl rollout restart deployment/app1 -n hardened
kubectl rollout restart deployment/app2 -n hardened
kubectl rollout status deployment/app1 -n hardened
kubectl rollout status deployment/app2 -n hardened
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
