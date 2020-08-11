package main

import (
	"net/http"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gin-gonic/gin"
)

var (
	clientError = gin.H{
		"error":   "Bad Request",
		"message": "Error processing your request, see logs for details",
	}
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      AppVersion,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func tweetHandler(c *gin.Context) {
	var t twitter.Tweet
	if err := c.ShouldBindJSON(&t); err != nil {
		logger.Printf("error binding tweet: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	logger.Printf("tweet: %v", t)

	c.JSON(http.StatusOK, gin.H{
		"request": time.Now(),
		"release": AppVersion,
		"status":  c.Request.RemoteAddr,
	})
}
