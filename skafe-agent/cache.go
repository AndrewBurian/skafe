package main

import (
	"log"
)

func Cache(receivedEvents <-chan AuditEvent, server chan<- AuditEvent, logger *log.Logger) {

	for {
		event := <-receivedEvents

		server <- event

	}

}
