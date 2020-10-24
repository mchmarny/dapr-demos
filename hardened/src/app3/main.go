package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net/http"
	"os"
	"strings"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/pkg/errors"
)

const (
	eventDataParsingError = "error parsing event data"
	addSubscriptionError  = "error adding topic subscription"
	startingServiceError  = "error starting service"
	intParsingError       = "error parsing data"
)

var (
	logger  = log.New(os.Stdout, "", 0)
	address = getEnvVar("ADDRESS", ":8083")

	pubSubName = getEnvVar("PUBSUB_NAME", "pubsub")
	topicName  = getEnvVar("TOPIC_NAME", "messages")
)

func main() {
	s := daprd.NewService(address)

	subscription := &common.Subscription{
		PubsubName: pubSubName,
		Topic:      topicName,
		Route:      fmt.Sprintf("/%s", topicName),
	}

	if err := s.AddTopicEventHandler(subscription, handler); err != nil {
		logger.Fatalf("%s: %v", addSubscriptionError, err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("%s: %v", startingServiceError, err)
	}
}

func handler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	logger.Printf("Event received (PubsubName:%s, Topic:%s, Data: %v", e.PubsubName, e.Topic, e.Data)

	var (
		count  int64 = 0
		cnvErr error
	)

	data := fmt.Sprintf("%v", e.Data)
	count, cnvErr = strconv.ParseInt(data, 10, 64)
	if cnvErr != nil {
		logger.Printf("%s %v: %v", intParsingError, data, cnvErr)
		return false, errors.Wrap(cnvErr, intParsingError)
	}

	logger.Printf("Even: %v", count%2 == 0)

	return false, nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
