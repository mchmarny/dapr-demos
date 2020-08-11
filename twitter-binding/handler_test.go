package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTweetHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/tweets", tweetHandler)
	w := httptest.NewRecorder()

	data, err := ioutil.ReadFile("./tweet.json")
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/tweets", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
