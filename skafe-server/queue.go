package main

import ()

// Queue for handling incoming events
func QueueEvents(conf *ServerConfig, in <-chan *AuditEvent, out chan<- *AuditEvent) {

	for {
		ev := <-in
		out <- ev
	}
}
