package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"
)

var (
	logger = log.New(os.Stdout, "", 0)
)

func main() {
	connStrPtr := flag.String("conn", "", "Event ID name")
	flag.Parse()

	ctx := context.Background()
	sender, err := newEventSender(ctx, *connStrPtr)
	if err != nil {
		logger.Fatalf("error while initializing sender: %v", err)
	}

	for {
		time.Sleep(time.Second * 1)
		sender.publish(ctx, getRoomReading())
	}
}
