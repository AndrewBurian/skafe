package main

import (
	"fmt"
	"log"
)

func ServerLink(events <-chan AuditEvent, logger *log.Logger) {

	for {

		event := <-events

		fmt.Println(event)
	}
}
