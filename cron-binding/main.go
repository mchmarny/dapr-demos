package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.OPTIONS("/run", optionsHandler)
	r.POST("/run", runHandler)
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func optionsHandler(c *gin.Context) {
	c.Header("Allow", "POST")
	c.Header("Content-Type", "application/json")
	c.AbortWithStatus(http.StatusOK)
}

func runHandler(c *gin.Context) {
	// TODO: do something interesting here
	log.Printf("invocation received: %v", time.Now())
	c.JSON(http.StatusOK, nil)
}
