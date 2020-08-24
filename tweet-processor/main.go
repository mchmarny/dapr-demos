package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/pkg/errors"
)

var (
	logger         = log.New(os.Stdout, "", 0)
	serviceAddress = getEnvVar("ADDRESS", ":60002")

	srcPubSubName = getEnvVar("SOURCE_PUBSUB_NAME", "tweeter-pubsub")
	srcTopicName  = getEnvVar("SOURCE_TOPIC_NAME", "tweets")

	resultPubSubName = getEnvVar("RESULT_PUBSUB_NAME", "processed-tweets-pubsub")
	resultTopicName  = getEnvVar("RESULT_TOPIC_NAME", "tweets")

	client dapr.Client
)

func main() {
	// create Dapr service
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	c, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("error creating Dapr client: %v", err)
	}
	client = c
	defer client.Close()

	// add handler to the service
	subscription := &common.Subscription{
		PubsubName: srcPubSubName,
		Topic:      srcTopicName,
	}
	s.AddTopicEventHandler(subscription, tweetHandler)

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func tweetHandler(ctx context.Context, e *common.TopicEvent) error {
	logger.Printf(
		"Processing Tweet (pubsub:%s/topic:%s) - ID:%s, Data: %s",
		e.PubsubName, e.Topic, e.ID, e.Data,
	)

	tweet, ok := e.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("event data not in expected type: %T", e.Data)
	}

	// TODO: invoke sentiment scoring service

	// TODO: augment tweet with sentiment

	// publish augmented tweet
	content, err := json.Marshal(tweet)
	if err != nil {
		return errors.Wrap(err, "unable to serialize tweet map")
	}

	if err := client.PublishEvent(ctx, resultPubSubName, resultTopicName, content); err != nil {
		return errors.Wrapf(err,
			"error while publishing content to pubsub:%s, topic:%s", resultPubSubName, resultTopicName,
		)
	}

	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
