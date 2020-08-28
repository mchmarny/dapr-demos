# Autoscaling Dapr service based on queue depth 

The autoscaling demo requires Keda which runs in Kubernates. To deploy demo, first apply the `Kafka` and `Keda` components and deployment:

> Note, if you didn't use the included [setup](../setup) to configure your Kubernates cluster you may have to make changes in both components. Otherwise the defaults are fine. 

## Subscriber (processing service)

To deploy the consumer service:

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

## Producer (generating load on the Kafka topic)

First create the `primes` topic 

```shell
make kafka-topic
```

Then deploy the load generator 

```shell
kubectl apply -n data -f producer/config/producer.yaml
kubectl rollout -n data status deployment/prime-calculator-request-producer
kubectl logs -n data -l demo=autoscaling-producer -f
```

Note, this will generate about 1,000 messages per second per "thread". To increase the volume, increase the "NUMBER_OF_THREADS" variable in deployment. Depending on the setup you can probably get this to about 10K events from a single deployment. To increase the volume further: 

```shell
kubectl scale -n data deployment/prime-calculator-request-producer --replicas=10 
```

> Careful, you can overflow your Kafka deployment. The topic we created does have a short TTL (10 min) but if you generate a lot of volume it will crash. 


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](./LICENSE)
