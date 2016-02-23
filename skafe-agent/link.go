package main

import (
	"fmt"
	"log"
)

func ServerLink(events <-chan AuditEvent, logger *log.Logger) {

	for {
		logger.Println("Server ready to dispatch")

		event := <-events

		logger.Println("Sending event")
		fmt.Println(event)
	}
}
