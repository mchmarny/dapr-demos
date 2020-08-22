package main

import (
	"context"
	"log"

	"net/http"
	"os"
	"strings"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	logger  = log.New(os.Stdout, "", 0)
	address = getEnvVar("ADDRESS", ":8080")
)

func main() {
	// create a Dapr service
	s := daprd.NewService(address)

	// add some input binding handler
	if err := s.AddBindingInvocationHandler("schedule", scheduleHandler); err != nil {
		logger.Fatalf("error adding binding handler: %v", err)
	}

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error starting service: %v", err)
	}
}

func scheduleHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	logger.Printf("Schedule - Metadata:%v, Data:%v", in.Metadata, in.Data)

	// TODO: do something with the cloud event data

	return nil, nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
