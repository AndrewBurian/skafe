package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	auditd := flag.Bool("audit", false, "Run as system auditer")
	flag.Parse()

	// create logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Println("Started skafe-agent")
	defer logger.Println("Terminated skafe-agent")

	// create channels to pass events between process steps
	newEventChan := make(chan AuditEvent) // auditer -> enricher

	enrichedEventChan := make(chan AuditEvent) // enricher -> cache

	sendEventChan := make(chan AuditEvent) // cache -> rateLimiter

	serverChan := make(chan AuditEvent) // rateLimiter -> serverLink

	go ServerLink(serverChan, logger)
	go RateLimit(sendEventChan, serverChan, 10, 1000000)
	go Cache(enrichedEventChan, sendEventChan, 10, nil)
	go Enricher(newEventChan, enrichedEventChan, logger)

	if ! *auditd {
		Audisp(newEventChan, os.Stdin)
	} else {
		Auditor(newEventChan, logger)
	}
}
