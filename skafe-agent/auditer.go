package main

import (
	"time"
)

type AuditEvent map[string]string

func Auditor(receivedEvents chan<- AuditEvent) {
	event := make(AuditEvent)

	event["test"] = "Hello world!"

	for {
		time.Sleep(time.Second * 5)
		receivedEvents <- event
	}
}

func runAudit() {

	// Load all rules
	log.Println("Reading rules file")
	content, err := ioutil.ReadFile("audit.rules.json")
	if err != nil {
		panic(err)
	}

	// Create the new audit socket
	s, err := audit.NewNetlinkConnection()

	if err != nil {
		fmt.Println(err)

	}

	defer s.Close()

	// Enable auditing
	log.Println("Checking audit enabled")
	isEnabled, err := audit.AuditIsEnabled(s)
	if err != nil {
		panic(err)
	}

	if isEnabled != 1 {
		log.Println("Enabling Audit")
		err := audit.AuditSetEnabled(s, 1)
		if err != nil {
			panic(err)
		}
	}

	// Register current pid with audit
	log.Println("Registering PID")
	err = audit.AuditSetPid(s, uint32(syscall.Getpid()))
	if err != nil {
		panic(err)
	}

	// Set audit rules
	log.Println("Setting audit rules")
	err = audit.SetRules(s, content)
	if err != nil {
		panic(err)
	}

	// Audit running
	log.Println("Audit running")

	// setup channel
	errChan := make(chan error)

	// Start Callback
	audit.GetAuditEvents(s, EventCallback, errChan)
	log.Println("Audit event callback started")

	// pop one event
	<-errChan

}

func EventCallback(msg *audit.AuditEvent, ce chan error, args ...interface{}) {
	fmt.Println(msg)
}
