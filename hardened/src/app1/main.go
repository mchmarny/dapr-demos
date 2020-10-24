package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/pkg/errors"
)

const (
	methodName      = "ping"
	invocationError = "error invoking"
)

var (
	logger = log.New(os.Stdout, "", 0)

	serviceAddress    = getEnvVar("ADDRESS", ":8081")
	targetServiceName = getEnvVar("TARGET_SERVICE_NAME", "app2")
	targetMethodName  = getEnvVar("TARGET_METHOD_NAME", "counter")

	client dapr.Client
)

func main() {
	s := daprd.NewService(serviceAddress)

	var clientErr error
	if client, clientErr = dapr.NewClient(); clientErr != nil {
		logger.Fatalf("error creating Dapr client: %v", clientErr)
	}
	defer client.Close()

	if err := s.AddServiceInvocationHandler(methodName, handler); err != nil {
		logger.Fatalf("error adding call invocation handler: %v", err)
	}

	logger.Printf("starting server at %s...", serviceAddress)
	if err := s.Start(); err != nil {
		logger.Fatalf("error starting server: %v", err)
	}
}

func handler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	logger.Printf("Method %s invoked  (ContentType:%s, Verb:%s, QueryString:%s, Data:%s)",
		methodName, in.ContentType, in.Verb, in.QueryString, in.Data)

	m := &dapr.DataContent{ContentType: in.ContentType, Data: in.Data}

	data, err := client.InvokeServiceWithContent(ctx, targetServiceName, targetMethodName, m)
	if err != nil {
		logger.Printf("%s %s/%s: %v", invocationError, targetServiceName, targetMethodName, err)
		return nil, errors.Wrapf(err, "%s %s/%s", invocationError, targetServiceName, targetMethodName)
	}
	logger.Printf("Response from invocation: %s", data)

	j := []byte(fmt.Sprintf(`{"on": %d, "count": %s}`, time.Now().UTC().UnixNano(), string(data)))

	out = &common.Content{ContentType: in.ContentType, Data: j}
	return
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
