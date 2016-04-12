package main

import (
	"testing"
)

// Test sane config
func TestConfig(t *testing.T) {

	conf, err := setupConfig("tests/sane.conf")
	if err != nil {
		t.Log(err)
		t.Fatalf("Error returned from sane config")
	}
	if conf == nil {
		t.Fatalf("Nil config returned")
	}

	if conf.port != 9999 {
		t.Errorf("Wrong port set in config")
		t.Log(conf.port)
	}

	if conf.addr != "skafe-server.local" {
		t.Errorf("Wrong addr set")
		t.Log(conf.addr)
	}

}

// Broken config file
func TestNoConfing(t *testing.T) {

	conf, err := setupConfig("/dev/null/broken.conf")
	if err == nil {
		t.Errorf("Expecting failed to open file error")
	}
	if conf != nil {
		t.Errorf("Expecting nil conf returned on file error")
	}
}

// No addr specified
func TestNoAddr(t *testing.T) {

	conf, err := setupConfig("tests/noserver.conf")
	if err == nil {
		t.Errorf("Expecting no addr error")
	}
	if conf != nil {
		t.Errorf("Expecting nil conf on error")
	}
}

// Port is NaN
func TestBadPort(t *testing.T) {

	conf, err := setupConfig("tests/badport.conf")
	if err == nil {
		t.Errorf("Expecting bad port error")
	}
	if conf != nil {
		t.Errorf("Expecting nil conf on bad port error")
	}
}

// Port > 65534
func TestHighPort(t *testing.T) {

	conf, err := setupConfig("tests/highport.conf")
	if err == nil {
		t.Errorf("Expecting bad port error")
	}
	if conf != nil {
		t.Errorf("Expecting nil conf on bad port error")
	}

}

// Logger file cannot be created
func TestLogFileBroken(t *testing.T) {

	conf, err := setupConfig("tests/brokenlog.conf")
	if err == nil {
		t.Errorf("Expecting bad log file error")
	}
	if conf != nil {
		t.Errorf("Expecting nil conf on bad log file")
	}
}
