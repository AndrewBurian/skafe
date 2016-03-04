package main

import (
	"testing"
	"time"
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

func TestRunRules(t *testing.T) {

	flagConfFile := string("./config/skafe-server.conf")

	// Setup server configuration
	conf, err := setupConfig(flagConfFile)
	if err != nil {
		t.Fatalf("Unable to setup server config: %s\n", err.Error())
	}

	// Setup loggers
	err = setupLoggers(conf)
	if err != nil {
		t.Fatalf("Unable to set up log files: %s\n", err.Error())
	}

	conf.serverLog.Println("SKAFE Server started!")
	defer conf.serverLog.Println("SKAFE Server terminated.")

	// setup the rule tree for the rule engine
	baseRule, err := SetupRuleTree(conf)
	if err != nil {
		conf.serverLog.Fatalf("Error loading rule engine: ", err)
		t.FailNow()
	}

	ruleChan := make(chan *AuditEvent)

	go RuleEngine(conf, baseRule, ruleChan)

	ruleChan <- &AuditEvent{
		"key":  "execve",
		"data": "testData",
		"cmd":  "cat /etc/shadow",
	}

	close(ruleChan)
	time.Sleep(3 * time.Second)

}
