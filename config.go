package main

import (
	"encoding/json"

	"github.com/r3labs/sse/v2"
)

type Service struct {
	ID       string
	Name     string
	Running  bool
	Building bool
}

var services = []Service{
	{
		ID:   "test",
		Name: "Test",
	},
}

func PublishServices() {
	v, err := json.Marshal(services)
	if err != nil {
		panic(err) // Should never happen
	}
	events.Publish("services", &sse.Event{
		Data: v,
	})
}
