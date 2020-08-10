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
)

func main() {
	// create serving server
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// add handler to the service
	s.AddServiceInvocationHandler("echo", echoHandler)

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func echoHandler(ctx context.Context, in *daprd.InvocationEvent) (out *daprd.Content, err error) {
	logger.Printf(
		"Invocation (ContentType:%s, Verb:%s, QueryString:%s, Data:%s)",
		in.ContentType, in.Verb, in.QueryString, string(in.Data),
	)

	// TODO: implement handling logic here
	out = &daprd.Content{
		ContentType: in.ContentType,
		Data:        in.Data,
	}

	return
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
