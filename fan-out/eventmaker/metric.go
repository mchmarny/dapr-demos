package main

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// RoomReading represents DT room
type RoomReading struct {
	ID          string  `json:"id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Time        int64   `json:"time"`
}

func getRoomReading() *RoomReading {
	min := 0.01
	max := 100.00
	return &RoomReading{
		ID:          uuid.New().String(),
		Temperature: min + rand.Float64()*(max-min),
		Humidity:    min + rand.Float64()*(max-min),
		Time:        time.Now().UTC().Unix(),
	}
}
