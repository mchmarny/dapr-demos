package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/pkg/errors"
)

const (
	languageDefault = "en"
	secretStoreName = "pipeline-secrets"
	secretStoreKey  = "Azure:CognitiveAPIKey"
)

var (
	logger = log.New(os.Stdout, "", 0)

	serviceAddress = getEnvVar("ADDRESS", ":60005")
	apiToken       = getEnvVar("API_TOKEN", "")
	apiDomain      = getEnvVar("API_DOMAIN", "tweet-sentiment")

	apiURL = fmt.Sprintf("https://%s.cognitiveservices.azure.com/text/analytics/v3.0/sentiment", apiDomain)
)

func main() {

	// create serving server
	s, err := daprd.NewService(serviceAddress)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// add handler to the service
	s.AddServiceInvocationHandler("sentiment", sentimentHandler)

	// start the server to handle incoming events
	log.Printf("starting server at %s...", serviceAddress)
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func sentimentHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	logger.Printf("Processing: %s", in.Data)
	var req map[string]string
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, errors.Wrapf(err, "error deserializing data: %s", in.Data)
	}

	score, err := getSentiment(ctx, req["language"], req["text"])
	if err != nil {
		logger.Printf("error scoring sentiment: %v", err)
		return nil, errors.Wrapf(err, "error scoring sentiment: %s", in.Data)
	}

	b, err := json.Marshal(score)
	if err != nil {
		return nil, errors.Wrapf(err, "error serializing score: %v", score)
	}

	logger.Printf("Processed: %s", b)
	return &common.Content{
		ContentType: "application/json",
		Data:        b,
	}, nil
}

// SentimentScore represents sentiment result
type SentimentScore struct {
	Sentiment  string  `json:"sentiment"`
	Confidence float64 `json:"confidence"`
}

func getSentiment(ctx context.Context, lang, text string) (out *SentimentScore, err error) {
	if text == "" {
		return nil, errors.New("text required")
	}

	if lang == "" {
		lang = languageDefault
	}

	if apiToken == "" {
		apiToken = getSecret(secretStoreName, secretStoreKey)
	}

	r := fmt.Sprintf(`{
        "documents": [{
			"language": "%s",
            "id": "1",
			"text": "%s"
          }]
      }`, lang, text)

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer([]byte(r)))
	if err != nil {
		return nil, errors.Wrapf(err, "error creating request from: %v", r)
	}

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Ocp-Apim-Subscription-Key", apiToken)

	client := http.Client{Timeout: time.Second * 5}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error posting to: %s", apiURL)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid API response status: %d", res.StatusCode)
	}
	defer res.Body.Close()

	// dump, _ := httputil.DumpResponse(res, true)
	// logger.Printf("response: %s", dump)

	var rez struct {
		Documents []struct {
			Sentiment string `json:"sentiment"`
			Scores    struct {
				Positive float64 `json:"positive"`
				Neutral  float64 `json:"neutral"`
				Negative float64 `json:"negative"`
			} `json:"confidenceScores"`
		} `json:"documents"`
	}

	if err := json.NewDecoder(res.Body).Decode(&rez); err != nil {
		return nil, errors.Wrap(err, "error decoding API response")
	}

	if len(rez.Documents) != 1 {
		return nil, errors.Wrapf(err, "invalid response, expected 1 document, got %d", len(rez.Documents))
	}

	doc := rez.Documents[0]
	out = &SentimentScore{
		Sentiment: doc.Sentiment,
	}

	switch out.Sentiment {
	case "positive":
		out.Confidence = rez.Documents[0].Scores.Positive
	case "negative":
		out.Confidence = rez.Documents[0].Scores.Negative
	case "neutral":
		out.Confidence = rez.Documents[0].Scores.Neutral
	default:
		return nil, fmt.Errorf("invalid sentiment: %s", out.Sentiment)
	}

	return
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}

func getSecret(store, key string) string {
	// try to find it in Dapr secret store
	c, err := dapr.NewClient()
	if err != nil {
		logger.Fatal("unable to create Dapr client")
	}
	if m, err := c.GetSecret(context.Background(), store, key, map[string]string{}); err == nil {
		return m[key]
	}
	logger.Fatalf("no item found in Dapr secret store for %s", key)
	return ""
}
