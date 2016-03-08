package main

import (
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"net"
)

const (
	DEFAULT_PORT uint16 = 6969
)

// Listen for connecting clients
func ClientLink(conf *ServerConfig, incomingEvents chan<- *AuditEvent) {
	listenStr := fmt.Sprintf(":%d", conf.port)

	var listenConn net.Listener
	var err error

	if conf.tls {
		listenConn, err = tls.Listen("tcp", listenStr, conf.tlsConf)
	} else {
		listenConn, err = net.Listen("tcp", listenStr)
	}

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

// Handle each individual client
func HandleClient(client net.Conn, events chan<- *AuditEvent) {

	defer client.Close()

	decoder := gob.NewDecoder(client)

	for {
		// create the audit event
		ev := &AuditEvent{}

		// decode into new event
		err := decoder.Decode(ev)
		if err != nil {
			fmt.Println(err)
			break
		}

		events <- ev
	}
}
