package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	logger = log.New(os.Stdout, "", 0)

	serviceAddress = getEnvVar("ADDRESS", ":60013")

	targetPubSubName = getEnvVar("TARGET_PUBSUB_NAME", "fanout-source-pubsub")
	targetTopicName  = getEnvVar("TARGET_TOPIC_NAME", "events")

	threadFreq = getEnvVar("THREAD_PUB_FREQ", "3s")
)

func main() {
	ctx := context.Background()

	tf, err := time.ParseDuration(threadFreq)
	if err != nil {
		logger.Fatalf("invalid thread frequency, expected duration: %s - %v", threadFreq, err)
	}
	logger.Printf("thread frequency: %s", tf)

	// create Dapr service
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		logger.Fatalf("failed to start the server: %v", err)
	}
	defer s.Stop()

	// dapr client
	c, err := dapr.NewClient()
	if err != nil {
		logger.Fatalf("failed to create Dapr client: %v", err)
	}
	defer c.Close()

	// timer
	timer := time.NewTicker(tf)
	defer timer.Stop()

	// produce
	go func() {
		if err := produce(ctx, c, timer); err != nil {
			logger.Fatalf("error: %v", err)
		}
	}()

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func produce(ctx context.Context, c dapr.Client, t *time.Ticker) error {
	for {
		select {
		case <-t.C:
			b, err := json.Marshal(getRoomReading())
			if err != nil {
				return errors.Wrap(err, "error serializing reading")
			}
			if err := c.PublishEvent(ctx, targetPubSubName, targetTopicName, b); err != nil {
				return errors.Wrap(err, "error publishing content")
			}
			logger.Printf("published: %s", b)
		}
	}
}

type roomReading struct {
	ID          string  `json:"id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Time        int64   `json:"time"`
}

func getRoomReading() interface{} {
	min := 0.01
	max := 100.00
	return &roomReading{
		ID:          uuid.New().String(),
		Temperature: min + rand.Float64()*(max-min),
		Humidity:    min + rand.Float64()*(max-min),
		Time:        time.Now().UTC().Unix(),
	}
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
