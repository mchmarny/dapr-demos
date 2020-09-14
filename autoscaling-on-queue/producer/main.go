package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/google/uuid"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	logger           = log.New(os.Stdout, "", 0)
	serviceAddress   = getEnvVar("ADDRESS", ":60034")
	pubSubName       = getEnvVar("PUBSUB_NAME", "autoscaling-pubsub")
	topicName        = getEnvVar("TOPIC_NAME", "metrics")
	numOfPublishers  = getEnvIntOrFail("NUMBER_OF_PUBLISHERS", "1")
	publishFrequency = getEnvDurationOrFail("PUBLISHERS_FREQ", "1s")
	publishDelay     = getEnvDurationOrFail("PUBLISHERS_DELAY", "10s")
	logFrequency     = getEnvDurationOrFail("LOG_FREQ", "3s")
	publishToConsole = getEnvBoolOrFail("PUBLISH_TO_CONSOLE", "false")

	client dapr.Client
)

func main() {
	if numOfPublishers < 1 {
		numOfPublishers = 1
	}
	logger.Printf("subscription name: %s", pubSubName)
	logger.Printf("topic name: %s", topicName)
	logger.Printf("number of publishers: %d", numOfPublishers)
	logger.Printf("publish frequency: %v", publishFrequency)
	logger.Printf("log frequency: %v", logFrequency)
	logger.Printf("publish delay: %v", publishDelay)

	// create Dapr service
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	c, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("error creating Dapr client: %v", err)
	}
	client = c
	defer client.Close()

	// handle signals
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	resultCh := make(chan bool, 100)
	stopCh := make(chan struct{})

	go func() {
		<-stop
		close(stopCh)
	}()

	// print results
	go monitor(resultCh, stopCh)

	// start producing
	for i := 1; i <= numOfPublishers; i++ {
		go publish(i, resultCh, stopCh)
	}

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}

}

func monitor(resultCh <-chan bool, stopCh <-chan struct{}) {
	var mux sync.Mutex
	var successCounter int64 = 0
	var errorCounter int64 = 0
	startTime := time.Now()
	tickerCh := time.NewTicker(logFrequency).C
	for {
		select {
		case r := <-resultCh:
			mux.Lock()
			if r {
				successCounter++
			} else {
				errorCounter++
			}
			mux.Unlock()
		case <-tickerCh:
			var avg float64 = 0
			if successCounter > 0 {
				avg = float64(successCounter) / time.Since(startTime).Seconds()
			}
			logger.Printf("%10d published, %3.0f/sec, %3d errors", successCounter, avg, errorCounter)
		case <-stopCh:
			os.Exit(0)
		}
	}
}

func publish(index int, resultCh chan<- bool, stopCh <-chan struct{}) {
	delayCh := time.NewTicker(publishDelay).C
	<-delayCh

	tickerCh := time.NewTicker(publishFrequency).C
	for {
		select {
		case <-stopCh:
			return
		case <-tickerCh:
			d := getEventData(index)
			if publishToConsole {
				logger.Printf("%s", d)
				resultCh <- true
				continue
			}
			resultCh <- client.PublishEvent(context.Background(), pubSubName, topicName, d) == nil
		}
	}
}

func getEventData(index int) []byte {
	r := requestContent{
		ID:   fmt.Sprintf("p%d-%s", index, uuid.New().String()),
		Data: []byte(getData(256)),
		Time: time.Now().UTC().Unix(),
	}

	// hash the entire message
	inSha := sha256.Sum256(r.Data)
	r.Sha = string(inSha[:])

	b, err := json.Marshal(r)
	if err != nil {
		logger.Fatalf("error generating request: %v", err)
	}
	return b
}

func getData(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[seededRand.Intn(len(chars))]
	}
	return string(b)
}

type requestContent struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
	Sha  string `json:"sha"`
	Time int64  `json:"time"`
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}

func getEnvIntOrFail(key, fallbackValue string) int {
	s := getEnvVar(key, fallbackValue)
	v, err := strconv.Atoi(s)
	if err != nil {
		logger.Fatalf("invalid number variable: %s - %v", s, err)
	}
	return v
}

func getEnvDurationOrFail(key, fallbackValue string) time.Duration {
	s := getEnvVar(key, fallbackValue)
	v, err := time.ParseDuration(s)
	if err != nil {
		logger.Fatalf("invalid duration variable: %s - %v", s, err)
	}
	return v
}

func getEnvBoolOrFail(key, fallbackValue string) bool {
	s := getEnvVar(key, fallbackValue)
	v, err := strconv.ParseBool(s)
	if err != nil {
		logger.Fatalf("invalid bool variable: %s - %v", s, err)
	}
	return v
}
