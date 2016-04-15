package main

import (
	"testing"
)

// Ensure the uid is being resolved
func TestUsername(t *testing.T) {
	name, err := GetUsername("0")
	if err != nil {
		t.Log(err)
		t.Errorf("Username failed to resolve")
	} else if name != "root" {
		t.Log(name)
		t.Errorf("Username incorrectly resolved")
	}
}

// test resolved uid is being placed in event
func TestSourceUser(t *testing.T) {

	ev := AuditEvent{
		"uid": "0",
	}

	GetSourceUser(&ev)

	if key, ok := ev["uname"]; ok {
		if key != "root" {
			t.Log(key)
			t.Error("Incorrect source user name")
		}
	} else {
		t.Log(ev)
		t.Error("user name not resolved")
	}
}

// test full command reconstruction
func TestFullCmd(t *testing.T) {

	ev := AuditEvent{
		"argc": "3",
		"a0":   "test",
		"a1":   "command",
		"a2":   "string",
	}

	GetFullCmd(&ev)

	if key, ok := ev["cmd"]; ok {
		if key != "test command string" {
			t.Log(key)
			t.Errorf("Command incorrectly reconstructed")
		}
	} else {
		t.Log(ev)
		t.Errorf("Failed to reconstruct command")
	}
}

// ensure no extra arguments are being used
func TestFullCmdExtra(t *testing.T) {

	ev := AuditEvent{
		"argc": "3",
		"a0":   "test",
		"a1":   "command",
		"a2":   "string",
		"a3":   "extra",
	}

	GetFullCmd(&ev)

	if key, ok := ev["cmd"]; ok {
		if key != "test command string" {
			t.Log(key)
			t.Errorf("Command incorrectly reconstructed")
		}
	} else {
		t.Log(ev)
		t.Errorf("Failed to reconstruct command")
	}
}
