package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type RbScriptWorker struct {
	lang   string
	cmd    *exec.Cmd
	stdin  bufio.Writer
	stdout bufio.Reader
}

func NewRbScriptWorker(bin string) (ScriptWorker, error) {

	worker := &RbScriptWorker{
		lang: SCRIPT_LANG_RUBY,
	}

	// Exec the ruby base script
	worker.cmd = exec.Command(bin, "scripts/ruby.rb")

	if err := worker.cmd.Start(); err != nil {
		return nil, err
	}

	return worker, nil
}

func (w *RbScriptWorker) Lang() string {
	return SCRIPT_LANG_RUBY
}

func (w *RbScriptWorker) Run(function string, ev *AuditEvent) bool {

	// Send the event to the script engine
	jsonEv, err := json.Marshal(ev)
	if err != nil {
		panic(err) //TODO
	}

	cmd := fmt.Sprintf("EVENT %s\n", jsonEv)
	_, err = w.stdin.WriteString(cmd)
	if err != nil {
		panic(err)
	}

	// Send command to script engine
	cmd = fmt.Sprintf("CMD %s\n", function)
	_, err = w.stdin.WriteString(cmd)
	if err != nil {
		panic(err) //TODO error handling
	}

	// receive responses
	for {
		// get message from script
		resp, err := w.stdout.ReadString('\n')
		if err != nil {
			panic(err) //TODO error handling
		}

		// Trim trailing newline
		resp = strings.TrimRight(resp, "\n")

		// Split message into parts
		parts := strings.Split(resp, " ")

		switch parts[0] {
		case "RESP":
			if parts[1] == "true" {
				return true
			} else {
				return false
			}
		default:
			panic("WHAT IS GOING ON") //TODO error handling
		}
	}

	// Should never arrive here
	panic("...")
	return false
}
