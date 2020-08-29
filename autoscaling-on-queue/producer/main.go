package main

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

const (
	min   = 3
	max   = 9999
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	logger = log.New(os.Stdout, "", 0)

	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	brokerAddress   = getEnvVar("KAFKA_BROKER", "localhost:9092")
	topicName       = getEnvVar("KAFKA_TOPIC", "messages")
	numOfThreadsStr = getEnvVar("NUMBER_OF_THREADS", "1")
	threadFreqStr   = getEnvVar("THREAD_PUB_FREQ", "10ms")
	threadFreq      time.Duration
)

func main() {
	numOfThreads, err := strconv.Atoi(numOfThreadsStr)
	if err != nil || numOfThreads < 1 {
		logger.Fatalf(
			"invalid number of thread (NUMBER_OF_THREADS must be positive int): %s - %v",
			numOfThreadsStr, err,
		)
	}
	logger.Printf("number of thread: %d", numOfThreads)

	freq, err := time.ParseDuration(threadFreqStr)
	if err != nil {
		logger.Fatalf(
			"invalid thread frequency (THREAD_PUB_FREQ) must be a duration): %s - %v",
			threadFreqStr, err,
		)
	}
	threadFreq = freq
	logger.Printf("thread frequency: %v", threadFreq)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	p, err := sarama.NewAsyncProducer(strings.Split(brokerAddress, ","), config)
	if err != nil {
		logger.Fatalf("error creating producer: %v", err)
	}
	defer p.AsyncClose()

	stopCh := make(chan struct{})
	resultCh := make(chan bool, 100)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go processResponse(p, resultCh)

	go func() {
		<-c
		close(stopCh)
	}()

	for i := 1; i <= numOfThreads; i++ {
		go publish(p, stopCh)
	}

	var mux sync.Mutex
	var successCounter int64 = 0
	var errorCounter int64 = 0
	startTime := time.Now()

	tickerCh := time.NewTicker(3 * time.Second).C
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
			logger.Printf("published: %10d, errors: %3d - %3.0f/sec ", successCounter, errorCounter, avg)
		case <-stopCh:
			os.Exit(0)
		}
	}
}

func processResponse(p sarama.AsyncProducer, outCh chan<- bool) {
	for {
		select {
		case <-p.Successes():
			outCh <- true
		case err := <-p.Errors():
			logger.Printf("error publishing: %v", err)
			outCh <- false
		}
	}
}

func publish(p sarama.AsyncProducer, stopCh <-chan struct{}) {
	tickerCh := time.NewTicker(threadFreq).C
	for {
		select {
		case <-stopCh:
			return
		case <-tickerCh:
			publishOne(p)
		}
	}
}

type validationRequest struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
	Sha  string `json:"sha"`
	Time int64  `json:"time"`
}

func publishOne(p sarama.AsyncProducer) {
	r := validationRequest{
		ID:   uuid.New().String(),
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
	p.Input() <- &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.ByteEncoder(b),
	}
}

func getData(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[seededRand.Intn(len(chars))]
	}
	return string(b)
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
