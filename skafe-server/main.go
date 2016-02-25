package main

import (
	"fmt"
)

type AuditEvent map[string]string

func main() {

	evChan := make(chan AuditEvent)

	go ClientLink(evChan)

	for {
		fmt.Println(<-evChan)
	}
}
