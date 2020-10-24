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
	savingStateError        = "error saving state"
	eventSerializationError = "error serializing event data"
	addInvokeHandlerError   = "error adding invocation handler"
	addSubscriptionError    = "error adding topic subscription"
	startingServiceError    = "error starting service"
	clientCreateError       = "error creating Dapr client"
)

var (
	logger     = log.New(os.Stdout, "", 0)
	address    = getEnvVar("ADDRESS", ":8081")
	methodName = getEnvVar("METHOD_NAME", "ping")
	pubSubName = getEnvVar("PUBSUB_NAME", "pubsub")
	topicName  = getEnvVar("TOPIC_NAME", "messages")
	storeName  = getEnvVar("STORE_NAME", "state")

	client dapr.Client
)

func main() {
	s := daprd.NewService(address)

	var clientErr error
	if client, clientErr = dapr.NewClient(); clientErr != nil {
		log.Fatalf("%s: %v", clientCreateError, clientErr)
	}
	defer client.Close()

	subscription := &common.Subscription{
		PubsubName: pubSubName,
		Topic:      topicName,
		Route:      fmt.Sprintf("/%s", topicName),
	}

	if err := s.AddServiceInvocationHandler(methodName, invokeHandler); err != nil {
		logger.Fatalf("%s: %v", addInvokeHandlerError, err)
	}

	if err := s.AddTopicEventHandler(subscription, eventHandler); err != nil {
		logger.Fatalf("%s: %v", addSubscriptionError, err)
	}

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("%s: %v", startingServiceError, err)
	}
}

func invokeHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	logger.Printf("Method %s invoked  (ContentType:%s, Verb:%s, QueryString:%s, Data:%s)",
		methodName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	now := time.Now()
	key := fmt.Sprintf("id-%d", now.UnixNano())

	if err := client.SaveState(ctx, storeName, key, in.Data); err != nil {
		logger.Printf("%s: %v", savingStateError, err)
		return nil, errors.Wrap(err, savingStateError)
	}

	data := fmt.Sprintf(`{"invocation": "%s", "key": "%s", "input": "%s"}`, in.Data, now.Format(time.RFC3339), key)
	out = &common.Content{ContentType: in.ContentType, Data: []byte(data)}
	return
}

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	logger.Printf("Event received (PubsubName:%s, Topic:%s, Data: %v", e.PubsubName, e.Topic, e.Data)

	data, err := json.Marshal(e.Data)
	if err != nil {
		logger.Printf("%s: %v", eventSerializationError, err)
		return false, errors.Wrap(err, eventSerializationError)
	}

	if err := client.SaveState(ctx, storeName, e.ID, data); err != nil {
		logger.Printf("%s: %v", savingStateError, err)
		return false, errors.Wrap(err, savingStateError)
	}
	logger.Printf("Event saved (ID:%s, Data: %v", e.ID, e.Data)

	return false, nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
