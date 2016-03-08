package main

import (
	"log"
	"regexp"
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

func TestMatchNodeBlankEvent(t *testing.T) {

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

func TestMatchNodeNoMatch(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: LOG,
		action:  MATCH,
		matches: map[string]*regexp.Regexp{
			"key": regexp.MustCompile("other val"),
		},
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key": "val",
	}

	RunNode(conf, node, event)

	if counter.Count != 0 {
		t.Errorf("Node did not trigger correct logs\nExpecting [0]\nGot [%d]", counter.Count)
	}
}

func TestMatchNodeHalfMatch(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		trigger: LOG,
		action:  MATCH,
		matches: map[string]*regexp.Regexp{
			"key":  regexp.MustCompile("val"),
			"key2": regexp.MustCompile("other val2"),
		},
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key":  "val",
		"key2": "val2",
	}

	RunNode(conf, node, event)

	if counter.Count != 0 {
		t.Errorf("Node did not trigger correct logs\nExpecting [0]\nGot [%d]", counter.Count)
	}
}

func TestRecursiveMatch(t *testing.T) {

	counter := &LogCounter{
		t: t,
	}
	logger := log.New(counter, "", 0)

	node := &RuleNode{
		name:    "first",
		trigger: LOG,
		action:  MATCH,
		matches: map[string]*regexp.Regexp{
			"key":  regexp.MustCompile("val"),
			"key2": regexp.MustCompile("val2"),
		},
		nodes: []*RuleNode{
			&RuleNode{
				name:   "second",
				trigger: LOG,
				action:  MATCH,
				matches: map[string]*regexp.Regexp{
					"key":  regexp.MustCompile("val"),
					"key2": regexp.MustCompile("val2"),
				},
			},
		},
	}

	conf := &ServerConfig{
		eventLog: logger,
	}

	event := &AuditEvent{
		"key":  "val",
		"key2": "val2",
	}

	RunNode(conf, node, event)

	if counter.Count != 2 {
		t.Errorf("Node did not trigger correct logs\nExpecting [2]\nGot [%d]", counter.Count)
	}
}
