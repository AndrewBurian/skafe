package main

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"strings"
)

func Enricher(newEvents <-chan AuditEvent, enrichedEvents chan<- AuditEvent, conf *AgentConfig) {

	for {
		// get the event
		event := <-newEvents

		// enrich the event
		GetSourceUser(&event)
		GetParentProcTitle(&event)
		GetFullCmd(&event)

		// dispatch the completed event
		enrichedEvents <- event
	}
}

func GetSourceUser(ev *AuditEvent) {
	uid, ok := (*ev)["uid"]
	if !ok {
		return
	}

	uname, err := GetUsername(uid)
	if err != nil {
		return
	}

	(*ev)["uname"] = uname
}

func GetUsername(uid string) (string, error) {
	usr, err := user.LookupId(uid)
	if err != nil {
		return "", err
	}
	return usr.Username, nil
}

func GetFullCmd(ev *AuditEvent) {
	argcStr, ok := (*ev)["argc"]
	if !ok {
		return
	}

	var argc int
	if n, err := fmt.Sscanf(argcStr, "%d", &argc); n != 1 || err != nil {
		return
	}

	var fullCmdBuf bytes.Buffer

	for i := 0; i < argc; i++ {
		argv := fmt.Sprintf("a%d", i)

		if str, ok := (*ev)[argv]; ok {
			fullCmdBuf.WriteString(str + " ")
		}
	}

	(*ev)["cmd"] = strings.TrimRight(fullCmdBuf.String(), " ")
}

func GetParentProcTitle(ev *AuditEvent) {

	// ensure this event has a ppid
	if ppid, ok := (*ev)["ppid"]; ok {
		name, err := GetProcTitle(ppid)
		if err == nil {
			(*ev)["pexe"] = name
		}
	}
}

// Get a process name from its PID
func GetProcTitle(pid string) (string, error) {

	procFile, err := os.Open("/proc/" + pid + "/status")
	if err != nil {
		return "", err
	}
	defer procFile.Close()

	var name string

	n, err := fmt.Fscanf(procFile, "Name: %63s", &name)
	if err != nil {
		return "", err
	}

	if n == 1 {
		return name, nil
	}

	return "", fmt.Errorf("Unknown error")
}
