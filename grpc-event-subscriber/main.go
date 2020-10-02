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
	logger         = log.New(os.Stdout, "", 0)
	serviceAddress = getEnvVar("ADDRESS", ":60002")
	pubSubName     = getEnvVar("PUBSUB_NAME", "grpc-events")
	topicName      = getEnvVar("TOPIC_NAME", "messages")
)

func main() {
	// create Dapr service
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		logger.Fatalf("failed to start the server: %v", err)
	}

	// add handler to the service
	subscription := &common.Subscription{
		PubsubName: pubSubName,
		Topic:      topicName,
	}
	s.AddTopicEventHandler(subscription, eventHandler)

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	logger.Printf(
		"event - PubsubName:%s, Topic:%s, ID:%s, Data: %s",
		e.PubsubName, e.Topic, e.ID, e.Data,
	)

	// TODO: do something with the cloud event data

	return false, nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
