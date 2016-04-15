package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {

	auditd := flag.Bool("audit", false, "Run as system auditer")
	confFile := flag.String("conf", "/etc/skafe/skafe-agent.conf", "Config file path")
	flag.Parse()

	// setup config
	conf, err := setupConfig(*confFile)
	if err != nil {
		log.Fatal(err)
	}

	conf.log.Println("Started skafe-agent")
	defer conf.log.Println("Terminated skafe-agent")

	// create channels to pass events between process steps
	newEventChan := make(chan AuditEvent) // auditer -> enricher

	enrichedEventChan := make(chan AuditEvent) // enricher -> cache

	sendEventChan := make(chan AuditEvent) // cache -> rateLimiter

	serverChan := make(chan AuditEvent) // rateLimiter -> serverLink

	go ServerLink(serverChan, conf)
	go RateLimit(sendEventChan, serverChan, 10, time.Second)
	go Cache(enrichedEventChan, sendEventChan, 10, nil)
	go Enricher(newEventChan, enrichedEventChan, conf)

	if !*auditd {
		Audisp(newEventChan, os.Stdin)
	} else {
		Auditor(newEventChan, conf)
	}
}
