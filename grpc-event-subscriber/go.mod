module github.com/mchmarny/dapr-grpc-event-subscriber-template

go 1.14

replace github.com/dapr/go-sdk => github.com/mchmarny/go-sdk v0.8.17

require (
	github.com/dapr/go-sdk v0.8.0
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	google.golang.org/genproto v0.0.0-20200702021140-07506425bd67 // indirect
	google.golang.org/grpc v1.30.0 // indirect
)
