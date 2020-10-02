package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/pkg/errors"
)

var (
	logger     = log.New(os.Stdout, "", 0)
	address    = getEnvVar("ADDRESS", ":8080")
	pubSubName = getEnvVar("PUBSUB_NAME", "tweeter-pubsub")
	topicName  = getEnvVar("TOPIC_NAME", "tweets")
	storeName  = getEnvVar("STORE_NAME", "tweet-store")
	client     dapr.Client
)

func main() {
	// create a Dapr service
	s := daprd.NewService(address)

	// create a Dapr client
	c, err := dapr.NewClient()
	if err != nil {
		logger.Fatalf("error creating Dapr client: %v", err)
	}
	client = c
	defer client.Close()

	// add twitter input binding handler
	if err := s.AddBindingInvocationHandler("tweets", tweetHandler); err != nil {
		logger.Fatalf("error adding binding handler: %v", err)
	}

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error starting service: %v", err)
	}
}

func tweetHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	logger.Printf("Tweet (query: %s, traceID: %s)", in.Metadata["Query"], in.Metadata["Traceparent"])
	var m map[string]interface{}
	if err := json.Unmarshal(in.Data, &m); err != nil {
		return nil, errors.Wrap(err, "error deserializing event data")
	}

	k := fmt.Sprintf("tw-%s", m["id_str"])
	if err := client.SaveState(ctx, storeName, k, in.Data); err != nil {
		return nil, errors.Wrapf(err, "error saving to store:%s with key:%s", storeName, k)
	}

	logger.Printf("Tweet saved in store: %s: %s", storeName, k)
	if err := client.PublishEvent(ctx, pubSubName, topicName, in.Data); err != nil {
		return nil, errors.Wrapf(err, "error publishing to %s/%s", pubSubName, topicName)
	}
	logger.Printf("Tweet published to %s/%s", pubSubName, topicName)
	return
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
