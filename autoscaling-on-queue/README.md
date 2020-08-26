# Autoscaling Dapr service based on queue depth 

The autoscaling demo requires Keda which runs in Kubernates. To deploy demo, first apply the `Kafka` and `Keda` components and deployment:

> Note, if you didn't use the included [setup](../setup) to configure your Kubernates cluster you may have to make changes in both components. Otherwise the defaults are fine. 

```shell
kubectl apply -f subscriber/k8s/pubsub.yaml
kubectl apply -f subscriber/k8s/state.yaml
kubectl apply -f subscriber/k8s/deployment.yaml
kubectl apply -f subscriber/k8s/keda.yaml
```

If you have changed an existing component, make sure to reload the deployment and wait until the new version is ready

```shell
kubectl rollout restart deployment/autoscaling-on-queue
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

Check Dapr logs to make sure there is no errors 

```shell
kubectl logs -l demo=autoscaling-demo -c daprd
```

Now export API token:

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

And publish a single request to make sure everything is up:

> To give the autoscaling demo something to scale on, the scaled service will calculate highest prime number up to the provided `max` number. 

```shell
curl -v -d '{"id":"id1","max":100,"Time":1598480443}' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/publish/autoscaling-kafka-queue/primes"
```

And check the service logs:

```shell
kubectl logs -l demo=autoscaling-demo -c service -f
```

If everything goes well you should see 

```shell
Request - PubSub:autoscaling-kafka-queue, Topic:primes, ID:39465d38-0a6e-4ef0-9011-6c978b18eb15
Highest prime for 100 is 97
Previous high: 97, New: 97
```

## Autoscaling Demo  

> TODO: Run some load on the Dapr binding 


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](./LICENSE)
