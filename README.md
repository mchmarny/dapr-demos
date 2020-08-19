# Dapr demos

Collection of personal Dapr demo

* [setup](./setup) - Deploys and configures Dapr in Kubernates 
* Bindings
  * [cron](./cron-binding) - Using scheduler to execute service 
  * [twitter](./twitter-binding) - Subscribing to a Twitter even stream
  * [state change handler](./state-change-handler) - RethinkDB state changes streamed into topic
* Eventing (Subscribing to topic and processing its events)
  * [gRPC event subscriber](./grpc-event-subscriber)
  * [HTTP event subscriber](./http-event-subscriber)
* Service 
  * [gRPC service](./grpc-service) - gRPC service example
* Integrations
  * [order cancellation](./order-cancellation) - multiple Dapr service integrations with observability
  * [Dapr API in ACI](./dapr-aci) - Dapr components as microservices 
  * [Dapr Pipeline](./dapr-pipeline) - 3 demos using Twitter binding to show incremental solution build


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](./LICENSE)
