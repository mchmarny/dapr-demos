package main

import (
	"context"
	"log"
	"os"
	"strings"

	daprd "github.com/dapr/go-sdk/service/grpc"
)

var (
	logger         = log.New(os.Stdout, "", 0)
	serviceAddress = getEnvVar("ADDRESS", ":50001")
	topicName      = getEnvVar("TOPIC_NAME", "events")
)

func main() {
	// create Dapr service
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// add handler to the service
	s.AddTopicEventHandler(topicName, eventHandler)

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func eventHandler(ctx context.Context, e *daprd.TopicEvent) error {
	logger.Printf(
		"event - Source: %s, Topic:%s, ID:%s, Content Type:%s, Data:%v",
		e.Source, e.Topic, e.ID, e.DataContentType, e.Data,
	)

	// TODO: do something with the cloud event data

	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
