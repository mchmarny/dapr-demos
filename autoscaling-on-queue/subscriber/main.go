package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"time"

	"net/http"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/pkg/errors"
)

const (
	primeStateKey = "high-prime"
)

var (
	logger = log.New(os.Stdout, "", 0)
	client dapr.Client

	address     = getEnvVar("ADDRESS", ":60022")
	bindingName = getEnvVar("BINDING_NAME", "autoscaling-kafka-queue")
	storeName   = getEnvVar("STORE_NAME", "prime-store")
)

func main() {
	// Dapr client
	c, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("error creating Dapr client: %v", err)
	}
	client = c
	defer client.Close()

	// Dapr service
	s, err := daprd.NewService(address)
	if err != nil {
		logger.Fatalf("failed to start the service: %v", err)
	}

	// Add binding
	if err := s.AddBindingInvocationHandler(bindingName, eventHandler); err != nil {
		logger.Fatalf("error adding topic subscription: %v", err)
	}

	// Start
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

func eventHandler(ctx context.Context, e *common.BindingEvent) (out []byte, err error) {
	if err := processRequest(ctx, e.Data); err != nil {
		logger.Printf("error processing request: %v", err)
		return nil, errors.Wrap(err, "error processing request")
	}
	return nil, nil
}

func processRequest(ctx context.Context, in []byte) error {
	var r calcRequest
	if err := json.Unmarshal(in, &r); err != nil {
		return errors.Wrap(err, "error serializing input data")
	}

	r.Prime = calcHighestPrime(&r)

	sr, err := getHighestPrime(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting highest prime")
	}

	logger.Printf("Highest prime for this request (max: %d): %d, all: %d.", r.Max, r.Prime, sr.Prime)
	if r.Prime > sr.Prime {
		bb, err := json.Marshal(r)
		if err != nil {
			return errors.Wrap(err, "error serializing request")
		}
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
