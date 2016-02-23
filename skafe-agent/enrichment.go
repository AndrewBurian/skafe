package main

import (
	"log"
)

func Enricher(newEvents <-chan AuditEvent, enrichedEvents chan<- AuditEvent, logger *log.Logger) {

	for {
		// get the event
		event := <-newEvents

		logger.Println("Event received by enricher")

		// enrich the event
		event["data"] = "enriched shit"

		logger.Println("Dispatching enriched event")
		enrichedEvents <- event
	}
}
