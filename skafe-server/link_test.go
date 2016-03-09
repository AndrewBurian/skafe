package main

import (
	"crypto/tls"
	"encoding/gob"
	"net"
	"sync"
	"testing"
)

func TestClientRejectNoTls(t *testing.T) {
	conf, err := setupConfig("tests/tls_goodca.conf")
	checkErr(nil, err, t)

	err = setupTLS(conf)
	checkErr(nil, err, t)

	listenConn, err := tls.Listen("tcp", "127.0.0.1:6969", conf.tlsConf)
	checkErr(nil, err, t)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	// Server
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c, err := listenConn.Accept()
		if err == nil {
			t.Log("Server Accepted insecure connection")
		}
		if err != nil {
			return
		}
		ev := make(chan *AuditEvent, 1)
		HandleClient(c, ev)
		select {
		case data := <-ev:
			t.Errorf("Server received insecure data %s", data)
		default:
		}
	}(wg)

	defer listenConn.Close()

	// Client
	conn, err := net.Dial("tcp", "127.0.0.1:6969")
	if err == nil {
		t.Log("Client established insecure connection")
	}
	if err != nil {
		return
	}

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(map[string]string{
		"key": "val",
	})
	if err == nil {
		t.Log("Client send insecure data")
	}
	conn.Close()
}

func TestClientRejectUntrusted(t *testing.T) {
	conf, err := setupConfig("tests/tls_goodca.conf")
	checkErr(nil, err, t)

	err = setupTLS(conf)
	checkErr(nil, err, t)

	listenConn, err := tls.Listen("tcp", "127.0.0.1:6969", conf.tlsConf)
	checkErr(nil, err, t)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	// Server
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c, err := listenConn.Accept()
		if err == nil {
			t.Log("Server Accepted insecure connection")
		}
		if err != nil {
			return
		}
		ev := make(chan *AuditEvent, 1)
		HandleClient(c, ev)
		select {
		case data := <-ev:
			t.Errorf("Server received insecure data %s", data)
		default:
		}
	}(wg)

	defer listenConn.Close()

	// Client
	tlsConf := &tls.Config{
		ServerName:         "demo.server.local",
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:6969", tlsConf)
	if err == nil {
		t.Log("Client established insecure connection")
	}
	if err != nil {
		return
	}

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(map[string]string{
		"key": "val",
	})
	if err == nil {
		t.Log("Client send insecure data")
	}
	conn.Close()
}

func TestClientAcceptTrusted(t *testing.T) {
	conf, err := setupConfig("tests/tls_goodca.conf")
	checkErr(nil, err, t)

	err = setupTLS(conf)
	checkErr(nil, err, t)

	listenConn, err := tls.Listen("tcp", "127.0.0.1:6969", conf.tlsConf)
	checkErr(nil, err, t)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	// Server
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c, err := listenConn.Accept()
		if err != nil {
			t.Errorf("Server failed to accept connection")
			return
		}
		ev := make(chan *AuditEvent, 1)
		HandleClient(c, ev)
		select {
		case _ = <-ev:
		default:
			t.Errorf("Server failed to receive data")
		}
	}(wg)

	defer listenConn.Close()

	// Client
	cert, err := tls.LoadX509KeyPair("../skafe-agent/tls/demoagent.pem", "../skafe-agent/tls/demoagent.key")
	checkErr(nil, err, t)

	tlsConf := &tls.Config{
		ServerName:         "demo.server.local",
		InsecureSkipVerify: true,
	}

	tlsConf.Certificates = append(tlsConf.Certificates, cert)

	conn, err := tls.Dial("tcp", "127.0.0.1:6969", tlsConf)
	if err != nil {
		t.Errorf("Client failed to connect: %s", err)
		return
	}

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(map[string]string{
		"key": "val",
	})
	if err != nil {
		t.Errorf("Client failed to send data: %s", err)
	}
	conn.Close()
}
