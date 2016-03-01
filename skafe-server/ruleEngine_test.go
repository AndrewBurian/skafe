package main

import (
	"testing"
)

func TestRuleTreeSetup(t *testing.T) {

	conf, err := setupConfig("config/skafe-server.conf")
	if err != nil {
		t.Fatal(err)
	}

	err = setupLoggers(conf)
	if err != nil {
		t.Fatal(err)
	}

	_, err = SetupRuleTree(conf)

	if err != nil {
		t.Fatal(err.Error())
	}
}
