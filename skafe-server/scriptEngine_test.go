package main

import (
	"testing"
)

/*
func Test(t *testing.T){

}
*/

func TestSetupScriptPoolRuby(t *testing.T) {
	conf := &ServerConfig{
		execRuby: "ruby",
		//scriptsdir_ruby: "scripts",
	}

	pool, err := SetupScriptPool(conf)
	checkErr(nil, err, t)

	if pool == nil {
		t.Fatalf("Nil script pool returned")
	}

	rubycmd, ok := pool.execCmd[SCRIPT_LANG_RUBY]
	if !ok {
		t.Fatalf("Ruby command not added")
	}
	if rubycmd != "ruby" {
		t.Fatalf("Ruby command incorrect")
	}
}

func TestGetWorkerWrongLang(t *testing.T) {
	pool := &ScriptPool{}

	worker, err := pool.GetWorker("erlang")
	if err == nil {
		t.Errorf("No error returned where one was expected")
	}

	if worker != nil {
		t.Errorf("Non nil worker returned")
	}
}
