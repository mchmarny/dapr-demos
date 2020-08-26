package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"net/http"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const (
	primeStateKey = "high-prime"
)

var (
	logger = log.New(os.Stdout, "", 0)
	client dapr.Client

	address    = getEnvVar("ADDRESS", ":8089")
	pubSubName = getEnvVar("PUBSUB_NAME", "autoscaling-kafka-queue")
	topicName  = getEnvVar("TOPIC_NAME", "primes")
	storeName  = getEnvVar("STORE_NAME", "prime-store")
)

func main() {
	// create Dapr client
	c, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("error creating Dapr client: %v", err)
	}
	client = c
	defer client.Close()

	// create a Dapr service
	s := daprd.NewService(address)

	// add some topic subscriptions
	subscription := &common.Subscription{
		PubsubName: pubSubName,
		Topic:      topicName,
		Route:      fmt.Sprintf("/%s", topicName),
	}

	if err := s.AddTopicEventHandler(subscription, eventHandler); err != nil {
		logger.Fatalf("error adding topic subscription: %v", err)
	}

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error starting service: %v", err)
	}
}

type calcRequest struct {
	ID    string `json:"id"`
	Max   int    `json:"max"`
	Prime int    `json:"prime"`
	Time  int64  `json:"time"`
}

func eventHandler(ctx context.Context, e *common.TopicEvent) error {
	logger.Printf("Request - PubSub:%s, Topic:%s, ID:%s", e.PubsubName, e.Topic, e.ID)
	if err := processRequest(ctx, e.Data); err != nil {
		logger.Printf("error processing request: %v", err)
		return errors.Wrap(err, "error processing request")
	}
	return nil
}

func processRequest(ctx context.Context, in interface{}) error {
	var r calcRequest
	if err := mapstructure.Decode(in, &r); err != nil {
		return errors.Wrap(err, "error serializing input data")
	}

	r.Prime = calcHighestPrime(&r)
	logger.Printf("Highest prime for %d is %d", r.Max, r.Prime)

	sr, err := getHighestPrime(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting highest prime")
	}

	logger.Printf("Previous high: %d, New: %d", sr.Prime, r.Prime)
	if r.Prime > sr.Prime {
		bb, err := json.Marshal(r)
		if err != nil {
			return errors.Wrap(err, "error serializing request")
		}
		logger.Printf("Saving new high: %d", r.Prime)
		if err := client.SaveState(ctx, storeName, primeStateKey, bb); err != nil {
			return errors.Errorf("error saving prime content: %v", err)
		}
	}

	return nil
}

func getHighestPrime(ctx context.Context) (r *calcRequest, err error) {
	item, err := client.GetState(ctx, storeName, primeStateKey)
	if err != nil {
		logger.Printf("error quering store: %v", err)
		return nil, errors.Wrapf(err, "error quering state store: %s for key: %s", storeName, primeStateKey)
	}
	if item == nil || item.Value == nil {
		return &calcRequest{
			Prime: 0,
			ID:    "id0",
			Max:   0,
			Time:  time.Now().UTC().Unix(),
		}, nil
	}
	var sr calcRequest
	if err := json.Unmarshal(item.Value, &sr); err != nil {
		return nil, errors.Wrap(err, "error parsing saved request content")
	}
	return &sr, nil
}

func calcHighestPrime(r *calcRequest) int {
	h := 0
	for i := 2; i <= r.Max; i++ {
		if isPrime(i) {
			h = i
		}
	}
	return h
}

func isPrime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
