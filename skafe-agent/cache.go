package main

import (
	"log"
)

func Cache(receivedEvents <-chan AuditEvent, server chan<- AuditEvent, logger *log.Logger) {

	for {
		event := <-receivedEvents
		logger.Println("Cache Received event")

		server <- event
		logger.Println("Cache Sent event")

	}

}
