package main

import (
	"context"
	"log"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/pkg/errors"
)

var (
	logger = log.New(os.Stdout, "", 0)

	serviceAddress    = getEnvVar("ADDRESS", ":8080")
	pubMethodName     = getEnvVar("METHOD_PUBLISH_NAME", "pub")
	callMethodName    = getEnvVar("METHOD_INVOKE_NAME", "call")
	pubSubName        = getEnvVar("PUBSUB_NAME", "pubsub")
	topicName         = getEnvVar("TOPIC_NAME", "messages")
	targetServiceName = getEnvVar("TARGET_SERVICE_NAME", "app2")
	targetMethodName  = getEnvVar("TARGET_METHOD_NAME", "ping")

	client dapr.Client
)

func main() {
	s := daprd.NewService(serviceAddress)

	var clientErr error
	if client, clientErr = dapr.NewClient(); clientErr != nil {
		logger.Fatalf("error creating Dapr client: %v", clientErr)
	}
	defer client.Close()

	if err := s.AddServiceInvocationHandler(callMethodName, callHandler); err != nil {
		logger.Fatalf("error adding call invocation handler: %v", err)
	}

	if err := s.AddServiceInvocationHandler(pubMethodName, pubHandler); err != nil {
		logger.Fatalf("error adding pub invocation handler: %v", err)
	}

	logger.Printf("starting server at %s...", serviceAddress)
	if err := s.Start(); err != nil {
		logger.Fatalf("error starting server: %v", err)
	}
}

func pubHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	logger.Printf("Method %s invoked  (ContentType:%s, Verb:%s, QueryString:%s, Data:%s)",
		pubMethodName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	if err := client.PublishEvent(ctx, pubSubName, topicName, in.Data); err != nil {
		logger.Printf("error publishing to %s/%s: %v", pubSubName, topicName, err)
		return nil, errors.Wrapf(err, "error publishing to %s/%s", pubSubName, topicName)
	}
	return
}

func callHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	logger.Printf("Method %s invoked  (ContentType:%s, Verb:%s, QueryString:%s, Data:%s)",
		callMethodName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	m := &dapr.DataContent{
		ContentType: in.ContentType,
		Data:        in.Data,
	}

	data, err := client.InvokeServiceWithContent(ctx, targetServiceName, targetMethodName, m)
	if err != nil {
		return nil, errors.Wrapf(err, "error invoking %s/%s", targetServiceName, targetMethodName)
	}
	logger.Printf("Response from invocation: %s", data)

	out = &common.Content{ContentType: in.ContentType, Data: data}
	return
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
