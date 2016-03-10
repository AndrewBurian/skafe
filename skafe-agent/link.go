package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

func ServerLink(events <-chan AuditEvent, logger *log.Logger) {

	var serverConn net.Conn
	var err error

Connecting:
	for {
		serverConn, err = net.Dial("tcp", "192.168.0.19:6969")
		if err != nil {
			fmt.Println("Connection failed", err)
		} else {
			break Connecting
		}
	}

	encoder := gob.NewEncoder(serverConn)

	for {

		event := <-events

		// glob onto the network
		encoder.Encode(event)
	}
}
