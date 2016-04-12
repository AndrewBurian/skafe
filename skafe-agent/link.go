package main

import (
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

const (
	DEFAULT_PORT uint16 = 6969
)

func ServerLink(events <-chan AuditEvent, conf *AgentConfig) {

	var serverConn net.Conn
	var err error

	strPort := fmt.Sprintf("%d", conf.port)

	// go forever
	for {

		// connection phase
		for {
			// attempt to connect with either TLS or plain
			if conf.tlsConf != nil {
				serverConn, err = tls.Dial("tcp", net.JoinHostPort(conf.addr, strPort), conf.tlsConf)
			} else {
				serverConn, err = net.Dial("tcp", net.JoinHostPort(conf.addr, strPort))
			}

			// check for connection errors
			if err != nil {

				// on failure, try again
				time.Sleep(10 * time.Second)
				continue

			}

			// on success, continue on
			conf.log.Println("Connected to skafe server")
			break
		}

		// create the encoder
		encoder := gob.NewEncoder(serverConn)

		// transmission phase
		for event := range events {

			// gob onto the network
			err := encoder.Encode(event)

			// check if connection lost
			if err != nil {

				// go back to connection phase
				conf.log.Println("Connection to server lost")
				break
			}
		}
	}
}
