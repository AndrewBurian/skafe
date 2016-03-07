package main

import (
	"testing"
)

// Ensure a valid cert and key are accepted
func TestGoodTLSSetup(t *testing.T) {

	conf := &ServerConfig{
		tls: true,
		tlsCertPath: "tls/demoserver.pem",
		tlsKeyPath: "tls/demoserver.key",
	}

	err := setupTLS(conf)
	if err != nil {
		t.Fatal(err)
	}
}

// Encure a bad cert crashes
func TestBadTLSSetup(t *testing.T) {

	conf := &ServerConfig{
		tls: true,
		tlsCertPath: "tls/nonexistant.pem",
		tlsKeyPath: "tls/demoserver.key",
	}

	err := setupTLS(conf)
	if err == nil {
		t.Fatalf("Bad cert did not cause an error")
	}
}
