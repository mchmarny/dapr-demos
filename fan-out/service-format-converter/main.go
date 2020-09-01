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

	serviceAddress = getEnvVar("ADDRESS", ":60012")

	sourcePubSubName = getEnvVar("SOURCE_PUBSUB_NAME", "fanout-source-pubsub")
	sourceTopicName  = getEnvVar("SOURCE_TOPIC_NAME", "events")

	targetServiceID  = getEnvVar("TARGET_SERVICE", "grpc-echo-service")
	targetMethodName = getEnvVar("TARGET_METHOD", "echo")
	targetFormat     = getEnvVar("TARGET_FORMAT", "xml")
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

func eventHandler(ctx context.Context, e *common.TopicEvent) error {
	logger.Printf("Event - PubsubName:%s, Topic:%s, ID:%s", e.PubsubName, e.Topic, e.ID)

	d, ok := e.Data.([]byte)
	if !ok {
		return errors.Errorf("invalid event data type: %T", e.Data)
	}

	var se SourceEvent
	if err := json.Unmarshal(d, &se); err != nil {
		return errors.Errorf("error parsing input content: %v", err)
	}

	var (
		me error
		b  []byte
		ct string
	)

	switch strings.ToLower(targetFormat) {
	case "json":
		b = d
		ct = "application/json"
	case "xml":
		if b, me = xml.Marshal(&e); me != nil {
			return errors.Errorf("error while converting content: %v", me)
		}
		ct = "application/xml"
	case "csv":
		b = []byte(fmt.Sprintf(`"%s",%f,%f,"%s"`,
			se.ID, se.Temperature, se.Humidity, time.Unix(se.Time, 0).Format(time.RFC3339)))
		ct = "text/csv"
	default:
		return errors.Errorf("invalid target format: %s", targetFormat)
	}
	logger.Printf("Target (%s): %s", targetFormat, b)

	content := &dapr.DataContent{Data: b, ContentType: ct}
	logger.Printf("Target: %+v", content)

	out, err := client.InvokeServiceWithContent(ctx, targetServiceID, targetMethodName, content)
	if err != nil {
		return errors.Wrap(err, "error invoking target binding")
	}
	logger.Printf("Response: %s", out)

	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
