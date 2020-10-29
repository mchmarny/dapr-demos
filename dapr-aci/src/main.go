package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"net/http"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/pkg/errors"
)

const (
	addInvokeHandlerError = "error adding invocation handler"
	startingServiceError  = "error starting service"
	clientCreateError     = "error creating Dapr client"
	addSubscriptionError  = "error subscribing to a topic"
	addInvocationError    = "error creating invcation handler"
	methodName            = "ping"
)

var (
	logger  = log.New(os.Stdout, "", 0)
	address = getEnvVar("ADDRESS", ":8082")

	pubSubName = getEnvVar("PUBSUB_NAME", "pubsub")
	topicName  = getEnvVar("TOPIC_NAME", "messages")

	storeName = getEnvVar("STORE_NAME", "store")

	client dapr.Client
)

func main() {
	s := daprd.NewService(address)

	var clientErr error
	if client, clientErr = dapr.NewClient(); clientErr != nil {
		logger.Fatalf("%s: %v", clientCreateError, clientErr)
	}
	defer client.Close()

	if err := s.AddServiceInvocationHandler(methodName, invokeHandler); err != nil {
		logger.Fatalf("%s: %v", addInvocationError, err)
	}

	subscription := &common.Subscription{
		PubsubName: pubSubName,
		Topic:      topicName,
		Route:      fmt.Sprintf("/%s", topicName),
	}

	if err := s.AddTopicEventHandler(subscription, eventHandler); err != nil {
		logger.Fatalf("%s: %v", addSubscriptionError, err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("%s: %v", startingServiceError, err)
	}
}

func invokeHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	logger.Printf("Method %s invoked  (ContentType:%s, Verb:%s, QueryString:%s, Data:%s)",
		methodName, in.ContentType, in.Verb, in.QueryString, in.Data)
	j := []byte(fmt.Sprintf(`{"on": %d, "greeting": "pong"}`, time.Now().UTC().UnixNano()))
	out = &common.Content{ContentType: in.ContentType, Data: j}
	return
}

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	logger.Printf("Event received (PubsubName:%s, Topic:%s, Data: %v", e.PubsubName, e.Topic, e.Data)

	data, ok := e.Data.([]byte)
	if !ok {
		data, err = json.Marshal(e.Data)
		if err != nil {
			return false, errors.Wrapf(err, "invalid data format: %T", e.Data)
		}
	}

	if err := client.SaveState(ctx, storeName, e.ID, data); err != nil {
		return false, errors.Wrapf(err, "error saving data to %s (%s)", storeName, data)
	}

	return false, nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
