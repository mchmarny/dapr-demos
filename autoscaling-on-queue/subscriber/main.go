package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"
	"os"
	"strings"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/pkg/errors"
)

const (
	primeStateKey = "high-prime"
)

var (
	logger     = log.New(os.Stdout, "", 0)
	reqProcDur time.Duration

	address         = getEnvVar("ADDRESS", ":60033")
	processDuration = getEnvVar("PROCESS_DURATION", "500ms")
	pubSubName      = getEnvVar("PUBSUB_NAME", "autoscaling-pubsub")
	topicName       = getEnvVar("TOPIC_NAME", "metrics")
)

func main() {
	// Dapr service
	s, err := daprd.NewService(address)
	if err != nil {
		logger.Fatalf("failed to start the service: %v", err)
	}

	d, err := time.ParseDuration(processDuration)
	if err != nil {
		logger.Fatalf("invalid parameter (PROCESS_DURATION) must be a duration): %s - %v", processDuration, err)
	}
	reqProcDur = d

	var mux sync.Mutex
	var successCount int64 = 1
	var errorCount int64 = 0

	resultCh := make(chan bool)
	startTime := time.Now()

	go func() {
		tickerCh := time.NewTicker(5 * time.Second).C
		for {
			select {
			case r := <-resultCh:
				mux.Lock()
				if r {
					successCount++
				} else {
					errorCount++
				}
				mux.Unlock()
			case <-tickerCh:
				var avg float64 = 0
				if successCount > 0 {
					avg = float64(successCount) / time.Since(startTime).Seconds()
				}
				logger.Printf("received: %10d, %3d errors - avg %3.0f/sec", successCount, errorCount, avg)
			}
		}
	}()

	// define subscription
	subscription := &common.Subscription{
		PubsubName: pubSubName,
		Topic:      topicName,
	}

	// subscribe
	if err := s.AddTopicEventHandler(subscription, func(ctx context.Context, e *common.TopicEvent) error {
		if err := processRequest(ctx, e.Data); err != nil {
			logger.Printf("error processing request: %v", err)
			resultCh <- false
			return errors.Wrap(err, "error processing request")
		}
		resultCh <- true
		return nil
	}); err != nil {
		logger.Fatalf("error adding topic subscription: %v", err)
	}

	// handle signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("error starting service: %v", err)
		}
	}()

	// Finish
	<-done
}

// does some computing to keep the process busy organically
func processRequest(ctx context.Context, in interface{}) error {
	tickerCh := time.NewTicker(reqProcDur).C
	<-tickerCh
	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
