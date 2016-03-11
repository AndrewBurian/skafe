package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	audisp := flag.Bool("audisp", false, "Run as Audisp plugin (read stdin)")
	flag.Parse()

	// create logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Println("Started skafe-agent")
	defer logger.Println("Terminated skafe-agent")

	// create channels to pass events between process steps
	newEventChan := make(chan AuditEvent) // auditer -> enricher

	enrichedEventChan := make(chan AuditEvent) // enricher -> cache

	sendEventChan := make(chan AuditEvent) // cache -> server

	go ServerLink(sendEventChan, logger)
	go Cache(enrichedEventChan, sendEventChan, logger)
	go Enricher(newEventChan, enrichedEventChan, logger)

	if *audisp {
		Audisp(newEventChan, os.Stdin)
	} else {
		Auditor(newEventChan, logger)
	}
}
