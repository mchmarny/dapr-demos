package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/olahol/melody.v1"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/pkg/errors"
)

var (
	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// service
	logger     = log.New(os.Stdout, "", 0)
	address    = getEnvVar("ADDRESS", ":8083")
	pubSubName = getEnvVar("PUBSUB_NAME", "queue")
	topicName  = getEnvVar("TOPIC_NAME", "processed")

	broadcaster *melody.Melody
	templates   *template.Template
)

func main() {

	// server mux
	mux := http.NewServeMux()

	// static content
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("resource/static"))))
	mux.HandleFunc("/favicon.ico", faviconHandler)

	// tempalates
	templates = template.Must(template.ParseGlob("resource/template/*"))

	// websocket upgrade
	broadcaster = melody.New()
	broadcaster.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// other handlers
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/ws", wsHandler)

	// create a Dapr service
	s := daprd.NewServiceWithMux(address, mux)

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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	broadcaster.HandleRequest(w, r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	proto := r.Header.Get("x-forwarded-proto")
	if proto == "" {
		proto = "http"
	}

	data := map[string]string{
		"host":    r.Host,
		"proto":   proto,
		"version": AppVersion,
	}

	err := templates.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./resource/static/img/favicon.ico")
}

func eventHandler(ctx context.Context, e *common.TopicEvent) error {
	logger.Printf(
		"event - PubsubName:%s, Topic:%s, ID:%s, Data: %v",
		e.PubsubName, e.Topic, e.ID, e.Data,
	)

	b, err := json.Marshal(e.Data)
	if err != nil {
		return errors.Wrap(err, "error marshaling data")
	}

	broadcaster.Broadcast(b)
	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
