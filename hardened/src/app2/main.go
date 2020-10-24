package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net/http"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/pkg/errors"
)

const (
	methodName      = "counter"
	stateCounterKey = "count"

	addInvokeHandlerError = "error adding invocation handler"
	startingServiceError  = "error starting service"
	clientCreateError     = "error creating Dapr client"
	publisingError        = "error publishing"
	getStateError         = "error getting state"
	saveStateError        = "error saving state"
	intParsingError       = "error parsing data"
)

var (
	logger  = log.New(os.Stdout, "", 0)
	address = getEnvVar("ADDRESS", ":8082")

	pubSubName = getEnvVar("PUBSUB_NAME", "pubsub")
	topicName  = getEnvVar("TOPIC_NAME", "messages")

	storeName = getEnvVar("STORE_NAME", "state")

	client dapr.Client
)

func main() {
	s := daprd.NewService(address)

	var clientErr error
	if client, clientErr = dapr.NewClient(); clientErr != nil {
		logger.Fatalf("%s: %v", clientCreateError, clientErr)
	}
	defer client.Close()

	if err := s.AddServiceInvocationHandler(methodName, handler); err != nil {
		logger.Fatalf("%s: %v", addInvokeHandlerError, err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("%s: %v", startingServiceError, err)
	}
}

func handler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	logger.Printf("Method %s invoked  (ContentType:%s, Verb:%s, QueryString:%s, Data:%s)",
		methodName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	// get previous number
	data, err := client.GetState(ctx, storeName, stateCounterKey)
	if err != nil {
		logger.Printf("%s: %v", getStateError, err)
		return nil, errors.Wrap(err, getStateError)
	}

	// convert the content to number
	var count int64 = 0
	var cnvErr error
	if data != nil && data.Value != nil {
		count, cnvErr = strconv.ParseInt(fmt.Sprintf("%s", data.Value), 10, 64)
		if cnvErr != nil {
			logger.Printf("%s %v: %v", intParsingError, data.Value, cnvErr)
			return nil, errors.Wrap(cnvErr, intParsingError)
		}
	}

	// increment
	count = count + 1
	logger.Printf("New count: %d", count)

	// data from counter
	b := []byte(fmt.Sprintf("%d", count))

	// save new number
	item := &dapr.SetStateItem{
		Etag:  data.Etag, // using the retreaved etag to ensure count consistency
		Key:   stateCounterKey,
		Value: b,
		Options: &dapr.StateOptions{
			Concurrency: dapr.StateConcurrencyLastWrite,
			Consistency: dapr.StateConsistencyStrong,
		},
	}
	if err := client.SaveStateItems(ctx, storeName, item); err != nil {
		logger.Printf("%s: %v", saveStateError, err)
		return nil, errors.Wrap(err, saveStateError)
	}

	// publish the results for processing down the stream
	if err := client.PublishEvent(ctx, pubSubName, topicName, b); err != nil {
		logger.Printf("%s to %s/%s: %v", publisingError, pubSubName, topicName, err)
		return nil, errors.Wrapf(err, "%s to %s/%s", publisingError, pubSubName, topicName)
	}

	out = &common.Content{ContentType: in.ContentType, Data: b}

	return
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
