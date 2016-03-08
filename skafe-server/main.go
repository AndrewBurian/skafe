package main

import (
	"flag"
	"log"
)

func main() {

	// flag for specifying a non-default config directory
	flagConfFile := flag.String("conf", "/etc/skafe/skafe-server.conf", "Set the directory containing config and rules")

	flag.Parse()

	// Setup server configuration
	conf, err := setupConfig(*flagConfFile)
	if err != nil {
		log.Fatalf("Unable to setup server config: %s\n", err.Error())
	}

	// Setup loggers
	err = setupLoggers(conf)
	if err != nil {
		log.Fatalf("Unable to set up log files: %s\n", err.Error())
	}

	// setup TLS
	err = setupTLS(conf)
	if err != nil {
		log.Fatalf("Unable to set up TLS: %s\n", err.Error())
	}

	conf.serverLog.Println("SKAFE Server started!")
	defer conf.serverLog.Println("SKAFE Server terminated.")

	// setup the rule tree for the rule engine
	baseRule, err := SetupRuleTree(conf)
	if err != nil {
		conf.serverLog.Fatalf("Error loading rule engine: ", err)
	}

	evChan := make(chan *AuditEvent)
	ruleChan := make(chan *AuditEvent)

	go ClientLink(conf, evChan)

	go RuleEngine(conf, baseRule, ruleChan)

	QueueEvents(conf, evChan, ruleChan)

}
