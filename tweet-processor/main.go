package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/pkg/errors"
)

var (
	logger         = log.New(os.Stdout, "", 0)
	serviceAddress = getEnvVar("ADDRESS", ":60002")

	srcPubSubName = getEnvVar("SOURCE_PUBSUB_NAME", "tweeter-pubsub")
	srcTopicName  = getEnvVar("SOURCE_TOPIC_NAME", "tweets")

	resultPubSubName = getEnvVar("RESULT_PUBSUB_NAME", "processed-tweets-pubsub")
	resultTopicName  = getEnvVar("RESULT_TOPIC_NAME", "processed-tweets")

	sentimentServiceName = getEnvVar("SENTIMENT_SERVICE_NAME", "sentiment-scorer")

	client dapr.Client
)

func main() {
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

	// add handler to the service
	subscription := &common.Subscription{
		PubsubName: srcPubSubName,
		Topic:      srcTopicName,
	}
	s.AddTopicEventHandler(subscription, tweetHandler)

	// start the server to handle incoming events
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func topicDataToSentimentRequest(b []byte) (s *SentimentRequest, err error) {
	var t TweetText
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, errors.Wrapf(err, "error deserializing tweet into request: %s", b)
	}

	s = &SentimentRequest{
		Text:     t.Text,
		Language: t.Lang,
	}

	if t.Extended.Text != "" {
		s.Text = t.Extended.Text
	}

	return
}

func getSentimentScore(ctx context.Context, req *SentimentRequest) (score *SentimentScore, err error) {
	cb, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to serialize sentiment request")
	}

	c := &dapr.DataContent{ContentType: "application/json", Data: cb}

	b, err := client.InvokeServiceWithContent(ctx, sentimentServiceName, "sentiment", c)
	if err != nil {
		return nil, errors.Wrap(err, "error invoking sentiment service")
	}

	var s SentimentScore
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, errors.Wrap(err, "error deserializing sentiment")
	}

	return &s, nil
}

func tweetHandler(ctx context.Context, e *common.TopicEvent) error {
	logger.Printf("Processing pubsub:%s/topic:%s id:%s", e.PubsubName, e.Topic, e.ID)

	b, ok := e.Data.([]byte)
	if !ok {
		return fmt.Errorf("invalid data type, expected []bytes: %T", e.Data)
	}

	sentReq, err := topicDataToSentimentRequest(b)
	if err != nil {
		return errors.Wrap(err, "error getting tweet text")
	}

	sentScore, err := getSentimentScore(ctx, sentReq)
	if err != nil {
		return errors.Wrap(err, "error getting sentiment score")
	}

	var tweetMap map[string]interface{}
	if err := json.Unmarshal(b, &tweetMap); err != nil {
		return errors.Wrap(err, "error deserializing content into map")
	}

	tweetMap["sentiment"] = sentScore
	content, err := json.Marshal(tweetMap)
	if err != nil {
		return errors.Wrap(err, "unable to serialize tweet map content")
	}

	if err := client.PublishEvent(ctx, resultPubSubName, resultTopicName, content); err != nil {
		return errors.Wrapf(err, "error publishing to %s/%s", resultPubSubName, resultTopicName)
	}

	logger.Printf("Processed tweet:%s - %v", e.ID, sentScore)
	return nil
}

// TweetText represents only the text of tweet for sentiment
type TweetText struct {
	Text     string `json:"text"`
	Lang     string `json:"lang"`
	Extended struct {
		Text string `json:"full_text"`
	} `json:"extended_tweet"`
}

// SentimentRequest represents the sentiment request
type SentimentRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

// SentimentScore represents sentiment result
type SentimentScore struct {
	Sentiment  string  `json:"sentiment"`
	Confidence float64 `json:"confidence"`
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
