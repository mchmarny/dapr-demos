package main

import (
	"context"
	"encoding/json"

	hub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/pkg/errors"
)

func newEventSender(ctx context.Context, connStr string) (*EventSender, error) {
	if connStr == "" {
		return nil, errors.New("connStr not defined")
	}
	c, err := hub.NewHubFromConnectionString(connStr)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating client from: '%s'", connStr)
	}

	_, err = c.GetRuntimeInformation(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}

	return &EventSender{
		client: c,
	}, nil
}

// EventSender sends events
type EventSender struct {
	client *hub.Hub
}

// Close closes the client connection
func (s *EventSender) Close() error {
	if s.client != nil {
		return s.client.Close(context.Background())
	}
	return nil
}

// Publish sends provied events to Event Hub
func (s *EventSender) publish(ctx context.Context, e *RoomReading) error {
	data, _ := json.Marshal(e)
	ev := hub.NewEvent(data)
	ev.ID = e.ID
	ev.Properties = make(map[string]interface{})
	ev.Properties["src"] = "Dapr DT Demo"

	logger.Printf("sending: %s", string(data))
	if err := s.client.Send(ctx, ev); err != nil {
		if !errors.Is(err, context.Canceled) {
			return errors.Wrapf(err, "error on publish: '%+v'", e)
		}
	}

	return nil
}
