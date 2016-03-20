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

func (w *RbScriptWorker) HasFunc(function string) bool {
	//TODO
	return false
}

func (w *RbScriptWorker) Run(function string, ev *AuditEvent) (bool, error) {

	// Send the event to the script engine
	jsonEv, err := json.Marshal(ev)
	if err != nil {
		return false, err
	}

	cmd := fmt.Sprintf("EVENT %s\n", jsonEv)
	_, err = w.stdin.WriteString(cmd)
	if err != nil {
		return false, err
	}

	// Send command to script engine
	cmd = fmt.Sprintf("CMD %s\n", function)
	_, err = w.stdin.WriteString(cmd)
	if err != nil {
		return false, err
	}

	// receive responses
	for {
		// get message from script
		resp, err := w.stdout.ReadString('\n')
		if err != nil {
			return false, err
		}

		// Trim trailing newline
		resp = strings.TrimRight(resp, "\n")

		// Split message into parts
		parts := strings.Split(resp, " ")

		// Base action on first message part
		switch parts[0] {

		// Response with a result
		case "RESP":
			switch parts[1] {
			case "true", "TRUE", "True":
				return true, nil
			case "false", "FALSE", "False":
				return false, nil
			default:
				return false, fmt.Errorf("Unknown response type %s", parts[1])
			}

		// Error occured
		case "ERR":
			return false, fmt.Errorf("An error occured in script execution: %s", resp)

		// Unknown response
		default:
			return false, fmt.Errorf("Unknown command from ruby script worker: %s", resp)
		}
	}

	// Should never arrive here
	return false, fmt.Errorf("Unhandled exception occured")
}
