package main

import (
	//audit "github.com/andrewburian/libaudit-go"
	//"fmt"
	audit "github.com/mozilla/libaudit-go"
	"io/ioutil"
	"log"
	"syscall"
	//"time"
)

type AuditEvent map[string]string

func Auditor(receivedEvents chan<- AuditEvent, logger *log.Logger) {

	newEvents, err := setupAudit()
	if err != nil {
		panic(err)
	}

	for ev := range newEvents {
		receivedEvents <- ev.Data
	}
}

func setupAudit() (<-chan *audit.AuditEvent, error) {

	// Load all rules
	log.Println("Reading rules file")
	content, err := ioutil.ReadFile("audit.rules.json")
	if err != nil {
		return nil, err
	}

	// Create the new audit socket
	s, err := audit.NewNetlinkConnection()
	if err != nil {
		return nil, err

	}

	// Enable auditing
	log.Println("Checking audit enabled")
	isEnabled, err := audit.AuditIsEnabled(s)
	if err != nil {
		return nil, err
	}

	if isEnabled != 1 {
		log.Println("Enabling Audit")
		err := audit.AuditSetEnabled(s, 1)
		if err != nil {
			return nil, err
		}
	}

	// Register current pid with audit
	log.Println("Registering PID")
	err = audit.AuditSetPid(s, uint32(syscall.Getpid()))
	if err != nil {
		return nil, err
	}

	// Set audit rules
	log.Println("Setting audit rules")
	err = audit.SetRules(s, content)
	if err != nil {
		return nil, err
	}

	// Audit running
	log.Println("Audit running")

	// setup channels
	errChan := make(chan error)
	eventChan := make(chan *audit.AuditEvent)

	// Start Callback
	audit.GetAuditEvents(s, EventCallback, errChan, eventChan)
	log.Println("Audit event callback started")

	return eventChan, nil
}

func EventCallback(msg *audit.AuditEvent, ce chan error, args ...interface{}) {
	eventChan := args[0].(chan *audit.AuditEvent)
	eventChan <- msg
}
