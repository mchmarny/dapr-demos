# dapr-demos

Collection of personal Dapr demo

* [setup](./setup) - Deploys and configures Dapr to Kubernates 
* Bindings
  * [cron](./cron-binding) - Using scheduler to execute service 
  * [twitter](./twitter-binding) - Subscribing to a Twitter even stream
  * [state change handler](./state-change-handler) - RethinkDB state changes streamed into topic
  * twitter - Using Twitter binding 
* Eventing (Subscribing to topic and processing its events)
  * [gRPC event subscriber](./grpc-event-subscriber)
  * [HTTP event subscriber](./http-event-subscriber)
* Service 
  * [gRPC service](./grpc-service) - gRPC service example
* Integration 
  * [order cancellation](./order-cancellation) - multiple Dapr service integrations with observability
  * [Dapr API in ACI](./dapr-aci) - Dapr components as microservices 
  * [Dapr Pipeline](./dapr-pipeline) - 3 demos using Twitter binding to show incremental solution build