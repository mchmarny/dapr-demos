# Autoscaling Dapr service based on queue depth 

The autoscaling demo requires Keda which runs in Kubernetes. To deploy demo, first apply the `Kafka` and `Keda` components to your Dapr deployment:

> Note, if you didn't use the included [setup](../setup) to configure your Kubernetes cluster you may have to make changes in both components. Otherwise the defaults are fine. 

## Keda 

To install [Keda](https://github.com/kedacore/keda) into the cluster using Help

```shell
kubectl create ns keda
kubectl apply -f https://github.com/kedacore/keda/releases/download/v2.0.0-beta/keda-2.0.0-beta.yaml
```

## Kafka 

To deploy in cluster version of Kafka

```shell
helm repo add confluentinc https://confluentinc.github.io/cp-helm-charts/
helm repo update
kubectl create ns kafka
helm install kafka confluentinc/cp-helm-charts -n kafka \
		--set cp-schema-registry.enabled=false \
		--set cp-kafka-rest.enabled=false \
		--set cp-kafka-connect.enabled=false \
		--set dataLogDirStorageClass=default \
		--set dataDirStorageClass=default \
		--set storageClass=default
kubectl rollout status deployment.apps/kafka-cp-control-center -n kafka
kubectl rollout status deployment.apps/kafka-cp-ksql-server -n kafka
kubectl rollout status statefulset.apps/kafka-cp-kafka -n kafka
kubectl rollout status statefulset.apps/kafka-cp-zookeeper -n kafka
```

When done, deploy Kafka client and wait until it's ready:

```shell
kubectl apply -n kafka -f deployment/kafka-client.yaml
kubectl wait -n kafka --for=condition=ready pod kafka-client --timeout=120s
```

When done, create the `metrics` topic: 

> These are just mocked events so no need to partition or replicate

```shell
kubectl -n kafka exec -it kafka-client -- kafka-topics \
		--zookeeper kafka-cp-zookeeper-headless:2181 \
		--topic metrics \
		--create \
		--partitions 1 \
		--replication-factor 1 \
		--if-not-exists
```

## Subscriber (processing service)

The `subscriber` doesn't really do anything so to resemble real-life processing it allows for explicit processing time setting. The default value is `500ms` but you can override it with an env vars:

```yaml
- name: PROCESS_DURATION
  value: "500ms"
```

To deploy the `subscriber` service:

```shell
kubectl apply -f deployment/kafka-pubsub.yaml
kubectl apply -f deployment/subscriber.yaml
kubectl apply -f deployment/keda-scaler.yaml
```

When done, start watching for the number of replicas of the deployed `subscriber` service 

```shell
watch kubectl get pods -l app=autoscaling-subscriber
```

To see the logs from `subscriber` service 

```shell
kubectl logs -l app=autoscaling-subscriber -c service -f
```

## Producer (generating load on the Kafka topic)

In a second terminal session now, deploy the `producer` and wait for it to be ready:

```shell
kubectl apply -f deployment/producer.yaml
kubectl rollout status deployment/autoscaling-producer
```

When done, start following the produces service logs 

```shell
kubectl logs -l app=autoscaling-producer -c service -f
```

If you need to stop the `producer`:

```shell
kubectl delete -f deployment/producer.yaml
```

## Demo 

Back in the initial terminal now, watch the number of `subscriber` pods being adjusted based on the depth of the queue:

```shell
NAME                                      READY   STATUS    RESTARTS   AGE
autoscaling-subscriber-674c7dc7b4-cjp48   2/2     Running   0          14m
autoscaling-subscriber-674c7dc7b4-fkg7m   2/2     Running   0          14m
autoscaling-subscriber-674c7dc7b4-sdj9z   2/2     Running   0          14m
```

To watch Keda scaling operator log for the depth of queue signal:

```shell
kubectl logs -l app=keda-operator -n keda -f
```

```json
{"level":"debug","ts":1598716685.4928422,"logger":"kafka_scaler","msg":"Group autoscaling has a lag of 2 for topic messages and partition 0\n"}
{"level":"debug","ts":1598716685.4929283,"logger":"scalehandler","msg":"Scaler for scaledObject is active","ScaledObject.Namespace":"default","ScaledObject.Name":"queue-outoscaling-scaler","ScaledObject.ScaleType":"deployment","Scaler":{}}
{"level":"debug","ts":1598716685.5025718,"logger":"scalehandler","msg":"ScaledObject's Status was properly updated","ScaledObject.Namespace":"default","ScaledObject.Name":"queue-outoscaling-scaler","ScaledObject.ScaleType":"deployment"}
```

You can also follow subscriber logs for the processing throughput:

```shell
kubectl logs -l app=autoscaling-subscriber -c service -f 
```

```shell
received:        746,   0 errors - avg   4/sec
received:        773,   0 errors - avg   4/sec
received:        794,   0 errors - avg   4/sec
```


If the `subscriber` is not being scaled, you may have to adjust the [deployment/keda.yaml](deployment/keda.yaml) parameters. The default minimum number of replicas is `0` and maximum is `10`. The `lagThreshold` (the number of topic messages the subscriber can be behind) is `5`.

You can also adjust the `producer` variables defined in the [deployment/producer.yaml](./deployment/producer.yaml) file to increase the number of events it publishes: 

* `NUMBER_OF_PUBLISHERS` - number of channels that are used to publish events (default: 1)
* `PUBLISHERS_FREQ` - frequency with which each channel publishes events (default: 1s) 
* `LOG_FREQ` - frequency with which the publisher prints out the processing throughout (default: 3s)

Depending on your network infrastructure (which may on some clouds be related to the size of the node VM), you can should get about 2,400 events by setting `NUMBER_OF_PUBLISHERS` to `4` and `PUBLISHERS_FREQ` to `100ms`. YMMV.

You can also scale the number of `autoscaling-producer` replicas:

```shell
kubectl scale -n kafka deployment/autoscaling-producer --replicas=10 
```

## Updating Components 

If you have changed an existing component, make sure to reload the deployments and wait until the new versions is ready

```shell
kubectl rollout restart deployment/autoscaling-subscriber
kubectl rollout status deployment/autoscaling-subscriber
kubectl rollout restart deployment/autoscaling-producer
kubectl rollout status deployment/autoscaling-producer
```

## Kafka Helpers 

Describe `metrics` topic

```shell
kubectl -n kafka exec -it kafka-client -- kafka-topics \
	--zookeeper kafka-cp-zookeeper:2181 \
	--topic metrics \
	--describe
```

Get the subscriber offsets for `metrics`

```shell
kubectl -n kafka exec -it kafka-client -- kafka-consumer-groups \
	--bootstrap-server kafka-cp-kafka:9092 \
	--describe \
	--group autoscaling-subscriber
```

Purge the `metrics` topic

```shell
kubectl -n kafka exec -it kafka-client -- kafka-topics \
	--zookeeper kafka-cp-zookeeper:2181 \
	--alter \
	--topic metrics \
	--config retention.ms=1000
sleep 15
kubectl -n kafka exec -it kafka-client -- kafka-topics \
	--zookeeper kafka-cp-zookeeper:2181 \
	--alter \
	--topic metrics \
	--delete-config retention.ms
```

Delete `metrics` topic

```shell
kubectl -n kafka exec -it kafka-client -- kafka-topics \
	--zookeeper kafka-cp-zookeeper:2181 \
	--delete \
	--topic metrics
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
