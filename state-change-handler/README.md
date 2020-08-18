# Dapr State Change Publisher 

Binding to detect RethinkDB state changes and stream them into a single topic. 

## Components 

### Source 

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: changes
spec:
  type: digitaltwins.rethinkdb.statechange
  metadata:
  - name: address
    value: "127.0.0.1:28015"
  - name: database
    value: "dapr"
```

### Target

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: events
spec:
  type: pubsub.redis
  metadata:
  - name: redisHost
    value: localhost:6379
  - name: redisPassword
    value: ""
```

## Service 

```shell
dapr run --app-id publisher \
	       --protocol grpc \
	       --app-port 50001 \
	       --components-path ./config \
	       go run main.go
```