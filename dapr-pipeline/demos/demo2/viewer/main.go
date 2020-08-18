package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"
	"gopkg.in/olahol/melody.v1"
)

var (
	logger = log.New(os.Stdout, "VIEWER == ", 0)

	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// service
	servicePort = env.MustGetEnvVar("PORT", "8083")

	sourceTopic = env.MustGetEnvVar("VIEWER_SOURCE_TOPIC_NAME", "processed")

	broadcaster *melody.Melody
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)

	// ws
	broadcaster = melody.New()
	broadcaster.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// static
	r.LoadHTMLGlob("resource/template/*")
	r.Static("/static", "./resource/static")
	r.StaticFile("/favicon.ico", "./resource/static/img/favicon.ico")

	// simple routes
	r.GET("/", rootHandler)

	// topic route
	viewerRoute := fmt.Sprintf("/%s", sourceTopic)
	logger.Printf("viewer route: %s", viewerRoute)
	r.POST(viewerRoute, eventHandler)

	// subscription
	r.GET("/dapr/subscribe", func(c *gin.Context) {
		data := []subscription{
			{
				Topic: sourceTopic,
				Route: viewerRoute,
			},
		}
		logger.Printf("subscription topics: %+v", data)
		c.JSON(http.StatusOK, data)
	})

	// websockets
	r.GET("/ws", func(c *gin.Context) {
		broadcaster.HandleRequest(c.Writer, c.Request)
	})

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := r.Run(hostPort); err != nil {
		logger.Fatal(err)
	}

}

type subscription struct {
	Topic string `json:"topic"`
	Route string `json:"route"`
}

// Options midleware
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "POST,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}
