package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
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

	threadCount = getEnvVar("NUMBER_OF_THREADS", "1")
	threadFreq  = getEnvVar("THREAD_PUB_FREQ", "1s")
)

func main() {
	// parse vars
	tc, err := strconv.Atoi(threadCount)
	if err != nil || tc < 1 {
		logger.Fatalf("invalid number of threads, expected positive int: %s - %v", threadCount, err)
	}
	logger.Printf("thread count: %d", tc)

	tf, err := time.ParseDuration(threadFreq)
	if err != nil {
		logger.Fatalf("invalid thread frequency, expected duration: %s - %v", threadFreq, err)
	}
	logger.Printf("thread frequency: %s", tf)

	// create Dapr service
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// handle signals
	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// produce
	go func() {
		for i := 1; i <= tc; i++ {
			if err := produce(i, tc, tf, doneCh); err != nil {
				logger.Fatalf("error: %v", err)
			}
		}
	}()

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func produce(ti, tc int, tf time.Duration, doneCh <-chan os.Signal) error {
	ctx := context.Background()
	// dapr client
	client, err := dapr.NewClient()
	if err != nil {
		return errors.Wrap(err, "failed to create Dapr client")
	}
	defer client.Close()

	tickerCh := time.NewTicker(tf).C
	for {
		select {
		case <-tickerCh:
			publishEvent(ctx, ti, client)
		case <-doneCh:
			return nil
		}
	}
}

func publishEvent(ctx context.Context, ti int, client dapr.Client) error {
	b, err := json.Marshal(getRoomReading())
	if err != nil {
		return errors.Wrap(err, "error serializing reading")
	}
	if err := client.PublishEvent(ctx, targetPubSubName, targetTopicName, b); err != nil {
		return errors.Wrap(err, "error publishing content")
	}
	logger.Printf("published: %s", b)
	return nil
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
