package main

import (
	"context"
	"log"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/pkg/errors"
)

var (
	client  dapr.Client
	logger  = log.New(os.Stdout, "", 0)
	address = getEnvVar("ADDRESS", ":50001")
	method  = getEnvVar("METHOD", "changes")
	pubsub  = getEnvVar("PUBSUB", "events")
	topic   = getEnvVar("TOPIC", "changes")
)

func main() {
	// create client
	c, err := dapr.NewClient()
	if err != nil {
		logger.Fatalf("error creating Dapr client: %v", err)
	}
	client = c
	defer client.Close()

	// create the service
	s, err := daprd.NewService(address)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// add method handler
	err = s.AddBindingInvocationHandler(method, bindingHandler)
	if err != nil {
		log.Fatalf("error adding binding handler: %v", err)
	}

	// start the service
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func bindingHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Printf("binding - Data:%s, Meta:%v", in.Data, in.Metadata)
	if in.Data != nil || len(in.Data) > 0 {
		if err := client.PublishEvent(ctx, pubsub, topic, in.Data); err != nil {
			logger.Printf("error publishing data to topic: %s", topic)
			return nil, errors.Wrapf(err, "error publishing data to topic: %s", topic)
		}
	}
	return nil, nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
