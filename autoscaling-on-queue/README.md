# Autoscaling Dapr service based on queue depth 

The autoscaling demo requires Keda which runs in Kubernates. To deploy demo, first apply the `Kafka` and `Keda` components and deployment:

> Note, if you didn't use the included [setup](../setup) to configure your Kubernates cluster you may have to make changes in both components. Otherwise the defaults are fine. 

## Keda 

To install [Keda](https://github.com/kedacore/keda) into the cluster 

```shell
helm repo add keda https://kedacore.github.io/charts
helm repo update
kubectl create namespace keda
helm install keda kedacore/keda -n keda --set logLevel=debug
```

## Kafka 

To deploy in cluster version of Kafka

```shell
helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
helm repo update
kubectl create ns data
helm install kafka incubator/kafka -n data \
    --set persistence.size=16Gi \
    --set zookeeper.storage=8Gi
```

## Subscriber (processing service)

The `subscriber` doesn't really do anything so to resemble real-life processing it allows for explicit processing time setting. The default value is `300ms` but you can override it with an env vars:

```yaml
- name: PROCESS_DURATION
  value: "300ms"
```

To deploy the `subscriber` service:

```shell
kubectl apply -f subscriber/k8s/pubsub.yaml
kubectl apply -f subscriber/k8s/deployment.yaml
kubectl apply -f subscriber/k8s/keda.yaml
```

If you have changed an existing Kafka component, make sure to reload the deployment and wait until the new version is ready

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

First create the `messages` topic 

```shell
make kafka-topic
```

Then deploy the load generator 

```shell
kubectl apply -n data -f producer/config/producer.yaml
kubectl rollout -n data status deployment/queue-outoscaling-producer
kubectl logs -n data -l demo=autoscaling-producer -f
```

To increase or decrease the number of messages the producer publishes to the topic you can edit the following env vars on the producer deployment: 

`NUMBER_OF_THREADS` dictates how many concurrent publishers the producer will run:

```yaml
- name: NUMBER_OF_THREADS
  value: "1"
```

And the `THREAD_PUB_FREQ` dictates the publish frequency per each thread

```yaml
- name: THREAD_PUB_FREQ
  value: "10ms"
```

> An average, targeting in cluster Kafka broker, with 4 threads and 1ms frequency you can post about 1000 messages per second

If you need to further increase the number of posted messages, simply increase the number of `producer` replicas 

```shell
kubectl scale -n data deployment/queue-outoscaling-producer --replicas=10 
```

> Careful, you can overflow your Kafka. The topic created in this demo does have a short TTL (10 min) but if you generate a lot of volume it still will overflow it.


## Demo 

Watch Keda scaling operator log 

```shell
kubectl logs -l app=keda-operator -n keda -f
```

Follow subscriber logs 

```shell
kubectl logs -l demo=autoscaling-demo -c service -f 
```

> WIP: replicas of `queue-outoscaling-subscriber` do not scale up based on the depth of the queue. With `lagThreshold` of `3` and Keda seeing queue lag of `7K+`, there number of replicas remains 1. Interestingly, with `minReplicaCount` is set to `0`, the pod will be scale to 0 when the producer is not sending messages, despite the fact that the queue still has thousands of unack'd messages 

```json
{"level":"debug","ts":1598713533.3401804,"logger":"scalehandler","msg":"Scaler for scaledObject is active","ScaledObject.Namespace":"default","ScaledObject.Name":"prime-calculator-scaler","ScaledObject.ScaleType":"deployment","Scaler":{}}
{"level":"debug","ts":1598713533.3656452,"logger":"scalehandler","msg":"ScaledObject's Status was properly updated","ScaledObject.Namespace":"default","ScaledObject.Name":"prime-calculator-scaler","ScaledObject.ScaleType":"deployment"}
{"level":"debug","ts":1598713536.3939075,"logger":"kafka_scaler","msg":"Group autoscaling has a lag of 70867 for topic messages and partition 0\n"}
{"level":"debug","ts":1598713536.3939333,"logger":"scalehandler","msg":"Scaler for scaledObject is active","ScaledObject.Namespace":"default","ScaledObject.Name":"prime-calculator-scaler","ScaledObject.ScaleType":"deployment","Scaler":{}}
{"level":"debug","ts":1598713536.403159,"logger":"scalehandler","msg":"ScaledObject's Status was properly updated","ScaledObject.Namespace":"default","ScaledObject.Name":"prime-calculator-scaler","ScaledObject.ScaleType":"deployment"}
{"level":"debug","ts":1598713539.4521942,"logger":"kafka_scaler","msg":"Group autoscaling has a lag of 70864 for topic messages and partition 0\n"}
{"level":"debug","ts":1598713539.452236,"logger":"scalehandler","msg":"Scaler for scaledObject is active","ScaledObject.Namespace":"default","ScaledObject.Name":"prime-calculator-scaler","ScaledObject.ScaleType":"deployment","Scaler":{}}
{"level":"debug","ts":1598713539.4604158,"logger":"scalehandler","msg":"ScaledObject's Status was properly updated","ScaledObject.Namespace":"default","ScaledObject.Name":"prime-calculator-scaler","ScaledObject.ScaleType":"deployment"}
````

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](./LICENSE)
