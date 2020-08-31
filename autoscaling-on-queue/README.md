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
helm repo add confluentinc https://confluentinc.github.io/cp-helm-charts/
helm repo update
kubectl create ns data
helm install kafka confluentinc/cp-helm-charts -n data \
		--set cp-schema-registry.enabled=false \
		--set cp-kafka-rest.enabled=false \
		--set cp-kafka-connect.enabled=false \
		--set dataLogDirStorageClass=default \
		--set dataDirStorageClass=default \
		--set storageClass=default
```

## Subscriber (processing service)

The `subscriber` doesn't really do anything so to resemble real-life processing it allows for explicit processing time setting. The default value is `300ms` but you can override it with an env vars:

```yaml
- name: PROCESS_DURATION
  value: "300ms"
```

To deploy the `subscriber` service:

```shell
kubectl apply -f subscriber/k8s/binding.yaml
kubectl apply -f subscriber/k8s/deployment.yaml
kubectl apply -f subscriber/k8s/keda.yaml
```

If you have changed an existing Kafka component, make sure to reload the deployment and wait until the new version is ready

```shell
kubectl rollout restart deployment/queue-outoscaling-subscriber
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```

## Producer (generating load on the Kafka topic)

First, deploy a Kafka client:

```shell
kubectl apply -n data -f producer/config/kafka-client.yaml
kubectl wait -n data --for=condition=ready pod kafka-client --timeout=120s
```

Then, create the `messages` topic: 

```shell
kubectl -n data exec -it kafka-client -- kafka-topics \
		--zookeeper kafka-cp-zookeeper-headless:2181 \
		--topic messages \
		--create \
		--partitions 5 \
		--replication-factor 1 \
		--if-not-exists
```

Then deploy the load `producer`:

```shell
kubectl apply -n data -f producer/config/producer.yaml
kubectl rollout -n data status deployment/queue-outoscaling-producer
kubectl logs -n data -l demo=autoscaling-producer -f
```

To stop the `producer`

```shell
kubectl delete -n data -f producer/config/producer.yaml
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

Watch the number of `subscriber` pods being adjusted based on the depth of the queue:

```shell
watch kubectl get pods -l app=queue-outoscaling-subscriber
```

```shell
Every 2.0s: kubectl get pods -l app=queue-outoscaling-subscriber

NAME                                            READY   STATUS    RESTARTS   AGE
queue-outoscaling-subscriber-674c7dc7b4-cjp48   2/2     Running   0          14m
queue-outoscaling-subscriber-674c7dc7b4-fkg7m   2/2     Running   0          14m
queue-outoscaling-subscriber-674c7dc7b4-sdj9z   2/2     Running   0          14m
```

Watch Keda scaling operator log for the depth of queue signal:

```shell
kubectl logs -l app=keda-operator -n keda -f
```

```json
{"level":"debug","ts":1598716685.4928422,"logger":"kafka_scaler","msg":"Group autoscaling has a lag of 2 for topic messages and partition 0\n"}
{"level":"debug","ts":1598716685.4929283,"logger":"scalehandler","msg":"Scaler for scaledObject is active","ScaledObject.Namespace":"default","ScaledObject.Name":"queue-outoscaling-scaler","ScaledObject.ScaleType":"deployment","Scaler":{}}
{"level":"debug","ts":1598716685.5025718,"logger":"scalehandler","msg":"ScaledObject's Status was properly updated","ScaledObject.Namespace":"default","ScaledObject.Name":"queue-outoscaling-scaler","ScaledObject.ScaleType":"deployment"}
```

Follow subscriber logs for the processing throughput:

```shell
kubectl logs -l demo=autoscaling-demo -c service -f 
```

```shell
received:        746,   0 errors - avg   4/sec
received:        773,   0 errors - avg   4/sec
received:        794,   0 errors - avg   4/sec
```

You can play with the volume of data published by the `producer`, the speed of processing by the `subscriber`, and the Keda parameters to show different scenarios

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](../LICENSE)
