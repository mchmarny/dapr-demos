package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

var (
	logger           = log.New(os.Stdout, "", 0)
	serviceAddress   = getEnvVar("ADDRESS", ":60030")
	sourcePubSubName = getEnvVar("SOURCE_PUBSUB_NAME", "fanout-source-pubsub")
	sourceTopicName  = getEnvVar("SOURCE_TOPIC_NAME", "test-topic")
)

func main() {
	// create Dapr service
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// add handler to the service
	sub := &common.Subscription{
		PubsubName: sourcePubSubName,
		Topic:      sourceTopicName,
	}
	s.AddTopicEventHandler(sub, eventHandler)

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func eventHandler(ctx context.Context, e *common.TopicEvent) error {
	logger.Printf("Event - PubsubName:%s, Topic:%s, ID:%s", e.PubsubName, e.Topic, e.ID)
	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
