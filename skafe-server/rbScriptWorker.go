package main

import (
	"io"
	"os/exec"
)

type RbScriptWorker struct {
	lang   string
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
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

func (w *RbScriptWorker) Run(function string) bool {

	//TODO

	return false
}
