package main

import (
	"log"
	"testing"
)

type LogCounter struct {
	Count int
	t     *testing.T
}

func (l *LogCounter) Write(b []byte) (int, error) {
	l.Count += 1
	l.t.Log(string(b))
	return len(b), nil
}

func TestMatchNodeBlank(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: LOG,
		action:  MATCH,
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key": "val",
	}

	RunNode(conf, node, event)

	if counter.Count != 1 {
		t.Errorf("Node did not trigger correct logs\nExpecting [1]\nGot [%d]", counter.Count)
	}
}

func TestMatchNodeBlankBoth(t *testing.T) {
	eventCounter := &LogCounter{
		t: t,
	}
	alertCounter := &LogCounter{t: t}

	eventLogger := log.New(eventCounter, "", 0)
	alertLogger := log.New(alertCounter, "", 0)

	node := &RuleNode{
		trigger: BOTH,
		action:  MATCH,
	}

	conf := &ServerConfig{
		eventLog: eventLogger,
		alertLog: alertLogger,
	}

	event := &AuditEvent{
		"key": "val",
	}

	RunNode(conf, node, event)

	if eventCounter.Count != 1 {
		t.Errorf("Node did not trigger correct event logs\nExpecting [1]\nGot [%d]", eventCounter.Count)
	}
	if alertCounter.Count != 1 {
		t.Errorf("Node did not trigger correct alert logs\nExpecting [1]\nGot [%d]", alertCounter.Count)
	}
}

func TestMatchNodeBlankAlert(t *testing.T) {
	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: ALERT,
		action:  MATCH,
	}

	conf := &ServerConfig{
		alertLog: logger,
	}

	event := &AuditEvent{
		"key": "val",
	}

	RunNode(conf, node, event)

	if counter.Count != 1 {
		t.Errorf("Node did not trigger correct logs\nExpecting [1]\nGot [%d]", counter.Count)
	}
}
