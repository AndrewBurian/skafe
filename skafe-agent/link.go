package main

import (
	"fmt"
	"log"
	"time"
)

func ServerLink(events <-chan AuditEvent, logger *log.Logger) {

	for {
		time.Sleep(time.Second * 7)
		logger.Println("Server ready to dispatch")

		event := <-events

		logger.Println("Sending event")
		fmt.Println(event)
	}
}
