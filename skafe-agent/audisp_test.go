package main

import (
	"bytes"
	"testing"
)

var (
	validEvent string = `audit(1364481363.243:24287): arch=c000003e syscall=2 success=no exit=-13 a0=7fffd19c5592 a1=0 a2=7fffd19c4b50 a3=a items=1 ppid=2686 pid=3538 auid=500 uid=500 gid=500 euid=500 suid=500 fsuid=500 egid=500 sgid=500 fsgid=500 tty=pts0 ses=1 comm="cat" exe="/bin/cat" subj=unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023 key="sshd_config"`

	fullEvent string = `type=DAEMON_START msg=audit(1460442803.395:4215): auditd start, ver=2.4.4 format=raw kernel=4.4.5-1-custom auid=4294967295 pid=5753 res=success`
)

// Check that the pure audit message is parsed properly
func TestAuditParseValid(t *testing.T) {

	data := validEvent

	_, err := ParseEvent(data)
	if err != nil {
		t.Errorf("Failed to parse event: %s", err)
	}

}

// Ensure it accepts auditd style messages
func TestAuditParseFull(t *testing.T) {

	data := fullEvent

	_, err := ParseEvent(data)
	if err != nil {
		t.Errorf("Failed to parse full event")
		t.Log(err)
	}
}

// Ensure it strips newlines properly
func TestAuditParteFullNewline(t *testing.T) {

	data := fullEvent + "\n"
	_, err := ParseEvent(data)
	if err != nil {
		t.Log(err)
		t.Errorf("Failed to parse event")
	}

}

// Audisp core logic is behaving as it should
func TestAudispValid(t *testing.T) {

	buf := bytes.NewBufferString(validEvent)

	evChan := make(chan AuditEvent, 2)

	err := Audisp(evChan, buf)

	if err != nil {
		t.Fatal(err)
	}

	select {
	case _, ok := <-evChan:
		if !ok {
			t.Fatalf("No event received")
		}
	default:
		t.Fatal("Event not received")
	}

	select {
	case _, ok := <-evChan:
		if ok {
			t.Fatalf("More than one event received")
		}
	default:
		t.Fatalf("Chanel not closed")
	}

}
