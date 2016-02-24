package main

import (
	"fmt"
	"log"
	"os"
)

func Enricher(newEvents <-chan AuditEvent, enrichedEvents chan<- AuditEvent, logger *log.Logger) {

	for {
		// get the event
		event := <-newEvents

		// enrich the event
		GetUser(&event)
		GetParentProcTitle(&event)

		// dispatch the completed event
		enrichedEvents <- event
	}
}

func GetUser(ev *AuditEvent) {
	(*ev)["username"] = "Geoff, probably"
}

func GetParentProcTitle(ev *AuditEvent) {

	// ensure this event has a ppid
	if ppid, ok := (*ev)["ppid"]; ok {

		procFile, err := os.Open("/proc/" + ppid + "/status")
		if err != nil {
			return
		}
		defer procFile.Close()

		var name string

		n, err := fmt.Fscanf(procFile, "Name: %63s", &name)
		if err != nil {
			fmt.Println(err)
		}

		if n == 1 {
			(*ev)["pexe"] = name
		}
	}
}
