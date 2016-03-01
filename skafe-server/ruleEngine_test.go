package main

import (
	"log"
	"os"
	"testing"
)

func TestRuleTreeSetup(t *testing.T) {

	conf := &ServerConfig{
		rulesDirPath: "./config/rules",
		serverLog:    log.New(os.Stdout, "Server log: ", log.LstdFlags),
	}

	_, err := SetupRuleTree(conf)

	if err != nil {
		t.Fatal(err.Error())
	}
}
