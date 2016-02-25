package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

func ClientLink(incomingEvents chan<- AuditEvent) {
	listenConn, err := net.Listen("tcp", ":6969")
	if err != nil {
		panic(err)
	}

	for {
		newConn, err := listenConn.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go HandleClient(newConn, incomingEvents)

	}

}

func HandleClient(client net.Conn, events chan<- AuditEvent) {

	defer client.Close()

	decoder := gob.NewDecoder(client)

	var ev AuditEvent

	for {
		err := decoder.Decode(&ev)
		if err != nil {
			fmt.Println(err)
			break
		}

		events <- ev
	}
}
