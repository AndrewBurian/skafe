package main

import (
	"testing"
)

func TestNewScriptWorker(t *testing.T) {
	w, err := NewRbScriptWorker("ruby", "./scripts")
	checkErr(nil, err, t)

	if w == nil {
		t.Errorf("Nil worker returned")
	}
}

func TestNewRbScriptWorkerBadBin(t *testing.T) {
	w, err := NewRbScriptWorker("/bin/false", "./scripts")
	if err == nil {
		t.Errorf("No error returned where one should have been")
	}
	if w != nil {
		t.Errorf("Non-nil worker returned")
	}
}
