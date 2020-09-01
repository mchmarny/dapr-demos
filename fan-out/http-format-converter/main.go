package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/pkg/errors"
)

var (
	logger = log.New(os.Stdout, "", 0)
	client dapr.Client

	serviceAddress = getEnvVar("ADDRESS", ":60011")

	sourcePubSubName = getEnvVar("SOURCE_PUBSUB_NAME", "fanout-source-pubsub")
	sourceTopicName  = getEnvVar("SOURCE_TOPIC_NAME", "events")

	targetBindingName = getEnvVar("TARGET_BINDING", "fanout-http-target-post-binding")
	targetFormat      = getEnvVar("TARGET_FORMAT", "json")
)

func main() {
	// create Dapr service
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	c, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("failed to create Dapr client: %v", err)
	}
	client = c
	defer client.Close()

	// add handler to the service
	sub := &common.Subscription{PubsubName: sourcePubSubName, Topic: sourceTopicName}
	s.AddTopicEventHandler(sub, eventHandler)

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// SourceEvent represents the input event
type SourceEvent struct {
	ID          string  `json:"id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Time        int64   `json:"time"`
}

func eventHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	logger.Printf("Source: %s", in.Data)

	var e SourceEvent
	if err := json.Unmarshal(in.Data, &e); err != nil {
		return nil, errors.Errorf("error parsing input content: %v", err)
	}

	var (
		me error
		b  []byte
	)

	switch strings.ToLower(targetFormat) {
	case "json":
		b = in.Data
	case "xml":
		if b, me = xml.Marshal(&e); me != nil {
			return nil, errors.Errorf("error while converting content: %v", me)
		}
	case "csv":
		b = []byte(fmt.Sprintf(`"%s",%f,%f,"%s"`,
			e.ID, e.Temperature, e.Humidity, time.Unix(e.Time, 0).Format(time.RFC3339)))
	default:
		return nil, errors.Errorf("invalid target format: %s", targetFormat)
	}
	logger.Printf("Target: %s", b)

	content := &dapr.BindingInvocation{
		Data: b,
		Metadata: map[string]string{
			"record-id":       e.ID,
			"conversion-time": time.Now().UTC().Format(time.RFC3339),
		},
		Name:      targetBindingName,
		Operation: "create",
	}

	if err := client.InvokeOutputBinding(ctx, content); err != nil {
		return nil, errors.Wrap(err, "error invoking target binding")
	}

	return nil, nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
