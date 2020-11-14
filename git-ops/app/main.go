package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	daprd "github.com/dapr/go-sdk/service/http"
)

const (
	staticMessage = "hello PDX"
)

var (
	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// BuildTime will be overritten during build
	BuildTime = "not set"

	// service
	logger  = log.New(os.Stdout, "", 0)
	address = getEnvVar("ADDRESS", ":8080")

	templates *template.Template
)

func main() {

	// server mux
	mux := http.NewServeMux()

	// static content
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("resource/static"))))
	mux.HandleFunc("/favicon.ico", faviconHandler)

	// tempalates
	templates = template.Must(template.ParseGlob("resource/template/*"))

	// other handlers
	mux.HandleFunc("/", rootHandler)

	// create a Dapr service
	s := daprd.NewServiceWithMux(address, mux)

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error starting service: %v", err)
	}
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
		"message": staticMessage,
		"built":   BuildTime,
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

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
